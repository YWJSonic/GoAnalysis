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
func (s *source) OnDeclarationsResult() []*dao.FuncParams {
	params := []*dao.FuncParams{}
	switch s.buf[s.r+1] {
	case '{', '\t', '}':
		return params
	case '(':
		s.nextCh()
		params = s.OnParams(true)

	default:
		info := s.OnDeclarationsType()
		param := dao.NewFuncParams()
		param.ContentTypeInfo = info
		params = append(params, param)
	}
	return params
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

			baseInfo, ok := dao.BaseTypeInfo[tmpStr]
			if ok {
				s.next()
				info = baseInfo
			} else {
				if tmpStr == _struct {
					info = s.OnStructType()
				} else if tmpStr == _chan {
					info = s.onChannelType()
				} else if tmpStr == _interface {
					s.OnInterfaceType()
				} else {
					panic("")
					// s.next()
					// s.rangeStr()
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
				case _func:
					info = s.OnFuncType('\n')
				case _map:
					info = s.OnMapType()
				case _interface:
					s.OnInterfaceType()
				case _struct:
					info = s.OnStructType()
				case _chan:
					info = s.onChannelType()
				default:
					if s.buf[nextTokenIdx] == '.' {
						info = s.OnQualifiedIdentType()
					}
				}
			}
		}
	}
	return
}
