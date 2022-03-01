package goanalysis

import (
	"codeanalysis/analysis/constant"
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

	// 判斷來源
	packageType := constant.From_Golang
	if Instants.ModuleInfo.IsLocalPackage(path) {
		packageType = constant.From_Local

	} else if Instants.ModuleInfo.IsVendorPackage(path) {
		packageType = constant.From_Vendor

	}

	// 解析 預設名稱
	splitStr := strings.Split(path, "/")
	name = splitStr[len(splitStr)-1]

	// 需要再次定義, 其他模式下的關聯名稱
	if importMod != "" || newName == "" {
		// 不明格式可能是 gopls 作用
		if len(name) > 3 && name[:3] == "go-" {
			newName = name[3:]
		} else if len(name) > 3 && name[len(name)-3:] == ".go" {
			newName = name[:len(name)-3]
		} else {
			newName = name
		}
	}

	if name == "" {
		panic("import name error")
	}
	// 根據預設名稱取得 package 關聯資料
	packageInfo, _ = Instants.LoadOrStoryPackage(packageType, path, dao.NewPackageInfo())
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
			return
		}
		s.nextCh()

		for {
			s.toNextCh()
			if s.checkCommon() {
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			} else {
				constInfos := s.ConstSpec()
				if len(constInfos) > 1 {
					panic("empty constInfos error")
				}
				if constInfos[0].ContentTypeInfo == nil && len(infos) > 0 {
					for _, constInfo := range constInfos {
						constInfo.ContentTypeInfo = infos[len(infos)-1].ContentTypeInfo
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

	if !s.isOnNewlineSymbol() { // 判斷區塊未結束
		s.toNextCh()

		// 確認下個區塊類型
		nextCh := rune(s.buf[s.r+1])
		if nextCh == '=' { // 隱藏型態 初始化
			// 解析  (xx = yy), (xx, xx, xx = yy) 格式
			s.next()

			// 解析表達式
			exps := s.scanExpressionList()
			for idx, info := range infos {
				info.Expression = exps[idx]
			}

			if s.ch == ' ' {
				s.toNextCh()
			}

			// 解析註解
			if s.checkCommon() {
				common := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
				for _, info := range infos {
					info.Common = common
				}
			}

		} else if s.checkCommon() {
			// 解析 (xx, xx, xx //)格式
			// 解析註解
			common := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			for _, info := range infos {
				info.Common = common
			}
		} else {
			// 一般初始化流程
			typeInfo := s.OnDeclarationsType()
			s.next()

			exps := s.scanExpressionList()
			for idx, info := range infos {
				info.ContentTypeInfo = typeInfo
				info.Expression = exps[idx]
			}

		}

	} else {
		s.toNextCh()
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
		name := s.scanIdentifier()
		info := dao.NewConstInfo(name)
		infos = append(infos, info)
	}

	return infos
}

//========== Var =================

// VarDecl = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
// VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
func (s *source) VariableDeclarations() {

	if s.buf[s.r+1] == '(' {
		if s.buf[s.r+2] == ')' {
			return
		}
		s.nextCh()

		for {
			s.toNextCh()
			for _, info := range s.VarSpec() {
				if info.GetName() == "_" {
					s.PackageInfo.ImplicitlyVarOrConstInfos = append(s.PackageInfo.ImplicitlyVarOrConstInfos, info)
				} else {
					s.PackageInfo.AllVarInfos[info.GetName()] = info
				}
			}
			s.toNextCh()
			if s.buf[s.r+1] == ')' {
				s.next()
				break
			}
		}
	} else {
		for _, info := range s.VarSpec() {
			if info.GetName() == "_" {
				s.PackageInfo.ImplicitlyVarOrConstInfos = append(s.PackageInfo.ImplicitlyVarOrConstInfos, info)
			} else {
				s.PackageInfo.AllVarInfos[info.GetName()] = info
			}
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
		exps := s.scanExpressionList()
		for idx, info := range infos {
			info.Expression = exps[idx]
		}

	} else {
		typeInfo := s.OnDeclarationsType()
		var exps []string
		if s.checkEnd() {
			// 指定初始化
			if s.buf[s.r+1] == '=' {
				s.next()
				// 解析表達式
				exps = s.scanExpressionList()

			}
		}
		// 默認初始化

		// 建立關聯
		for idx, info := range infos {
			// 指定型別
			info.ContentTypeInfo = typeInfo

			// 關聯 表達式內容
			if len(exps) > idx {
				info.Expression = exps[idx]
			}
		}
	}

	if s.checkCommon() {
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

			if s.checkCommon() {
				s.nextTargetToken('\n')
				continue
			}

			info := s.TypeSpec()
			s.PackageInfo.AllTypeInfos[info.GetName()] = info
			infos = append(infos, info)
			s.toNextCh()
		}
	} else {
		s.toNextCh()
		info := s.TypeSpec()
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

	s.toNextCh()
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
			info = s.OnSliceType()
		} else { // array type
			info = s.OnArrayType()
		}

	case '<': // OutPutChanelType
		info = s.OnChannelType()
	case '.': // short ArranType
		info = s.onShortSliceType()
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
					info = s.OnChannelType()
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
					info = s.OnChannelType()
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
	info.SetName(s.scanIdentifier())

	// Signature
	info.ParamsInPoint, info.ParamsOutPoint = s.onSignature()

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
	name := s.scanIdentifier()
	info.SetName(name)

	// Signature
	info.ParamsInPoint, info.ParamsOutPoint = s.onSignature()

	// FunctionBody
	info.Body = s.funcBodyBlock()
	s.PackageInfo.AllFuncInfo[info.GetName()] = info
}

// func 名稱
func (s *source) onSignature() (paramsInPoint, paramsOutPoint []dao.FuncParams) {
	paramsInPoint = s.OnParameters()

	// 遇到換行符號提早結束
	if s.isOnNewlineSymbol() {
		return
	}

	// result 前一定會有空白
	if s.ch != ' ' {
		return
	}
	// 調整格式
	// 用於判斷 func()		{} 格式
	nextCh := s.buf[s.r+1]
	if nextCh == ' ' || nextCh == '\t' {
		s.toNextCh()
	}

	if !s.checkCommon() {
		paramsOutPoint = s.OnDeclarationsResult()
	}
	return
}

// 處理 func 多參數資料
/* 進入指標應當只在 abc(w *int) (x, z int)
 * 輸入參數指標 b=任意 r="("
 * =====================
 * 輸出參數指標 b=任意 r="("
 */
func (s *source) OnParameters() []dao.FuncParams {
	strLit := []dao.FuncParams{}
	// nextCh := rune(s.buf[s.r+1])
	isCommonLine := false     // 前一行是否為註解
	lineCommon := ""          // 註解內文
	istype := false           // 匿名參數
	identifiers := []string{} // 紀錄存文字 identifier *token 可判斷為類型所以不用儲存

	// 確認 identifiers 資料型態後呼叫
	clearIdentifiers := func(info dao.ITypeInfo) {
		if info == nil {
			// 格式為匿名宣告
			for _, identifier := range identifiers {
				if info, ok := dao.BaseTypeInfo[identifier]; ok {
					params := dao.NewFuncParams()
					params.SetName("_")
					params.ContentTypeInfo = info
					strLit = append(strLit, params)
				} else {
					info := s.PackageInfo.GetType(identifier)
					params := dao.NewFuncParams()
					params.SetName("_")
					params.ContentTypeInfo = info
					strLit = append(strLit, params)
				}
			}
			identifiers = []string{}
		} else {
			// 格式為多重宣告
			for _, identifier := range identifiers {
				params := dao.NewFuncParams()
				params.SetName(identifier)
				params.ContentTypeInfo = info
				strLit = append(strLit, params)
			}
			identifiers = []string{}
		}
	}

	// 解析 ParameterList
	for {

		// 換行處理
		if s.isOnNewlineSymbol() {
			s.nextCh()

			// 註解重製判斷
			if !isCommonLine && lineCommon != "" {
				lineCommon = ""
			}

			// 每行需要重製的資料
			isCommonLine = false
			continue
		} else if s.ch == ',' {
			s.nextCh()
			continue
		} else if s.ch == ')' {
			clearIdentifiers(nil)
			s.nextCh()
			break
		}

		// 排除空格
		if s.ch == ' ' || s.ch == '\t' {
			s.toNextCh()
		}

		// 判斷整行註解
		if s.checkCommon() {
			lineCommon += s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			isCommonLine = true
			continue
		} else {
			isCommonLine = false
		}

		// ParameterDecl
		nextCh := rune(s.buf[s.r+1])

		// 無參數 或 格式條整
		if nextCh == ')' || nextCh == '\n' {
			s.nextCh()
			continue
		}

		// 解析 匿名參數
		// 一個匿名全部都匿名
		if istype || util.IsToken(nextCh) {
			info := s.OnDeclarationsType()
			params := dao.NewFuncParams()
			params.SetName("_")
			params.ContentTypeInfo = info
			istype = true
			clearIdentifiers(nil)
			strLit = append(strLit, params)
			continue
		}

		// 還無法判斷 name 或非 token 開頭的 type
		s.nextCh()
		nameOrType := s.scanIdentifier()

		// 判斷區塊結束, 資料為 type
		if s.ch == ',' || s.ch == ')' {
			// 匿名參數 or 多重參數
			identifiers = append(identifiers, nameOrType)

		} else {
			if util.IsToken(s.ch) {
				if s.ch == '.' {
					// Qualified 類型
					s.nextCh()
					typeName := s.scanIdentifier()
					fullName := fmt.Sprintf("%s.%s", nameOrType, typeName)
					info := dao.NewTypeInfoQualifiedIdent()
					info.SetName(fullName)
					info.ImportLink, info.ContentTypeInfo = s.PackageInfo.GetPackageType(nameOrType, typeName)

					params := dao.NewFuncParams()
					params.SetName("_")
					params.ContentTypeInfo = info
					istype = true
					strLit = append(strLit, params)
				} else {
					info := s.OnTypeSwitch(nameOrType)
					params := dao.NewFuncParams()
					params.SetName("_")
					params.ContentTypeInfo = info
					istype = true
					clearIdentifiers(nil)
					strLit = append(strLit, params)
				}
			} else if s.ch == ' ' {
				if nameOrType == "chan" {
					info := s.OnTypeSwitch(nameOrType)
					params := dao.NewFuncParams()
					params.SetName("_")
					params.ContentTypeInfo = info
					istype = true
					clearIdentifiers(nil)
					strLit = append(strLit, params)
				} else {
					// 標準參數格式
					info := s.OnDeclarationsType()
					params := dao.NewFuncParams()
					params.SetName(nameOrType)
					params.ContentTypeInfo = info
					clearIdentifiers(info)
					strLit = append(strLit, params)
				}
			}
		}
	}
	return strLit
}

// func 輸出參數
func (s *source) OnDeclarationsResult() []dao.FuncParams {
	params := []dao.FuncParams{}
	nextCh := s.buf[s.r+1]

	// 判斷無 result
	// ',' 用於判斷 func() (func(), error) 格式
	// '{' 用於判斷 func() {} 格式
	// ')' 用於判斷 func(func()) () 兩個 '))' 格式
	if s.ch == '\n' || s.ch == ',' || s.ch == ')' || nextCh == '{' {
		return params
	}

	if nextCh == '\t' || nextCh == ' ' {
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
