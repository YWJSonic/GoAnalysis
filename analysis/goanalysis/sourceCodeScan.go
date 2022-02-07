package goanalysis

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/tool"
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

// 解析 const 區塊
func (s *source) ConstantDeclarations() {
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
				infos := s.ConstSpec()
				for _, info := range infos {
					s.PackageInfo.AllTypeInfos[info.Name] = info
				}
			}

			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		infos := s.ConstSpec()
		for _, info := range infos {
			s.PackageInfo.AllTypeInfos[info.Name] = info
		}
	}
}

// ConstSpec 解析
func (s *source) ConstSpec() []*dao.ConstInfo {
	var datas []string
	var infos []*dao.ConstInfo
	nextIdx := s.nextIdx()

	// top commin
	if s.buf[nextIdx] == '\n' {
		s.next()
		name := s.rangeStr()
		s.toNextCh()

		datas = append(datas, name, "_", "_")

		if s.CheckCommon() {
			s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		}
		return infos
	}

	infos = s.ConstantIdentifierList()

	if s.buf[s.r+1] == '=' { // 隱藏型態 初始化
		s.next()
		exps := s.OnConstantExpression()
		for idx, info := range infos {
			info.Expressions = exps[idx]
		}

	} else { // 一般 初始化

		typeInfo := s.OnDeclarationsType()
		for idx, info := range infos {

			info.TypeInfo = typeInfo

			s.next()
			info.Expressions = s.OnConstantExpression()[idx]
		}
	}

	if s.CheckCommon() {
		common := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		for _, info := range infos {
			info.Common = common
		}
	}
	return infos
}

func (s *source) ConstantIdentifierList() []*dao.ConstInfo {
	var infos []*dao.ConstInfo
	nextIdx := s.nextIdx()
	// 變數名稱解析
	if s.buf[nextIdx-1] == ',' { // 多參數初始化
		for {
			s.next()
			name := s.rangeStr()
			if name[len(name)-1] == ',' {
				info := dao.NewConstInfo(name[:len(name)-1])
				infos = append(infos, info)
				continue

			} else {
				info := dao.NewConstInfo(name)
				infos = append(infos, info)
				break
			}
		}

	} else { // 單一參數初始化

		s.next()
		name := s.rangeStr()
		info := dao.NewConstInfo(name)
		infos = append(infos, info)
	}

	return infos
}

// const a, b, c = Expression, Expression, Expression

// 解析表達式
//
// @param dao.ITypeInfo	表達式型別
// @param int			表達是數量
//
// @return []string 表達式內容
func (s *source) OnConstantExpression() []*dao.Expressions {

	var infos []*dao.Expressions

	for {
		info := dao.NewExpressions()
		info.Objs = s.OnConstantExpressionMath()
		infos = append(infos, info)
		if s.ch != ',' {
			break
		}
		s.toNextCh()
	}

	return infos
}

// 解析 算式陣列
func (s *source) OnConstantExpressionMath() []dao.ITypeInfo {
	var infos []dao.ITypeInfo
	// 暫時不解析隱藏宣告
	baseInfo := dao.BaseTypeInfo[_string]
	infos = append(infos, baseInfo)
	toNext := false
	for {
		s.nextCh()
		switch s.ch {
		case ',', '\n':
			toNext = true
		case ' ':
			if s.buf[s.r+1] == '/' {
				toNext = true
			}
		}
		if toNext {
			break
		}

	}
	return infos

	s.toNextCh()
	ch := rune(s.buf[s.r+1])

	if tool.IsDecimal(ch) || ch == '.' && tool.IsDecimal(rune(s.buf[s.r+2])) {
		info := s.scanNumber()
		infos = append(infos, info)

		// s.toNextCh()
		switch s.ch {
		case ',':
			break

		case ' ':
			s.nextCh()
			switch s.ch {
			case '+':
				// s.nextCh()
				infos = append(infos, s.OnConstantExpressionMath()...)
			}
		case '/':
			infos = append(infos, s.OnConstantExpressionMath()...)
		}

	} else if ch == '"' {

		info := s.scanString()
		infos = append(infos, info)

		switch s.ch {
		case ',':
			break

		case ' ':
			s.nextCh()
			switch s.ch {
			case '+':
				s.nextCh()
				infos = append(infos, s.OnConstantExpressionMath()...)
			}

		}
	}

	return infos
}

//========== Var =================

// VarDecl = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
// VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
func (s *source) VariableDeclarations() {

	if s.buf[s.r+1] == '(' {
		if s.buf[s.r+2] == ')' {
			fmt.Println(":: var () ::")
			return
		}
		s.nextCh()

		for {
			s.toNextCh()
			for _, info := range s.VarSpec() {
				s.PackageInfo.AllTypeInfos[info.Name] = info
			}
			s.toNextCh()
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		for _, info := range s.VarSpec() {
			s.PackageInfo.AllTypeInfos[info.Name] = info
		}
	}
}

// ShortVarDecl = IdentifierList ":=" ExpressionList .
func (s *source) ShortVariableDeclarations() {

}

func (s *source) VarSpec() []*dao.VarInfo {
	var datas []string
	var infos []*dao.VarInfo
	nextIdx := s.nextIdx()

	if s.buf[nextIdx] == '\n' {
		s.next()
		name := s.rangeStr()
		s.toNextCh()

		datas = append(datas, name, "_", "_")

		if s.CheckCommon() {
			s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		}
		return infos
	}

	infos = s.VariableIdentifierList()

	if s.buf[s.r+1] == '=' {
		s.nextCh()

		return infos
	} else {
		typeInfo := s.OnDeclarationsType()
		for idx, info := range infos {

			info.Expressions.Objs = append(info.Expressions.Objs, typeInfo)

			s.next()
			info.Expressions = s.OnVariableExpression()[idx]
		}
	}

	if s.CheckCommon() {
		common := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		for _, info := range infos {
			info.Common = common
		}
	}

	return infos
}

func (s *source) VariableIdentifierList() []*dao.VarInfo {
	var infos []*dao.VarInfo
	nextIdx := s.nextIdx()
	// 變數名稱解析
	if s.buf[nextIdx-1] == ',' { // 多參數初始化
		for {
			s.next()
			name := s.rangeStr()
			if name[len(name)-1] == ',' {
				info := dao.NewVarInfo(name[:len(name)-1])
				infos = append(infos, info)
				continue

			} else {
				info := dao.NewVarInfo(name)
				infos = append(infos, info)
				break
			}
		}

	} else { // 單一參數初始化

		s.next()
		name := s.rangeStr()
		info := dao.NewVarInfo(name)
		infos = append(infos, info)
	}

	return infos
}

// 解析表達式
//
// @param dao.ITypeInfo	表達式型別
// @param int			表達是數量
//
// @return []string 表達式內容
func (s *source) OnVariableExpression() []*dao.Expressions {

	var infos []*dao.Expressions

	for {
		info := dao.NewExpressions()
		info.Objs = s.OnConstantExpressionMath()
		infos = append(infos, info)
		if s.ch != ',' {
			break
		}
		s.toNextCh()
	}

	return infos
}

// 解析 算式陣列
func (s *source) OnVariableExpressionMath() []dao.ITypeInfo {
	var infos []dao.ITypeInfo
	// 暫時不解析隱藏宣告
	baseInfo := dao.BaseTypeInfo[_string]
	infos = append(infos, baseInfo)
	toNext := false
	for {
		s.nextCh()
		switch s.ch {
		case ',', '\n':
			toNext = true
		case ' ':
			if s.buf[s.r+1] == '/' {
				toNext = true
			}
		}
		if toNext {
			break
		}

	}
	return infos

	s.toNextCh()
	ch := rune(s.buf[s.r+1])

	if tool.IsDecimal(ch) || ch == '.' && tool.IsDecimal(rune(s.buf[s.r+2])) {
		info := s.scanNumber()
		infos = append(infos, info)

		// s.toNextCh()
		switch s.ch {
		case ',':
			break

		case ' ':
			s.nextCh()
			switch s.ch {
			case '+':
				// s.nextCh()
				infos = append(infos, s.OnConstantExpressionMath()...)
			}
		case '/':
			infos = append(infos, s.OnConstantExpressionMath()...)
		}

	} else if ch == '"' {

		info := s.scanString()
		infos = append(infos, info)

		switch s.ch {
		case ',':
			break

		case ' ':
			s.nextCh()
			switch s.ch {
			case '+':
				s.nextCh()
				infos = append(infos, s.OnConstantExpressionMath()...)
			}

		}
	}

	return infos
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
func (s *source) TypeDeclarations() {
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
			s.TypeSpec()
			s.toNextCh()

		}
	} else {
		s.TypeSpec()
	}
}

// 處理 Type 宣告
// AliasDecl = identifier "=" Type .
// TypeDef = identifier Type .
// ex: _type t1 string
// r="_"
func (s *source) TypeSpec() (typeInfo dao.ITypeInfo) {

	s.next()
	name := strings.TrimSpace(s.rangeStr())
	isExist := s.PackageInfo.ExistType(name)
	if isExist {
		panic("")
	} else if s.buf[s.r+1] == '=' {
		s.next()
	}
	s.toNextCh()
	typeInfo = s.OnDeclarationsType()
	typeInfo.SetName(name)
	return typeInfo
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
	panic("")
	s.OnParams(true)
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
	panic("")
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
	s.OnParams(true)
	return nil
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
func (s *source) OnParams(isParameters bool) []*dao.FuncParams {
	strLit := []*dao.FuncParams{}
	var nextIdx int = s.r
	var unknowCount int
	unknowData := []string{}

	if s.buf[s.r+1] == ')' {
		s.next()
		return strLit
	}

	for isNext := true; isNext; {
		nextIdx = s.nextIdxByPosition(nextIdx)

		switch s.buf[nextIdx-1] {
		case ',':
			unknowCount++

		case ')':
			panic("")
			// 隱藏名稱宣告
			// for i := 0; i < unknowCount; i++ {
			// 	data := s.OnType(')')
			// 	s.nextCh()
			// 	strLit = append(strLit, []string{"", data})
			// }
			// strLit = append(strLit, []string{"", s.OnType(')')})
			// unknowCount = 0

			// isNext = false

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
			panic("")
			// paramName := strings.TrimSpace(s.rangeStr())
			// paramType := s.OnType(')')

			// for _, name := range unknowData {
			// 	strLit = append(strLit, []string{name, paramType})
			// }
			// strLit = append(strLit, []string{paramName, paramType})
			// unknowData = make([]string, 0)

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
		s.OnParams(true)

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
