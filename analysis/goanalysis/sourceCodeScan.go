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
func (s *source) ImportDeclarations() {
	var packageLinks []*dao.ImportInfo
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
			packageLinks = append(packageLinks, dao.NewImportLink(name, importPath, packageInfo))
			s.toNextCh()

		}
	} else {
		name, importPath, packageInfo := s.importSpec()
		packageLinks = append(packageLinks, dao.NewImportLink(name, importPath, packageInfo))
	}

	for _, link := range packageLinks {
		s.PackageInfo.ImportLink[link.GetName()] = link
	}
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
	var infos []*dao.ConstInfo

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
				constInfos := s.ConstSpec()
				for _, constInfo := range constInfos {
					infos = append(infos, constInfo)
				}
			}

			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		constInfos := s.ConstSpec()
		for _, constInfo := range constInfos {
			infos = append(infos, constInfo)
		}
	}
	for _, info := range infos {
		s.PackageInfo.AllTypeInfos[info.GetName()] = info
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
	s.toNextCh()

	if s.buf[s.r+1] == '=' { // 隱藏型態 初始化
		s.next()
		exps := s.OnConstantExpression()
		for idx, info := range infos {
			info.Expressions = exps[idx]
		}

	} else { // 一般 初始化

		typeInfo := s.OnDeclarationsType()
		s.next()
		exps := s.OnConstantExpression()
		for idx, info := range infos {

			info.TypeInfo = typeInfo
			info.Expressions = exps[idx]
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
	baseInfo := dao.BaseTypeInfo["string"]
	infos = append(infos, baseInfo)
	toNext := false
	for {
		s.nextCh()
		switch s.ch {
		case ',', '\n':
			toNext = true
		case ' ':
			if s.CheckCommon() {
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
				s.PackageInfo.AllTypeInfos[info.GetName()] = info
			}
			s.toNextCh()
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		for _, info := range s.VarSpec() {
			s.PackageInfo.AllTypeInfos[info.GetName()] = info
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
	s.toNextCh()

	if s.buf[s.r+1] == '=' {
		s.next()
		exps := s.OnVariableExpression()
		for idx, info := range infos {
			info.Expressions = exps[idx]
		}

	} else {
		typeInfo := s.OnDeclarationsType()
		var exps []*dao.Expressions
		// 指定初始化
		if s.buf[s.r+1] == '=' {
			s.next()
			// 解析表達式
			exps = s.OnVariableExpression()

		}
		// 默認初始化

		// 建立關聯
		for idx, info := range infos {
			// 指定型別
			info.TypeInfo = typeInfo

			// 關聯 表達式內容
			if len(exps) > idx {
				info.Expressions = exps[idx]
			}
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
		info.Objs = s.OnVariableExpressionMath()
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
	baseInfo := dao.BaseTypeInfo["string"]
	infos = append(infos, baseInfo)
	toNext := false
	for {
		s.nextCh()
		switch s.ch {
		case ',', '\n':
			toNext = true
		case ' ':
			if s.CheckCommon() {
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
	infos := []dao.ITypeInfo{}
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
			info := s.TypeSpec()
			s.PackageInfo.AllTypeInfos[info.GetName()] = info
			infos = append(infos, info)
			s.toNextCh()
		}
	} else {
		info := s.TypeSpec()
		s.PackageInfo.AllTypeInfos[info.GetName()] = info
		infos = append(infos, info)
	}
}

// 處理 Type 宣告
// AliasDecl = identifier "=" Type .
// TypeDef = identifier Type .
// ex: _type t1 string
// r="_"
func (s *source) TypeSpec() dao.ITypeInfo {
	var typeInfo dao.ITypeInfo
	s.next()
	name := strings.TrimSpace(s.rangeStr())
	isExist := s.PackageInfo.ExistType(name)
	if isExist {
		panic("")
	} else if s.buf[s.r+1] == '=' {
		s.next()
		s.toNextCh()

		info := dao.NewTypeAliasDecl()
		info.ContentTypeInfo = s.OnDeclarationsType()
		info.SetName(name)
		typeInfo = info
	} else {
		s.toNextCh()
		info := dao.NewTypeDef()
		info.ContentTypeInfo = s.OnDeclarationsType()
		info.SetName(name)
		typeInfo = info
	}
	return typeInfo
}

//========== Func =================

// FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
// FunctionName = identifier .
// FunctionBody = Block .
func (s *source) FunctionDeclarations() {
	s.nextCh()
	info := dao.NewFuncInfo()

	// FunctionName
	info.SetName(s.scanIdentifiers())

	// Signature
	info.ParamsInPoint = s.OnParameters()
	info.ParamsOutPoint = s.OnDeclarationsResult()

	// FunctionBody
	info.Body = s.funcBodyBlock()
	s.PackageInfo.FuncInfo[info.GetName()] = info
}

// MethodDecl = "func" Receiver MethodName Signature [ FunctionBody ] .
// Receiver   = Parameters .
func (s *source) MethodDeclarations() {
	s.nextCh()

	// Receiver
	info := dao.NewFuncInfo()
	receiver := s.OnParameters()
	if len(receiver) == 0 {
		panic("")
	}
	info.Receiver = receiver[0]

	// MethodName
	s.nextCh()
	info.SetName(s.scanIdentifiers())

	// Signature
	info.ParamsInPoint = s.OnParameters()
	info.ParamsOutPoint = s.OnDeclarationsResult()

	// FunctionBody
	info.Body = s.funcBodyBlock()
	s.PackageInfo.FuncInfo[info.GetName()] = info
}

// func 名稱
func (s *source) OnFuncName() string {
	s.nextToken()
	return strings.TrimSpace(s.rangeStr())
}

// 處理 func 多參數資料
/* 進入指標應當只在 abc(w *int) (x, z int)
 * 輸入參數指標 b=任意 r="("
 * =====================
 * 輸出參數指標 b=任意 r="("
 */
func (s *source) OnParameters() []*dao.FuncParams {
	strLit := []*dao.FuncParams{}
	nextCh := rune(s.buf[s.r+1])

	// 無參數
	if nextCh == ')' {
		s.next()
		return strLit
	}

	// 解析 ParameterList
	istype := false // 匿名參數
	for {

		// 解析 ParameterDecl
		identifiers := []string{}
		if istype {
			// 匿名參數解析
			info := s.OnDeclarationsType()
			params := dao.NewFuncParams()
			params.SetName("_")
			params.ContentTypeInfo = info
			strLit = append(strLit, params)

		} else {
			// 第一次處理與多名稱宣告解析
			parameterDeclDone := false // 多名稱宣告是否完成
			for !parameterDeclDone && !istype {

				// 提前判斷到隱藏宣告
				if !tool.IsLetter(nextCh) && !tool.IsDecimal(nextCh) && s.ch != '_' {
					// 處理未解析完成的資料
					if len(identifiers) > 0 {
						for _, identifier := range identifiers {
							// 基礎類型
							info, ok := dao.BaseTypeInfo[identifier]
							if !ok {
								// 自定義類型
								info = s.PackageInfo.GetType(identifier)
								if info == nil {
									panic("")
								}
							}

							params := dao.NewFuncParams()
							params.SetName("_")
							params.ContentTypeInfo = info
							strLit = append(strLit, params)

						}
					}

					info := s.OnDeclarationsType()
					params := dao.NewFuncParams()
					params.SetName("_")
					params.ContentTypeInfo = info
					strLit = append(strLit, params)
					istype = true
					parameterDeclDone = true
					continue
				}

				s.nextCh()
				identifierOrType := s.scanIdentifiers()

				switch s.ch {
				case '.':
					// Qualified 類型
					s.nextCh()
					typeName := s.scanIdentifiers()
					fullName := fmt.Sprintf("%s.%s", identifierOrType, typeName)
					info := dao.NewTypeInfoQualifiedIdent()
					info.SetName(fullName)
					info.ImportLink, info.ContentTypeInfo = s.PackageInfo.GetPackageType(identifierOrType, typeName)

					params := dao.NewFuncParams()
					params.SetName("_")
					params.ContentTypeInfo = info
					strLit = append(strLit, params)
					istype = true

				default:
					// 多重隱藏宣告 基礎類型
					info, ok := dao.BaseTypeInfo[identifierOrType]
					if ok {
						params := dao.NewFuncParams()
						params.SetName("_")
						params.ContentTypeInfo = info
						strLit = append(strLit, params)
						istype = true
					} else {
						identifiers = append(identifiers, identifierOrType)
					}
				}

				// 結束判斷
				if s.ch != ',' {
					// 多重宣告
					if len(identifiers) > 0 {
						// 多重隱藏宣告 自定義類型
						if s.ch == ')' {
							// (int, float) 類型解析
							info := s.PackageInfo.GetType(identifierOrType)
							if info == nil {
								panic("")
							}
							params := dao.NewFuncParams()
							params.SetName("_")
							params.ContentTypeInfo = info
							strLit = append(strLit, params)

						} else {
							// (a, b Type) 類型解析
							info := s.OnDeclarationsType()
							for _, identifier := range identifiers {
								params := dao.NewFuncParams()
								params.SetName(identifier)
								params.ContentTypeInfo = info
								strLit = append(strLit, params)
							}
						}
					}

					parameterDeclDone = true
				} else {
					s.next()
					nextCh = rune(s.buf[s.r+1])
				}
			}
		}

		if s.ch == ',' {
			s.nextCh()
		} else if s.ch == ')' {
			s.nextCh()
			break
		}

	}

	return strLit
}

// func 輸出參數
func (s *source) OnDeclarationsResult() []*dao.FuncParams {
	params := []*dao.FuncParams{}

	nextCh := s.buf[s.r+1]
	// 判斷無 result
	if s.ch == '\n' || nextCh == '{' {
		return params
	}

	if nextCh == '\t' {
		s.toNextCh()
		nextCh = s.buf[s.r+1]
	}

	if nextCh == '(' {
		s.nextCh()
		params = s.OnParameters()
	} else {
		info := s.OnDeclarationsType()
		param := dao.NewFuncParams()
		param.ContentTypeInfo = info
		params = append(params, param)
	}

	return params
}

// _{
//   return
//  }
// b=任意 r="_"
func (s *source) funcBodyBlock() string {
	if s.buf[s.r+1] != '{' || s.ch == '\n' {
		return ""
	}
	s.nextCh()
	offset := s.r
	scanDone := false
	var blockCount int

	for !scanDone {
		s.toNextToken()
		switch s.buf[s.r+1] {
		case '}':
			s.nextCh()
			if blockCount > 0 {
				blockCount--
				continue
			}
			scanDone = true
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
	s.nextCh()
	return string(s.buf[offset:s.r])
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
