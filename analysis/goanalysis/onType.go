package goanalysis

import (
	"codeanalysis/analysis/dao"
	"fmt"
)

// // StructType    = "struct" "{" { FieldDecl ";" } "}" .
// // FieldDecl     = (IdentifierList Type | EmbeddedField) [ Tag ] .
// // EmbeddedField = [ "*" ] TypeName .
// // Tag           = string_lit .
// /* 進入指標應當只在
//  * r="_" range="struct" // 應該調整為統一在
//  */
// func (s *source) OnStructType() *dao.TypeInfoStruct {
// 	info := dao.NewTypeInfoStruct()
// 	info.SetTypeName("struct")
// 	s.nextToken()
// 	if s.buf[s.r+1] == '}' {
// 		s.nextCh()
// 		s.nextCh()
// 		return info
// 	}
// 	s.nextCh()
// 	// TODO
// 	// 會出現兩個不同結尾 須調整
// 	// FieldDecl
// 	s.OnFieldDecl(info)
// 	if s.ch == '}' {
// 		panic("struct error")
// 	}
// 	return info
// }

func (s *source) OnFieldDecl(info *dao.TypeInfoStruct) {

	var lineCommon string
	isCommonLine := false

	// 從 struct{'\n' <- 這邊開始解析
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
		} else if s.ch == '}' {
			s.nextCh()
			break
		}

		s.toNextCh()

		// 判斷整行註解
		if s.checkCommon() {
			lineCommon += s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			isCommonLine = true
			continue
		} else {
			isCommonLine = false
		}

		var tmpVarInfos []*dao.VarInfo // 此區塊宣告的參數
		nextCh := s.buf[s.r+1]
		if nextCh == '}' {
			s.nextCh()
			s.nextCh()
			break
		}

		if nextCh == '*' {
			// 解析隱藏參數 (struct{ *TypeName}) 格式
			varInfo := s.scanEmbeddedField()
			info.ImplicitlyVarInfos = append(info.ImplicitlyVarInfos, varInfo)
			tmpVarInfos = append(tmpVarInfos, varInfo)
		} else {
			s.nextCh()

			names := s.scanIdentifiers()
			// 調整 (struct{ xxx\n}), (struct{ xxx\t\tType}) 格式
			if !s.isOnNewlineSymbol() {
				s.toNextCh()
			}

			if len(names) == 0 {
				// 解析錯誤
				panic("")

			} else if len(names) == 1 {
				if s.ch == ' ' {
					// 單一名稱
					// 解析 IdentifierList Type 格式

					nextCh = s.buf[s.r+1]
					if nextCh == '`' {
						// 解析 (XXX    `json:""`) 格式
						// 無指標隱藏宣告 後續接到 Tag
						varInfo := dao.NewVarInfo("_")
						varInfo.ContentTypeInfo = s.PackageInfo.GetType(names[0])
						info.ImplicitlyVarInfos = append(info.ImplicitlyVarInfos, varInfo)
						tmpVarInfos = append(tmpVarInfos, varInfo)

					} else if s.checkCommon() {
						// 解析 (XXX    //) 格式
						// 無指標隱藏宣告 後續接到註解
						varInfo := dao.NewVarInfo("_")
						varInfo.ContentTypeInfo = s.PackageInfo.GetType(names[0])
						info.ImplicitlyVarInfos = append(info.ImplicitlyVarInfos, varInfo)
						tmpVarInfos = append(tmpVarInfos, varInfo)

					} else {
						// 正常宣告
						varInfo := dao.NewVarInfo(names[0])
						varInfo.ContentTypeInfo = s.OnDeclarationsType()
						info.VarInfos[varInfo.GetName()] = varInfo
						tmpVarInfos = append(tmpVarInfos, varInfo)
					}

				} else if s.ch == '.' {
					// 解析 EmbeddedField = QualifiedIdent 格式
					// 無指標隱藏宣告
					s.nextCh()
					typeName := s.scanIdentifier()
					fullName := fmt.Sprintf("%s.%s", names[0], typeName)
					quaInfo := dao.NewTypeInfoQualifiedIdent()
					quaInfo.SetName(fullName)
					quaInfo.ImportLink, quaInfo.ContentTypeInfo = s.PackageInfo.GetPackageType(names[0], typeName)

					varInfo := dao.NewVarInfo("_")
					varInfo.ContentTypeInfo = quaInfo
					info.ImplicitlyVarInfos = append(info.ImplicitlyVarInfos, varInfo)
					tmpVarInfos = append(tmpVarInfos, varInfo)
				} else if s.ch == '\n' {
					// 解析 (XXX\n) 格式
					// 無指標隱藏宣告 後續換行
					varInfo := dao.NewVarInfo("_")
					varInfo.ContentTypeInfo = s.PackageInfo.GetType(names[0])
					info.ImplicitlyVarInfos = append(info.ImplicitlyVarInfos, varInfo)
					tmpVarInfos = append(tmpVarInfos, varInfo)

				}
			} else {
				// 解析 EmbeddedField = identifier, identifier type 格式
				contextInfo := s.OnDeclarationsType()
				for _, name := range names {
					varInfo := dao.NewVarInfo(name)
					varInfo.ContentTypeInfo = contextInfo
					info.VarInfos[varInfo.GetName()] = varInfo
					tmpVarInfos = append(tmpVarInfos, varInfo)
				}
			}
		}

		// 處理不同的結束格式
		// 解析 Tag, Common
		if !s.isOnNewlineSymbol() {
			s.toNextCh()

			// 解析 Tag 格式
			if s.buf[s.r+1] == '`' {
				if len(tmpVarInfos) != 1 {
					panic("")
				}

				s.nextCh()
				tag := s.scanStringLit('`')
				tmpVarInfos[0].Tag = tag
				// tmpStrLit = append(tmpStrLit, s.OnJsonTag())
				if !s.isOnNewlineSymbol() {
					s.toNextCh()
				}
			}

			// 解析後注解
			if s.checkCommon() {
				backCommon := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
				if lineCommon == "" {
					lineCommon = backCommon
				}

			}

			// 註解回填到指定變數說明上
			if lineCommon != "" && len(tmpVarInfos) > 0 {
				tmpVarInfos[0].Common = lineCommon
				lineCommon = ""
			}

			if !s.isOnNewlineSymbol() {
				// 單行宣告 * struct{ xxx, yyy int }
				if s.buf[s.r+1] == '}' {
					s.next()
					if !s.isOnNewlineSymbol() {
						s.toNextCh()
						break
					} else {
						break
					}
				} else {
					panic("struct 解析錯誤")
				}
			}
		}
	}
}

// func (s *source) onPointType() *dao.TypeInfoPointer {
// 	s.nextCh()

// 	info := dao.NewTypeInfoPointer()
// 	info.SetTypeName("point")
// 	info.ContentTypeInfo = s.OnDeclarationsType()
// 	return info
// }

// 處理 package type 類型
// QualifiedIdent = PackageName "." identifier .
/* 進入指標應當只在
 * package.ide
 * b="p" r="." range="package"
 */
func (s *source) OnQualifiedIdentType() *dao.TypeInfoQualifiedIdent {
	var info = dao.NewTypeInfoQualifiedIdent()

	packageName := s.rangeStr()
	s.onQualifiedIdentifier(packageName, info)
	return info
}

func (s *source) onQualifiedIdentifier(packageName string, info *dao.TypeInfoQualifiedIdent) {
	s.nextCh()
	typeName := s.scanIdentifier()
	fullName := fmt.Sprintf("%s.%s", packageName, typeName)

	info.SetName(fullName)
	link, iItype := s.PackageInfo.GetPackageType(packageName, typeName)
	info.ImportLink, info.ContentTypeInfo = link, iItype
}

func (s *source) OnSliceType() *dao.TypeInfoSlice {
	info := dao.NewTypeInfoSlice()
	info.SetTypeName("slice")
	s.nextCh()
	s.nextCh()
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

func (s *source) onShortSliceType() *dao.TypeInfoSlice {
	info := dao.NewTypeInfoSlice()
	info.SetTypeName("slice")
	if string(s.buf[s.r+1:s.r+4]) == "..." {
		s.nextCh()
		s.nextCh()
		s.nextCh()
		info.ContentTypeInfo = s.OnDeclarationsType()
	} else {
		panic("")
	}
	return info
}

// ArrayType   = "[" ArrayLength "]" ElementType .
// ArrayLength = Expression .
// ElementType = Type .

// 解析 Array 類型
func (s *source) OnArrayType() *dao.TypeInfoArray {
	s.nextCh()

	info := dao.NewTypeInfoArray()
	info.SetTypeName("array")

	// ArrayLength
	s.nextCh()

	info.Size = s.onFirstScanExpression(']')
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

func (s *source) OnChannelType() *dao.TypeInfoChannel {
	info := dao.NewTypeInfoChannel()
	info.SetTypeName("chan")
	if s.buf[s.r+1] == '<' { // 單出
		s.nextCh()
		s.nextCh()
		info.FlowType = 1

	}

	s.nextCh()
	info.SetName(s.scanIdentifier())
	if s.ch == '<' {
		s.nextCh()
		s.nextCh()
		info.FlowType = 2
	}

	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

// 處理 map 類型
// 進入指標應當只在 _map[]
// r="_"
func (s *source) OnMapType() dao.ITypeInfo {
	info := dao.NewTypeInfoMap()
	info.SetTypeName("map")
	s.nextToken()
	info.KeyType = s.OnDeclarationsType()
	if s.ch != ']' {
		s.nextCh()
	}
	info.ValueType = s.OnDeclarationsType()
	return info
}

// InterfaceType      = "interface" "{" { ( MethodSpec | InterfaceTypeName ) ";" } "}" .
// MethodSpec         = MethodName Signature .
// MethodName         = identifier .
// InterfaceTypeName  = TypeName .
/* 進入指標應當只在
 * interface{
 * b=i r={ rang="interface"
 */

// 解析 interface 類型
func (s *source) OnInterfaceType() *dao.TypeInfoInterface {
	s.nextCh()

	s.scanIdentifier()
	info := dao.NewTypeInfoInterface()
	info.SetTypeName("interface")
	info.MatchInfos = s.onMethodSpec()

	return info
}
func (s *source) onMethodSpec() (matchInfos []dao.ITypeInfo) {
	if s.ch == '{' {
		s.nextCh()
		s.nextCh()
	} else {
		// MethodSpec | InterfaceTypeName
		s.next()
		s.toNextCh()

		// 解析內文
		for {
			if s.buf[s.r+1] == '}' {
				break
			}

			// 解析註解
			if s.checkCommon() {
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
				s.toNextCh()
				continue
			}

			// 解析 interface 方法
			s.nextCh()
			name := s.scanIdentifier()
			switch s.ch {
			case '(':
				// 解析 MathodSpec
				matchInfo := dao.NewFuncInfo()
				matchInfo.SetName(name)
				matchInfo.ParamsInPoint, matchInfo.ParamsOutPoint = s.onSignature()
				matchInfos = append(matchInfos, matchInfo)

			case '\n':
				// 解析 InterfaceTypeName
				iInfo := s.PackageInfo.GetType(name)
				matchInfos = append(matchInfos, iInfo)
			}

			s.toNextCh()
		}
		s.nextCh()
	}
	return
}

// ========== Func 類別宣告  ============

// FunctionType 格式
// FunctionType   = "func" Signature .
// Signature      = Parameters [ Result ] .
// Result         = Parameters | Type .
// Parameters     = "(" [ ParameterList [ "," ] ] ")" .
// ParameterList  = ParameterDecl { "," ParameterDecl } .
// ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
// 處理 func 類型
// 進入指標應當只在 _func(
// r="_"

// 解析 func 類型 *不包含方法
func (s *source) OnFuncType() *dao.TypeInfoFunction {
	info := dao.NewTypeInfoFunction()
	info.SetTypeName("func")
	s.nextToken()

	info.ParamsInPoint, info.ParamsOutPoint = s.onSignature()

	return info
}

// 處理 `_`json:"data"`` reflection interface 區塊
// https://golang.org/ref/spec#Tag
// https://golang.org/pkg/reflect/#StructTag
// 進入 r="_"
// 離開 r="`" 結束符號上
func (s *source) OnJsonTag() (str string) {
	s.nextCh()
	s.nextTargetToken('`')
	str = s.rangeStr()
	return
}

// 處理 `_//` 形式註解區塊
// r="_"
func (s *source) OnTag() (str string) {
	str = s.OnComments("//")
	return
}

// 處理 `_//` 與 `_/* */` 形式註解區塊
// r= "_"
func (s *source) OnComments(commentType string) (str string) {
	if commentType == "//" {
		s.nextTargetToken('\n')
		str = s.rangeStr()
	} else if commentType == "/*" {
		for {
			s.nextTargetToken('*')
			str += s.rangeStr()
			if s.buf[s.r+1] == '/' {
				s.nextCh()
				break
			}
		}
	} else {
		panic("OnComments Error")
	}
	return
}

// type switch ===============================

func (s *source) OnTypeSwitch(key string) (info dao.ITypeInfo) {
	switch key {
	case "*", "&":
		pointInfo := dao.NewTypeInfoPointer()
		pointInfo.SetTypeName("point")
		pointInfo.ContentTypeInfo = s.OnDeclarationsType()
		info = pointInfo
	case "struct":
		structInfo := dao.NewTypeInfoStruct()
		structInfo.SetTypeName("struct")
		if s.buf[s.r+1] == '}' {
			s.nextCh()
			s.nextCh()
		} else {
			s.nextCh()
			s.OnFieldDecl(structInfo)
		}
		info = structInfo

	case "[":
		if s.buf[s.r+1] == ']' {
			// slice
			s.nextCh()
			sliceInfo := dao.NewTypeInfoSlice()
			sliceInfo.SetTypeName("slice")
			sliceInfo.ContentTypeInfo = s.OnDeclarationsType()
			info = sliceInfo
		} else {
			// array
			s.nextCh()
			arrayInfo := dao.NewTypeInfoArray()
			arrayInfo.SetTypeName("array")
			arrayInfo.Size = s.onFirstScanExpression(']')
			arrayInfo.ContentTypeInfo = s.OnDeclarationsType()
			info = arrayInfo
		}
	case "chan":
		chanInfo := dao.NewTypeInfoChannel()
		chanInfo.SetName("chan")
		if s.ch == '<' {
			s.nextCh()
			s.nextCh()
			chanInfo.FlowType = 2
		}

		chanInfo.ContentTypeInfo = s.OnDeclarationsType()
		info = chanInfo
	case "func":
		funcInfo := dao.NewTypeInfoFunction()
		funcInfo.SetTypeName("func")
		funcInfo.ParamsInPoint, funcInfo.ParamsOutPoint = s.onSignature()
		info = funcInfo
	case "interface":
		interfaceInfo := dao.NewTypeInfoInterface()
		interfaceInfo.SetTypeName("interface")
		interfaceInfo.MatchInfos = s.onMethodSpec()
		info = interfaceInfo

	case "map":
		mapInfo := dao.NewTypeInfoMap()
		mapInfo.SetTypeName("map")
		mapInfo.KeyType = s.OnDeclarationsType()
		if s.ch != ']' {
			s.nextCh()
		}
		mapInfo.ValueType = s.OnDeclarationsType()
		info = mapInfo
	}
	return
}
