package goanalysis

import (
	"codeanalysis/analysis/dao"
	"fmt"
	"strings"
)

// StructType		= "_struct" "{" { FieldDecl ";" } "}" .
/* 進入指標應當只在
 * r="_" range="struct" // 應該調整為統一在
 */
func (s *source) OnStructType() *dao.TypeInfoStruct {
	unknowData := []string{}
	info := dao.NewTypeInfoStruct()
	info.TypeName = _struct
	s.nextToken()
	if s.buf[s.r+1] == '}' {
		s.next()
		return info
	}
	s.nextCh()

	for {
		s.toNextCh()

		if string(s.buf[s.r+1:s.r+3]) == "//" {
			s.OnComments("//")
		} else {
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
			if s.buf[s.r+1] == ' ' {
				s.toNextCh()
			}

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

		if s.buf[s.r+1] == '}' {
			s.nextTargetToken('}')
			break
		}
	}

	return info
}

func (s *source) onPointType() *dao.TypeInfoPointer {
	s.nextCh()

	info := dao.NewTypeInfoPointer()
	info.TypeName = "*"
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
	s.next()
	typeName := strings.TrimSpace(s.rangeStr())
	fullName := fmt.Sprintf("%s.%s", packageName, typeName)

	info.Name = fullName
	info.ImportLink, info.ContentTypeInfo = s.PackageInfo.GetPackageType(packageName, typeName)
	return info
}

func (s *source) OnSliceType(endTag byte) *dao.TypeInfoSlice {
	info := dao.NewTypeInfoSlice()
	info.TypeName = "[]"
	s.nextCh()
	s.nextCh()
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}

func (s *source) onShortArrayType() string {
	var str string
	if string(s.buf[s.r+1:s.r+4]) == "..." {
		s.nextCh()
		s.nextCh()
		s.nextCh()
		str = "..."
		// str += s.subDeclarationsType()
	}
	return str
}

func (s *source) onArrayType() string {
	s.nextCh()
	str := string(s.ch)
	// str += s.subDeclarationsType()
	return str

}

func (s *source) OnArrayType() string {
	s.nextCh()
	str := string(s.ch)
	s.nextTargetToken(']')
	str += s.rangeStr()
	str += string(s.ch)
	str += s.OnType()
	return str
}

func (s *source) onChannelType() *dao.TypeInfoChannel {

	info := dao.NewTypeInfoChannel()
	if s.buf[s.r+1] == '<' { // 單出
		s.next()
		info.Name = s.rangeStr()
		info.FlowType = 1

	} else {
		s.next()

		info.Name = s.rangeStr()

		// channel 方向
		if info.Name == _chan { // 雙向
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

	str = str + " " + s.OnType()
	return str
}

// 處理 map 類型
// 進入指標應當只在 _map[]
// r="_"
func (s *source) OnMapType() dao.ITypeInfo {
	info := dao.NewTypeInfoMap()
	s.nextToken()
	info.KeyType = s.OnDeclarationsType()
	if s.ch != ']' {
		s.nextCh()
	}
	info.ValueType = s.OnDeclarationsType()
	return info
}

// 處理 interface 類型
/* 進入指標應當只在
 * interface{
 * b=i r={ rang="interface"
 */
func (s *source) OnInterfaceType() string {
	strLit := [][2]string{}
	str := ""
	str = "interface"
	s.nextToken()

	if s.buf[s.r+1] == '}' {
		str += "{}"
		s.nextCh()
		s.nextCh()

	} else {
		str = "interface {"
		s.toNextCh()
		for {
			// name := strings.TrimSpace(s.rangeStr())
			// strLit = append(strLit, [2]string{name, s.OnFuncType()})
			name := s.OnFuncName()
			Signatures := s.OnParams(true)
			Results := s.OnDeclarationsResult()
			strLit = append(strLit, [2]string{name, fmt.Sprint(Signatures, Results)})

			idx := s.nextIdx()
			if strings.TrimSpace(string(s.buf[s.r+1:idx])) == "}" {
				s.next()
				break
			}
		}
		str += fmt.Sprint(strLit) + "}"
	}
	return str
}

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
func (s *source) OnFuncType(endTag byte) *dao.FuncInfo {
	info := dao.NewFuncInfo()
	info.TypeName = _func
	s.nextToken()
	info.ParamsInPoint = s.OnParams(true)
	if endTag == '\n' {
		info.ParamsOutPoint = s.OnDeclarationsResult()
	} else {
		panic("")
		s.OnResult(endTag)

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
