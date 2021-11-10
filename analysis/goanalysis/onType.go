package goanalysis

import (
	"fmt"
	"strings"
)

// 處理 type 參數類型
/* 進入指標應當只在 _types
 * b=任意 r="_"
 */
func (s *source) OnType(params ...byte) (str string) {
	// s.nextCh()
	var endTag byte
	if len(params) > 0 {
		endTag = params[0]
	}
	switch s.buf[s.r+1] {
	case '*': // PointerType
		s.nextCh()
		str = string(s.ch)
		str += s.OnType(params...)
	case '[': // ArrayType, SliceType
		if s.buf[s.r+2] == ']' { // slice type
			str = s.OnSliceType(endTag)
		} else { // array type
			str = s.OnArrayType()
		}
	case '<': // OutPutChanelType
		str = s.OnChannelType()
	case '.': // ArranType
		if string(s.buf[s.r+1:s.r+4]) == "..." {
			s.nextCh()
			s.nextCh()
			s.nextCh()
			str = "..."
			str += s.OnType(params...)
		}
	default:

		nextTokenIdx := s.nextTokenIdx()
		endTokenIdx := s.nextTargetTokenIdx(endTag)
		if nextTokenIdx == endTokenIdx {
			s.nextToken()
			nextRangeStr := string(s.rangeStr())
			str = nextRangeStr
		} else if endTokenIdx < nextTokenIdx {
			nextRangeStr := string(s.buf[s.r+1 : endTokenIdx])
			nextRangeStr = strings.TrimSpace(nextRangeStr)
			switch nextRangeStr {
			case _interface:
				str = s.OnInterfaceType()
			case _func:
				str = s.OnFuncType(endTag)
			case _struct:
				str = s.OnStructType()
			default:
				if s.buf[nextTokenIdx] == '.' {
					str = s.OnQualifiedIdentType()
				} else {
					s.nextTargetToken(rune(endTag))
					str = nextRangeStr
				}
			}
		} else {
			nextRangeStr := string(s.buf[s.r+1 : nextTokenIdx])
			nextRangeStr = strings.TrimSpace(nextRangeStr)
			switch nextRangeStr {
			case _interface:
				str = s.OnInterfaceType()
			case _func:
				str = s.OnFuncType(endTag)
			case _struct:
				str = s.OnStructType()
			case _map:
				str = s.OnMapType()
			default:
				if s.buf[nextTokenIdx] == '.' {
					str = s.OnQualifiedIdentType()
				} else {
					s.nextToken()
					str = nextRangeStr
				}
			}
		}

	}

	return
}

// 處理 package type 類型
// QualifiedIdent = PackageName "." identifier .
/* 進入指標應當只在
 * package.ide
 * b="p" r="." range="package"
 */
func (s *source) OnQualifiedIdentType() string {
	s.nextToken()
	str := s.rangeStr() + "."
	s.toFirstNextTarget([]byte{')', ',', ' ', ']', '\n'})
	str += s.rangeStr()
	return str
}

func (s *source) OnSliceType(endTag byte) string {
	s.nextCh()
	str := string(s.ch)
	s.nextCh()
	str += string(s.ch)
	str += s.OnType(endTag)
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

// 處理 struct 類型
//StructType    = "struct" "{" { FieldDecl ";" } "}" .
/* 進入指標應當只在
 * struct_{  }
 * r="_" range="struct" // 應該調整為統一在 "{"
 * struct{}
 * r="{" range="struct"
 */
func (s *source) OnStructTypeOld() string {
	str := ""
	if s.buf[s.r+1] == '}' {
		str = "{}"
	} else {
		s.next()
		strLit := [][2]string{}
		unknowData := []string{}
		var name string
		for {
			nextIdx := s.nextIdx()

			if !(s.buf[nextIdx] == '\n') {
				s.next()
				name = strings.TrimSpace(s.rangeStr())
				if s.buf[s.r-1] == ',' {
					unknowData = append(unknowData, name[:len(name)-1])
					continue
				}
			}

			typename := s.OnType()
			for _, nameData := range unknowData {
				strLit = append(strLit, [2]string{nameData, typename})
			}
			unknowData = make([]string, 0)
			strLit = append(strLit, [2]string{name, typename})

			// 清除格式
			if s.ch == '\n' {
				for {
					if s.buf[s.r+1] == '\t' {
						s.nextCh()
					} else {
						break
					}
				}
			}
			if s.buf[s.r] == '}' {
				s.next()
				break
			}
		}
		str = fmt.Sprint(strLit)
	}

	return str
}

// StructType		= "_struct" "{" { FieldDecl ";" } "}" .
/* 進入指標應當只在
 * r="_" range="struct" // 應該調整為統一在
 */
func (s *source) OnStructType() (str string) {
	str = "struct"
	strLit := [][]string{}
	unknowData := []string{}
	s.nextToken()
	if s.buf[s.r+1] == '}' {
		str += "{}"
		s.nextCh()
		return
	}
	s.nextCh()

	for {
		s.toNextCh()

		if string(s.buf[s.r+1:s.r+3]) == "//" {
			s.OnTag()
		} else {
			var name string
			nextIdx := s.nextIdx()

			if !(s.buf[nextIdx] == '\n') {
				s.next()
				name = strings.TrimSpace(s.rangeStr())
				if s.buf[s.r-1] == ',' {
					unknowData = append(unknowData, name[:len(name)-1])
					continue
				}
			}

			if s.buf[s.r+1] == ' ' {
				s.toNextCh()
			}

			typename := s.OnDeclarationsType()
			for _, nameData := range unknowData {
				strLit = append(strLit, []string{nameData, typename})
			}
			unknowData = make([]string, 0)
			tmpStrLit := []string{name, typename}
			s.toNextCh()
			for s.ch == ' ' {
				if string(s.buf[s.r+1:s.r+3]) == "//" {
					s.OnTag()
				} else if s.buf[s.r+1] == '`' {
					tmpStrLit = append(tmpStrLit, s.OnJsonTag())
					s.toNextCh()
				}
			}
			strLit = append(strLit, tmpStrLit)
		}

		// idx := s.nextEndIdx('}')
		if s.buf[s.r+1] == '}' {
			s.nextTargetToken('}')
			break
		}
	}
	str += fmt.Sprint(strLit)

	return str
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
func (s *source) OnFuncType(endTag byte) string {
	s.nextToken()
	str := "func("
	strLit := s.OnParams(true)
	for i, s2 := range strLit {
		if i == len(strLit)-1 {
			str = str + s2[0] + " " + s2[1]
		} else {
			str = str + s2[0] + " " + s2[1] + ", "
		}
	}
	str += ")"
	if endTag == '\n' {
		strLit = s.OnDeclarationsResult()
	} else {
		strLit = s.OnResult(endTag)

	}
	if len(strLit) == 1 && strLit[0][0] == "" {
		str = str + strLit[0][0] + " " + strLit[0][1]
	} else {
		for i, s2 := range strLit {
			if i == len(strLit)-1 {
				str = str + s2[0] + " " + s2[1] + ")"
			} else {
				str = str + s2[0] + " " + s2[1] + ", "
			}
		}
	}

	return str
}

// 處理 map 類型
// 進入指標應當只在 _map[]
// r="_"
func (s *source) OnMapType() (str string) {
	str = "map["
	s.nextToken()
	str += s.OnType(']') + "]"
	if s.ch != ']' {
		s.nextCh()
	}
	str += s.OnType()
	return
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
