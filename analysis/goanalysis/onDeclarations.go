package goanalysis

import (
	"codeanalysis/analysis/dao"
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
		s.OnDeclarationsType()
		var name string
		strLit = append(strLit, []string{"", name})
	}
	return strLit
}

// 處理 type 參數類型
/* 進入指標應當只在 _types
 * b=任意 r="_"
 */
func (s *source) OnDeclarationsType() (structInfo dao.IStructInfo) {

	// s.nextCh()
	switch s.buf[s.r+1] {
	case '*': // PointerType
		s.onPointType()
	case '[': // ArrayType, SliceType
		if s.buf[s.r+2] == ']' { // slice type
			s.OnSliceType('\n')
		} else { // array type
			s.OnArrayType()
		}

	case '<': // OutPutChanelType
		s.onChannelType()
	case '.': // short ArranType
		s.onShortArrayType()
	default:
		nextEndIdx := s.nextIdx()
		nextTokenIdx := s.nextTokenIdx()

		if nextEndIdx < nextTokenIdx {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextEndIdx]))
			if tmpStr == _struct {
				structInfo = s.OnStructType()
			} else if tmpStr == _chan {
				s.next()
				structChannelInfo := &dao.StructChannelInfo{}
				structChannelInfo.ChanFlow = s.rangeStr()
				structChannelInfo.StructInfo = s.OnDeclarationsType().(*dao.StructInfo)
				structInfo = structChannelInfo
			} else if tmpStr == _interface {
				s.OnInterfaceType()
			} else {
				s.next()
				s.rangeStr()
			}
		} else if nextTokenIdx == nextEndIdx {
			s.nextToken()
		} else if s.buf[nextTokenIdx] == '.' {
			s.OnQualifiedIdentType()
		} else if s.buf[nextTokenIdx] == '(' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _func {
				s.OnFuncType('\n')
			}
		} else if s.buf[nextTokenIdx] == '{' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _interface {
				s.OnInterfaceType()
			} else if tmpStr == _struct {
				s.OnStructType()
			}
		} else if s.buf[nextTokenIdx] == '[' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _map {
				s.OnMapType()
			}
		} else if s.buf[nextTokenIdx] == '<' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _chan {
				if s.buf[s.r+1] == '<' { // 單出
					s.next()
					s.rangeStr()

				} else {
					s.next()
					s.rangeStr()
				}

				s.subDeclarationsType()
			}
		}
	}

	return
}

func (s *source) subDeclarationsType() (str string) {
	// s.nextCh()
	switch s.buf[s.r+1] {
	case '*': // PointerType
		s.nextCh()
		str = string(s.ch)
		str += s.subDeclarationsType()
	case '[': // ArrayType, SliceType
		if s.buf[s.r+2] == ']' { // slice type
			s.nextCh()
			str += string(s.ch)
			s.nextCh()
			str += string(s.ch)
			str += s.subDeclarationsType()
			return str
		} else { // array type
			s.nextCh()
			str += string(s.ch)
			s.nextTargetToken(']')
			str += s.rangeStr()
			str += string(s.ch)
			str += s.subDeclarationsType()
		}
	case '<': // OutPutChanelType
		if s.buf[s.r+1] == '<' { // 單出
			s.next()
			str = s.rangeStr()

		} else {
			s.next()
			str = s.rangeStr()
		}

		str = str + " " + s.subDeclarationsType()
	case '.': // ArranType
		if string(s.buf[s.r+1:s.r+4]) == "..." {
			s.nextCh()
			s.nextCh()
			s.nextCh()
			str = "..."
			str += s.subDeclarationsType()
		}
	default:
		nextEndIdx := s.nextIdx()
		nextTokenIdx := s.nextTokenIdx()

		if nextEndIdx < nextTokenIdx {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextEndIdx]))
			if tmpStr == _struct {
				s.OnStructType()
			} else if tmpStr == _chan {
				if s.buf[s.r+1] == '<' { // 單出
					s.next()
					str = s.rangeStr()

				} else {
					s.next()
					str = s.rangeStr()
				}
				str = str + " " + s.subDeclarationsType()
			} else if tmpStr == _interface {
				str = s.OnInterfaceType()
			} else {
				s.next()
				str = s.rangeStr()
			}
		} else if nextTokenIdx == nextEndIdx {
			s.nextToken()
		} else if s.buf[nextTokenIdx] == '.' {
			s.OnQualifiedIdentType()
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
				// str = s.OnStructType()
			}
		} else if s.buf[nextTokenIdx] == '[' {
			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
			if tmpStr == _map {
				str = s.OnMapType()
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

				str = str + " " + s.subDeclarationsType()
			}
		}
	}

	return
}
