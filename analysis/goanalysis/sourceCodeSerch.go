package goanalysis

import (
	"codeanalysis/analysis/dao"
	"fmt"
	"strings"
)

type source struct {
	buf         []byte           // 檔案內文資料
	b, r, e     int              // 掃描錨點資料
	ch          rune             // 當前文字資料
	isEnd       bool             // 掃描進度
	PackageInfo *dao.PackageInfo // 檔案資訊
}

//========== Import =================

// 解析 import 區塊
//
// @return []*PackageLink 未關聯的 import 資料
func (s *source) ImportDeclarations() []*dao.PackageLink {
	var packageLinks []*dao.PackageLink
	if s.buf[s.r+1] == '(' {
		s.nextCh()
		s.toNextCh()

		for {
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}

			if string(s.buf[s.r+1:s.r+3]) == "//" {
				s.nextTargetToken('\n')
				continue
			}
			name, importPath, packageInfo := s.importSpec()
			packageLinks = append(packageLinks, dao.NewPackageLink(name, importPath, packageInfo))
			s.toNextCh()

		}
	} else {
		name, importPath, packageInfo := s.importSpec()
		packageLinks = append(packageLinks, dao.NewPackageLink(name, importPath, packageInfo))

	}
	return packageLinks
}

// 解析 import 內文
//
// @return name 		import檔案的替換名稱
// @return importPath	import檔案的路徑
func (s *source) importSpec() (name, importPath string, packageInfo *dao.PackageInfo) {
	s.nextToken()
	name = strings.TrimSpace(s.rangeStr())
	s.nextTargetToken('"')
	importPath = strings.TrimSpace(s.rangeStr())

	packageInfo, ok := Instants.AllPackageMap[importPath]
	if !ok {
		packageInfo = dao.NewPackageInfo(importPath)
		Instants.AllPackageMap[importPath] = packageInfo
	}

	return name, importPath, packageInfo
}

//========== Const =================

// ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
// ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .
// IdentifierList = identifier { "," identifier } .
// ExpressionList = Expression { "," Expression } .
func (s *source) ConstantDeclarations() {
	var constSpec [][]string
	if s.buf[s.r+1] == '(' {
		if s.buf[s.r+2] == ')' {
			fmt.Println(":: const () ::")
			return
		}
		s.nextCh()

		for {
			s.toNextCh()
			if s.CheckCommon() {
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			} else {
				constSpec = append(constSpec, s.ConstSpec())
			}

			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		constSpec = append(constSpec, s.ConstSpec())
	}
	fmt.Println(":: const (", constSpec, ") ::")
}

func (s *source) ConstSpec() []string {
	var datas []string
	nextIdx := s.nextIdx()

	if s.buf[nextIdx] == '\n' {
		s.next()
		name := s.rangeStr()
		s.toNextCh()

		datas = append(datas, name, "_", "_")

		if s.CheckCommon() {
			s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		}
		return datas

	}
	if s.buf[nextIdx-1] == ',' { // 多參數初始化
		var nameList []string
		for {
			s.nextToken()
			nameList = append(nameList, s.rangeStr())
			s.nextCh()
			nextIdx = s.nextIdx()
			if s.buf[nextIdx-1] != ',' {
				break
			}
		}
		s.next()
		nameList = append(nameList, s.rangeStr())
		s.toNextCh()
		if s.buf[s.r+1] == '=' {
			s.next()
			expressionList := s.OnExpressionList(len(nameList))
			datas = append(datas, nameList...)
			datas = append(datas, "_")
			datas = append(datas, expressionList...)
		} else {
			typeName := s.OnType()
			expressionList := s.OnExpressionList(len(nameList))
			datas = append(datas, nameList...)
			datas = append(datas, typeName)
			datas = append(datas, expressionList...)
		}
	} else { // 單一參數初始化

		s.next()
		name := s.rangeStr()
		s.toNextCh()

		switch s.buf[s.r+1] {
		case '=': // 隱藏型態
			s.next()
			expressionList := s.OnExpressionList(1)
			datas = append(datas, name, "_")
			datas = append(datas, expressionList...)

		case '/': // Common
			if s.CheckCommon() {
				datas = append(datas, name, "_", "_")
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
				return datas
			}
		case '\n': // 段落已結束
			datas = append(datas, name, "_", "_")
			return datas

		default: // 標準初始化
			typeName := s.OnType()
			s.next()
			expressionList := s.OnExpressionList(1)
			datas = append(datas, name, typeName)
			datas = append(datas, expressionList...)
		}

	}
	if s.CheckCommon() {
		s.OnComments(string(s.buf[s.r+1 : s.r+3]))
	}
	return datas
}

func (s *source) OnExpressionList(count int) (strLit []string) {
	str := ""
	if count == 1 {

		for {
			if s.ch == '\n' {
				break
			} else if s.buf[s.r+1] == '`' {
				str += s.OnJsonTag()
			} else if string(s.buf[s.r+1:s.r+3]) == "//" {
				s.nextTargetToken('\n')
				break
			} else {
				s.next()
				str += s.rangeStr()
				break
			}
			// s.next()
			// str += s.rangeStr()
			// if s.ch == '\n' {
			// 	break
			// } else if string(s.buf[s.r+1:s.r+3]) == "//" {
			// 	s.nextTargetToken('\n')
			// 	break
			// }
		}
		if str != "" {
			strLit = append(strLit, str)
		}
	} else {
		for i := 0; i < count-1; i++ {
			s.nextTargetToken(',')
			strLit = append(strLit, s.rangeStr())
			s.nextCh()
		}
		for {
			s.next()
			str += s.rangeStr()
			if s.ch == '\n' {
				break
			} else if string(s.buf[s.r+1:s.r+3]) == "//" {
				s.nextTargetToken('\n')
				break
			}
		}
		strLit = append(strLit, str)
	}
	return
}

//========== Var =================

// VarDecl = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
// VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
func (s *source) VariableDeclarations() {
	var varSpec [][]string
	if s.buf[s.r+1] == '(' {
		if s.buf[s.r+2] == ')' {
			fmt.Println(":: var () ::")
			return
		}
		s.nextCh()

		for {
			s.toNextCh()
			varSpec = append(varSpec, s.VarSpec())
			s.toNextCh()
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		varSpec = append(varSpec, s.VarSpec())
	}
	fmt.Println(":: var (", varSpec, ") ::")
}

// ShortVarDecl = IdentifierList ":=" ExpressionList .
func (s *source) ShortVariableDeclarations() {

}

func (s *source) VarSpec() []string {
	var datas []string
	nextIdx := s.nextIdx()

	if s.buf[nextIdx] == '\n' {
		s.next()
		name := s.rangeStr()
		s.toNextCh()

		datas = append(datas, name, "_", "_")
		return datas

	}
	if s.buf[nextIdx-1] == ',' {
		var nameList []string
		for {
			s.nextToken()
			nameList = append(nameList, s.rangeStr())
			s.nextCh()
			nextIdx = s.nextIdx()
			if s.buf[nextIdx-1] != ',' {
				break
			}
		}
		s.next()
		nameList = append(nameList, s.rangeStr())
		s.toNextCh()
		if s.buf[s.r+1] == '=' {
			s.next()
			expressionList := s.OnExpressionList(len(nameList))
			datas = append(datas, nameList...)
			datas = append(datas, "_")
			datas = append(datas, expressionList...)
		} else {
			typeName := s.OnType()
			expressionList := s.OnExpressionList(len(nameList))
			datas = append(datas, nameList...)
			datas = append(datas, typeName)
			datas = append(datas, expressionList...)
		}
		return datas
	}

	s.next()
	name := s.rangeStr()
	s.toNextCh()

	if s.buf[s.r+1] == '=' {
		s.nextCh()
		expressionList := s.OnExpressionList(1)
		datas = append(datas, name, "_")
		datas = append(datas, expressionList...)
		return datas
	}

	s.OnDeclarationsType()
	var typeName string
	s.toNextCh()
	expressionList := s.OnExpressionList(1)
	datas = append(datas, name, typeName)
	datas = append(datas, expressionList...)

	return datas
}

//========== Type =================

// Type declarations
// TypeDecl = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
// TypeSpec = AliasDecl | TypeDef .
// AliasDecl = identifier "=" Type .
// TypeDef = identifier Type .

// 解析宣告區塊
//
// @return []*StructInfo 宣告內容
func (s *source) TypeDeclarations() []dao.IStructInfo {
	var structs []dao.IStructInfo
	if s.buf[s.r+1] == '(' {
		s.nextCh()
		s.toNextCh()

		for {
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}

			if string(s.buf[s.r+1:s.r+3]) == "//" {
				s.nextTargetToken('\n')
				continue
			}
			structs = append(structs, s.TypeSpec())
			s.toNextCh()

		}
	} else {

		structs = append(structs, s.TypeSpec())
	}
	return structs
}

// 處理 Type 宣告
// AliasDecl = identifier "=" Type .
// TypeDef = identifier Type .
// ex: _type t1 string
// r="_"
func (s *source) TypeSpec() (structInfo dao.IStructInfo) {

	s.next()
	name := strings.TrimSpace(s.rangeStr())
	if s.buf[s.r+1] == '=' {
		s.next()
	}
	s.toNextCh()
	structInfo = s.OnDeclarationsType()
	structInfo.SetName(name)
	return structInfo
}

func (s *source) typeformat() {
	var typeSpec [][]string
	if s.buf[s.r+1] == '(' {
		if s.buf[s.r+2] == ')' {
			fmt.Println(":: type () ::")
			return
		}
		s.next()
		for {
			if s.buf[s.r+1] == ')' {
				s.nextCh()
				break
			}
			s.next()
			name := strings.TrimSpace(s.rangeStr())
			if name == "" {
				continue
			}
			if s.buf[s.r+1] == '=' {
				s.next()
			}
			s.OnDeclarationsType()
			var typename string
			typeSpec = append(typeSpec, []string{name, typename})
		}
	} else {
		s.next()
		name := strings.TrimSpace(s.rangeStr())
		if s.buf[s.r+1] == '=' {
			s.next()
		}
		s.OnDeclarationsType()
		var typename string
		typeSpec = append(typeSpec, []string{name, typename})
	}
	fmt.Println(":: type (", typeSpec, ") ::")
}

//========== Func =================

// FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
// FunctionName = identifier .
// FunctionBody = Block .
func (s *source) FunctionDeclarations() {
	name := s.OnFuncName()
	signatures, results := s.Signature()
	s.funcBodyBlock()
	fmt.Println(":: func (", name, signatures, results, ") ::")
}

func (s *source) Signature() (Signatures [][]string, Results [][]string) {
	Signatures = s.OnParams(true)
	Results = s.OnResult(' ')
	return
}

// _{
//   return
//  }
// b=任意 r="_"
func (s *source) funcBodyBlock() {
	var blockCount int

	s.nextToken()
	for {
		s.toNextToken()
		switch s.buf[s.r+1] {
		case '}':
			s.nextCh()
			if blockCount > 0 {
				blockCount--
				continue
			}
			return
		case '\'':
			s.nextCh()
			if s.buf[s.r+1] == '\\' { // 跳脫字元
				s.nextCh()
				s.nextCh()
			}
			s.nextTargetToken('\'')
		case '"':
			s.nextCh()
			s.nextTargetToken('"')
		case '/':
			commentType := string(s.buf[s.r+1 : s.r+3])
			if commentType == "//" || commentType == "/*" {
				s.OnComments(commentType)
			}
		case '{':
			s.nextCh()
			blockCount++
		default:
			s.nextCh()
		}
	}
}

// MethodDecl = "func" Receiver MethodName Signature [ FunctionBody ] .
// Receiver   = Parameters .
func (s *source) MethodDeclarations() {
	s.nextCh()
	receiver := s.OnReceiver()
	name := s.OnFuncName()
	if name == "OnDeclarationsResult" {
		fmt.Println("")
	}
	signatures, results := s.Signature()
	if string(s.buf[s.r]) != " " {
		fmt.Println("")
	}
	s.funcBodyBlock()
	fmt.Println(":: method (", receiver, name, signatures, results, ") ::")
}

func (s *source) OnReceiver() [][]string {
	// s.nextTargetToken(')')
	// receivers := strings.TrimSpace(s.rangeStr())
	// receiversLost := strings.Split(receivers, " ")

	// if len(receiversLost) == 1 {
	// 	seltype = receiversLost[0]
	// 	isPoint = receiversLost[0][0] == '*'
	// } else {
	// 	selname = receiversLost[0]
	// 	seltype = receiversLost[1]
	// 	isPoint = receiversLost[1][0] == '*'
	// }
	// s.nextCh()

	return s.OnParams(true)
}

// func 名稱
func (s *source) OnFuncName() string {
	s.nextToken()
	return strings.TrimSpace(s.rangeStr())
}

// 處理 func 多參數資料
/* @params isParameters 是否為輸入參數
 * 進入指標應當只在 abc(w *int) (x, z int)
 * 輸入參數指標 b=任意 r="("
 * =====================
 * 輸出參數指標 b=任意 r="("
 */
func (s *source) OnParams(isParameters bool) [][]string {
	strLit := [][]string{}
	var nextIdx int = s.r
	var unknowCount int
	unknowData := []string{}

	if s.buf[s.r+1] == ')' {
		// case in
		//	func name()
		s.nextCh()
		s.nextCh()
		return strLit
	}

	for isNext := true; isNext; {
		nextIdx = s.nextIdxByPosition(nextIdx)

		switch s.buf[nextIdx-1] {
		case ',':
			unknowCount++

		case ')':
			// 隱藏名稱宣告
			for i := 0; i < unknowCount; i++ {
				data := s.OnType(')')
				s.nextCh()
				strLit = append(strLit, []string{"", data})
			}
			strLit = append(strLit, []string{"", s.OnType(')')})
			unknowCount = 0

			isNext = false

		default:
			// 存在簡易變數宣告
			if unknowCount > 0 {
				for i := 0; i < unknowCount; i++ {
					s.nextToken()
					data := strings.TrimSpace(s.rangeStr())
					unknowData = append(unknowData, data)
				}
				s.nextCh()
				unknowCount = 0
			}
			s.next()
			paramName := strings.TrimSpace(s.rangeStr())
			paramType := s.OnType(')')

			for _, name := range unknowData {
				strLit = append(strLit, []string{name, paramType})
			}
			strLit = append(strLit, []string{paramName, paramType})
			unknowData = make([]string, 0)

			if s.ch == ')' {
				isNext = false
			} else {
				s.next()
				nextIdx = s.r
			}
		}
	}
	s.nextCh()
	return strLit
}

// func 輸出參數
func (s *source) OnResult(endTag ...byte) [][]string {
	strLit := [][]string{}
	switch s.buf[s.r+1] {
	case '{', '\t', '}':
		return strLit
	case '(':
		s.nextCh()
		strLit = s.OnParams(true)

	default:
		strLit = append(strLit, []string{"", s.OnType(endTag[0])})
	}
	return strLit
}

// 確認後續是否為註解
// 指標指定位置 _//
// r="_"
func (s *source) CheckCommon() bool {
	if s.r+1 == s.e {
		return false
	}

	if s.buf[s.r+1] != '/' {
		return false
	}

	previrwCh := string(s.buf[s.r+1 : s.r+3])
	if previrwCh != "//" && previrwCh != "/*" {
		return false
	}

	return true
}
