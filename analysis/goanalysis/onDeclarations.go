package goanalysis

import (
	"fmt"
	"strings"
)

func (s *source) OnDeclarationsInterfaceType() string {
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
		s.next()
		for {
			// name := strings.TrimSpace(s.rangeStr())
			// strLit = append(strLit, [2]string{name, s.OnFuncType()})
			name := s.OnFuncName()
			Signatures := s.OnParams(true)
			Results := s.OnDeclarationsInterfaceType()
			strLit = append(strLit, [2]string{name, fmt.Sprint(Signatures, Results)})

			idx := s.nextIdx()
			if strings.TrimSpace(string(s.buf[s.r+1:idx])) == "}" {
				s.next()
				break
			}
		}
		str = fmt.Sprint(strLit)
	}
	return str
}

// func 輸出參數
func (s *source) OnDeclarationsResult() [][]string {
	strLit := [][]string{}
	switch s.buf[s.r+1] {
	case '{', '\t', '}':
		return strLit
	case '(':
		s.nextCh()
		strLit = s.OnParams(true)

	default:
		strLit = append(strLit, []string{"", s.OnDeclarationsType()})
	}
	return strLit
}

// 處理 type 參數類型
/* 進入指標應當只在 _types
 * b=任意 r="_"
 */
func (s *source) OnDeclarationsType() (str string) {
	// s.nextCh()
	switch s.buf[s.r+1] {
	case '*': // PointerType
		s.nextCh()
		str = string(s.ch)
		str += s.OnDeclarationsType()
	case '[': // ArrayType, SliceType
		if s.buf[s.r+2] == ']' { // slice type
			s.nextCh()
			str += string(s.ch)
			s.nextCh()
			str += string(s.ch)
			str += s.OnDeclarationsType()
			return str
		} else { // array type
			s.nextCh()
			str += string(s.ch)
			s.nextTargetToken(']')
			str += s.rangeStr()
			str += string(s.ch)
			str += s.OnDeclarationsType()
		}
	case '<': // OutPutChanelType
		if s.buf[s.r+1] == '<' { // 單出
			s.next()
			str = s.rangeStr()

		} else {
			s.next()
			str = s.rangeStr()
		}

		str = str + " " + s.OnDeclarationsType()
	case '.': // ArranType
		if string(s.buf[s.r+1:s.r+4]) == "..." {
			s.nextCh()
			s.nextCh()
			s.nextCh()
			str = "..."
			str += s.OnDeclarationsType()
		}
	default:
		nextEndIdx := s.nextIdx()
		nextTokenIdx := s.nextTokenIdx()

		if nextEndIdx < nextTokenIdx {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextEndIdx]))
			if tmpStr == _struct {
				str = s.OnStructType()
			} else if tmpStr == _chan {
				if s.buf[s.r+1] == '<' { // 單出
					s.next()
					str = s.rangeStr()

				} else {
					s.next()
					str = s.rangeStr()
				}
				str = str + " " + s.OnDeclarationsType()
			} else if tmpStr == _interface {
				str = s.OnInterfaceType()
			} else {
				s.next()
				str = s.rangeStr()
			}
		} else if nextTokenIdx == nextEndIdx {
			s.nextToken()
		} else if s.buf[nextTokenIdx] == '.' {
			str = s.OnQualifiedIdentType()
		} else if s.buf[nextTokenIdx] == '(' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _func {
				str = s.OnFuncType('\n')
			}
		} else if s.buf[nextTokenIdx] == '{' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _interface {
				str = s.OnInterfaceType()
			} else if tmpStr == _struct {
				str = s.OnStructType()
			}
		} else if s.buf[nextTokenIdx] == '[' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _map {
				str = s.OnDeclarationsMapType()
			}
		} else if s.buf[nextTokenIdx] == '<' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _chan {
				if s.buf[s.r+1] == '<' { // 單出
					s.next()
					str = s.rangeStr()

				} else {
					s.next()
					str = s.rangeStr()
				}

				str = str + " " + s.OnDeclarationsType()
			}
		}
	}

	return
}

// 處理 map 類型
// 進入指標應當只在 _map[]
// r="_"
func (s *source) OnDeclarationsMapType() (str string) {
	str = "map["
	s.nextToken()
	str += s.OnType(']') + "]"
	if s.ch != ']' {
		s.nextCh()
	}
	str += s.OnDeclarationsType()
	return
}
