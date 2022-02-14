package goanalysis

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/util"

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
			packageLinks = append(packageLinks, s.importSpec())
			s.toNextCh()

		}
	} else {
		packageLinks = append(packageLinks, s.importSpec())
	}

	for _, link := range packageLinks {
		s.PackageInfo.AllImportLink[link.Path] = link
	}
}

// 解析 import 內文
//
// @return name 		import檔案的替換名稱
// @return importPath	import檔案的路徑
func (s *source) importSpec() *dao.ImportInfo {
	s.nextToken()

	var importMod string
	var newName, name, path string
	var packageInfo *dao.PackageInfo

	// 解析 新名稱
	newName = strings.TrimSpace(s.rangeStr())
	s.nextTargetToken('"')
	if newName == "." || newName == "_" {
		importMod = newName
	}

	// 解析 路徑
	path = strings.TrimSpace(s.rangeStr())

	// 解析 預設名稱
	splitStr := strings.Split(path, "/")
	name = splitStr[len(splitStr)-1]

	// 需要再次定義, 其他模式下的關聯名稱
	if importMod != "" || newName == "" {
		// 不明格式可能是 gopls 作用
		if len(name) > 3 && name[:3] == "go-" {
			newName = name[3:]
		} else {
			newName = name
		}
	}

	if name == "" {
		fmt.Println("")
	}
	// 根據預設名稱取得 package 關聯資料
	packageInfo, _ = Instants.LoadOrStoryPackage(path, dao.NewPackageInfo())
	importInfo := dao.NewImportLink()
	importInfo.NewName = newName
	importInfo.ImportMod = importMod
	importInfo.Path = path
	importInfo.Package = packageInfo
	return importInfo
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
				if len(constInfos) > 1 {
					panic("")
				}
				if constInfos[0].TypeInfo == nil && len(infos) > 0 {
					for _, constInfo := range constInfos {
						constInfo.TypeInfo = infos[len(infos)-1].TypeInfo
					}
				}
				infos = append(infos, constInfos...)
			}

			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		constInfos := s.ConstSpec()
		infos = append(infos, constInfos...)
	}
	for _, info := range infos {
		s.PackageInfo.AllConstInfos[info.GetName()] = info
	}
}

// ConstSpec 解析
func (s *source) ConstSpec() []*dao.ConstInfo {
	infos := s.ConstantIdentifierList()

	if s.ch != '\n' { // 判斷區塊未結束
		s.toNextCh()

		// 確認下個區塊類型
		nextCh := rune(s.buf[s.r+1])
		if nextCh == '=' { // 隱藏型態 初始化
			s.next()

			// 解析表達式
			exps := s.OnConstantExpression()
			for idx, info := range infos {
				info.Expressions = exps[idx]
			}

			// 解析註解
			if s.CheckCommon() {
				common := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
				for _, info := range infos {
					info.Common = common
				}
			}

		} else if s.CheckCommon() {
			// 解析註解
			common := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			for _, info := range infos {
				info.Common = common
			}
		} else {
			// 一般初始化流程
			typeInfo := s.OnDeclarationsType()
			s.next()
			exps := s.OnConstantExpression()
			for idx, info := range infos {

				info.TypeInfo = typeInfo
				info.Expressions = exps[idx]
			}
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

		s.nextCh()
		name := s.scanIdentifiers()
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
	// 後續補完此解析
	return infos

	// 解析流程
	s.toNextCh()
	ch := rune(s.buf[s.r+1])

	if util.IsDecimal(ch) || ch == '.' && util.IsDecimal(rune(s.buf[s.r+2])) {
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
				s.PackageInfo.AllVarInfos[info.GetName()] = info
			}
			s.toNextCh()
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		for _, info := range s.VarSpec() {
			s.PackageInfo.AllVarInfos[info.GetName()] = info
		}
	}
}

// ShortVarDecl = IdentifierList ":=" ExpressionList .
func (s *source) ShortVariableDeclarations() {

}

func (s *source) VarSpec() []*dao.VarInfo {
	var infos []*dao.VarInfo
	nextIdx := s.nextIdx()

	if s.buf[nextIdx] == '\n' {
		panic("")
		// s.next()
		// name := s.rangeStr()
		// s.toNextCh()

		// if s.CheckCommon() {
		// 	s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		// }
		// return infos
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
	// 後續補完解析
	return infos

	// 解析流程
	s.toNextCh()
	ch := rune(s.buf[s.r+1])

	if util.IsDecimal(ch) || ch == '.' && util.IsDecimal(rune(s.buf[s.r+2])) {
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
// @return []dao.ITypeInfo 宣告內容
func (s *source) TypeDeclarations() []dao.ITypeInfo {
	infos := []dao.ITypeInfo{}
	if s.buf[s.r+1] == '(' {
		s.nextCh()
		s.toNextCh()
		for {
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}

			if s.CheckCommon() {
				s.nextTargetToken('\n')
				continue
			}

			info := s.TypeSpec()
			if info.GetName() == "" {
				fmt.Println("")
			}
			s.PackageInfo.AllTypeInfos[info.GetName()] = info
			infos = append(infos, info)
			s.toNextCh()
		}
	} else {
		s.toNextCh()
		info := s.TypeSpec()
		if info.GetName() == "" {
			fmt.Println("")
		}
		s.PackageInfo.AllTypeInfos[info.GetName()] = info
		infos = append(infos, info)
	}
	return infos
}

// 處理 Type 宣告
// AliasDecl = identifier "=" Type .
// TypeDef = identifier Type .
// ex: _type t1 string
// r="_"
func (s *source) TypeSpec() dao.ITypeInfo {
	var typeInfo *dao.TypeInfo
	s.next()
	name := strings.TrimSpace(s.rangeStr())
	if name == "ITheme" {
		fmt.Println("")
	}
	if s.buf[s.r+1] == '=' {
		s.next()
		s.toNextCh()

		info := dao.NewTypeAliasDecl()
		info.SetName(name)
		info.ContentTypeInfo = s.OnDeclarationsType()
		typeInfo = info
	} else {
		s.toNextCh()
		info := dao.NewTypeDef()
		if name == "Gc_OpenByEcSiteReq" {
			fmt.Println("")
		}
		info.SetName(name)
		info.ContentTypeInfo = s.OnDeclarationsType()
		typeInfo = info
	}

	isExist := s.PackageInfo.ExistType(typeInfo.GetName())
	if isExist {
		info := s.PackageInfo.GetType(name).(*dao.TypeInfo)
		switch typeInfo.DefType {
		case "Decl":
			info.DefType = typeInfo.DefType
			info.ContentTypeInfo = typeInfo.ContentTypeInfo

		case "Def":
			info.DefType = typeInfo.DefType
			info.ContentTypeInfo = typeInfo.ContentTypeInfo
		}

	}
	return typeInfo
}

// 處理 type 參數類型
/* 進入指標應當只在 _types
 * b=任意 r="_"
 */
func (s *source) OnDeclarationsType() (info dao.ITypeInfo) {

	// s.nextCh()
	switch s.buf[s.r+1] {
	case '*': // PointerType
		info = s.onPointType()
	case '[': // ArrayType, SliceType
		if s.buf[s.r+2] == ']' { // slice type
			info = s.OnSliceType('\n')
		} else { // array type
			info = s.OnArrayType()
		}

	case '<': // OutPutChanelType
		s.onChannelType()
	case '.': // short ArranType
		info = s.onShortArrayType()
	default:
		nextEndIdx := s.nextIdx()
		nextTokenIdx := s.nextTokenIdx()

		if nextEndIdx < nextTokenIdx {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextEndIdx]))

			baseInfo, ok := dao.BaseTypeInfo[tmpStr]
			if ok {
				s.next()
				info = baseInfo
			} else {
				if tmpStr == "struct" {
					info = s.OnStructType()
				} else if tmpStr == "chan" {
					info = s.onChannelType()
				} else if tmpStr == "interface" {
					info = s.OnInterfaceType()
				} else {
					info = s.PackageInfo.GetType(tmpStr)
					s.next()
				}
			}

		} else if nextTokenIdx == nextEndIdx {
			s.nextToken()
		} else {

			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))

			baseInfo, ok := dao.BaseTypeInfo[tmpStr]
			if ok {
				s.nextToken()
				info = baseInfo
			} else {

				switch tmpStr {
				case "func":
					info = s.OnFuncType()
				case "map":
					info = s.OnMapType()
				case "interface":
					info = s.OnInterfaceType()
				case "struct":
					info = s.OnStructType()
				case "chan":
					info = s.onChannelType()
				default:
					if s.buf[nextTokenIdx] == '.' {
						info = s.OnQualifiedIdentType()

					} else {
						info = s.PackageInfo.GetType(tmpStr)
						s.nextToken()
					}

				}
			}
		}
	}
	return
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
	s.PackageInfo.AllFuncInfo[info.GetName()] = info
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
	s.PackageInfo.AllFuncInfo[info.GetName()] = info
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
func (s *source) OnParameters() []dao.FuncParams {
	strLit := []dao.FuncParams{}
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
				if !util.IsLetter(nextCh) && !util.IsDecimal(nextCh) && s.ch != '_' {
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
func (s *source) OnDeclarationsResult() []dao.FuncParams {
	params := []dao.FuncParams{}

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
		param.SetName("_")
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
			} else {
				s.nextCh()
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
