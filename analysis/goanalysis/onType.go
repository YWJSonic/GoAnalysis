package goanalysis

import (
	"codeanalysis/analysis/dao"
	"fmt"
	"strings"
)

// StructType    = "struct" "{" { FieldDecl ";" } "}" .
// FieldDecl     = (IdentifierList Type | EmbeddedField) [ Tag ] .
// EmbeddedField = [ "*" ] TypeName .
// Tag           = string_lit .
/* 進入指標應當只在
 * r="_" range="struct" // 應該調整為統一在
 */
func (s *source) OnStructType() *dao.TypeInfoStruct {
	info := dao.NewTypeInfoStruct()
	info.SetTypeName("struct")
	s.nextToken()
	if s.buf[s.r+1] == '}' {
		s.next()
		return info
	}
	s.nextCh()

	// FieldDecl
	s.OnFieldDecl(info)

	return info
}

func (s *source) OnFieldDecl(info *dao.TypeInfoStruct) {
	unknowData := []string{}
	for {
		s.toNextCh()

		// 判斷整行註解
		if s.CheckCommon() {
			s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		} else {

			if s.buf[s.r+1] == '}' {
				s.nextTargetToken('}')
				break
			}

			var name string
			var tmpVarInfos []*dao.VarInfo
			nextIdx := s.nextIdx()

			// 解析變數名稱
			if !(s.buf[nextIdx] == '\n') {
				s.next()
				name = strings.TrimSpace(s.rangeStr())

				// 宣告語法糖
				if s.buf[s.r-1] == ',' {
					unknowData = append(unknowData, name[:len(name)-1])
					continue
				}
			}

			// 清除格式空格
			s.toNextCh()

			// 解析變數型態
			contextInfo := s.OnDeclarationsType()
			if len(unknowData) > 0 {
				for _, nameData := range unknowData {
					varInfo := dao.NewVarInfo(nameData)
					varInfo.TypeInfo = contextInfo
					info.VarInfos[nameData] = varInfo
					tmpVarInfos = append(tmpVarInfos, varInfo)
				}
			} else {
				varInfo := dao.NewVarInfo(name)
				varInfo.TypeInfo = contextInfo
				info.VarInfos[name] = varInfo
				tmpVarInfos = append(tmpVarInfos, varInfo)

			}
			unknowData = make([]string, 0)
			s.toNextCh()
			for s.ch == ' ' {
				if string(s.buf[s.r+1:s.r+3]) == "//" {
					common := s.OnComments("//")
					for _, varInfo := range tmpVarInfos {
						varInfo.Common = common
					}
				} else if s.buf[s.r+1] == '`' {
					if len(tmpVarInfos) != 1 {
						panic("")
					}
					s.next()
					tag := s.rangeStr()
					tmpVarInfos[0].Tag = tag
					// tmpStrLit = append(tmpStrLit, s.OnJsonTag())
					s.toNextCh()
				}
			}
		}

		// if s.buf[s.r+1] == '}' {
		// 	s.nextTargetToken('}')
		// 	break
		// }
	}
}

func (s *source) onPointType() *dao.TypeInfoPointer {
	s.nextCh()

	info := dao.NewTypeInfoPointer()
	info.SetTypeName("*")
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

// 處理 package type 類型
// QualifiedIdent = PackageName "." identifier .
/* 進入指標應當只在
 * package.ide
 * b="p" r="." range="package"
 */
func (s *source) OnQualifiedIdentType() *dao.TypeInfoQualifiedIdent {
	var info = dao.NewTypeInfoQualifiedIdent()

	s.nextToken()
	packageName := s.rangeStr()
	s.nextCh()
	typeName := s.scanIdentifiers()
	fullName := fmt.Sprintf("%s.%s", packageName, typeName)

	info.SetName(fullName)
	link, iItype := s.PackageInfo.GetPackageType(packageName, typeName)
	info.ImportLink, info.ContentTypeInfo = link, iItype
	return info
}

func (s *source) OnSliceType(endTag byte) *dao.TypeInfoSlice {
	info := dao.NewTypeInfoSlice()
	info.SetTypeName("slice")
	s.nextCh()
	s.nextCh()
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

func (s *source) onShortArrayType() *dao.TypeInfoSlice {
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
	info.Size = s.scanExpression()
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

func (s *source) onChannelType() *dao.TypeInfoChannel {

	info := dao.NewTypeInfoChannel()
	if s.buf[s.r+1] == '<' { // 單出
		s.next()
		info.SetName(s.rangeStr())
		info.FlowType = 1

	} else {
		s.next()

		info.SetName(s.rangeStr())

		// channel 方向
		if info.GetName() == "chan" { // 雙向
			info.FlowType = 0

		} else { // 單入
			info.FlowType = 2

		}
	}

	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

// 處理 channel 類型
// ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
/* 進入指標應當只在
 * _chan
 * r="_"
 * _chan<-
 * r="_"
 * _<-chan
 * r="_"
 */
func (s *source) OnChannelType() string {
	str := ""
	if s.buf[s.r+1] == '<' { // 單出
		s.next()
		str = s.rangeStr()

	} else {
		s.next()
		str = s.rangeStr()
	}

	// str = str + " " + s.OnType()
	return str
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

	typeName := s.scanIdentifiers()
	info := dao.NewTypeInfoInterface()
	info.SetTypeName(typeName)

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
			if s.CheckCommon() {
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
				s.toNextCh()
				continue
			}

			// 解析 interface 方法
			s.nextCh()
			name := s.scanIdentifiers()
			switch s.ch {
			case '(':
				// 解析 MathodSpec
				matchInfo := dao.NewFuncInfo()
				matchInfo.SetName(name)
				matchInfo.ParamsInPoint = s.OnParameters()
				matchInfo.ParamsOutPoint = s.OnDeclarationsResult()
				info.MatchInfos = append(info.MatchInfos, matchInfo)

			case '\n':
				// 解析 InterfaceTypeName
				iInfo := s.PackageInfo.GetType(name)
				info.MatchInfos = append(info.MatchInfos, iInfo)
			}

			s.toNextCh()
		}
		s.nextCh()
	}

	return info
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
	info.ParamsInPoint = s.OnParameters()
	info.ParamsOutPoint = s.OnDeclarationsResult()

	if s.CheckCommon() {
		common := s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		info.Common = common
	}

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
