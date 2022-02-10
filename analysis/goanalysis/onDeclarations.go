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
			Signatures := s.OnParameters()
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
