package goanalysis

// 處理 type 參數類型
/* 進入指標應當只在 _types
 * b=任意 r="_"
 */
func (s *source) OnType(params ...byte) (str string) {
	// // s.nextCh()
	// var endTag byte
	// if len(params) > 0 {
	// 	endTag = params[0]
	// }
	// switch s.buf[s.r+1] {
	// case '*': // PointerType
	// 	s.nextCh()
	// 	str = string(s.ch)
	// 	str += s.OnType(params...)
	// case '[': // ArrayType, SliceType
	// 	if s.buf[s.r+2] == ']' { // slice type
	// 		str = s.OnSliceType(endTag)
	// 	} else { // array type
	// 		str = s.OnArrayType()
	// 	}
	// case '<': // OutPutChanelType
	// 	str = s.OnChannelType()
	// case '.': // ArranType
	// 	if string(s.buf[s.r+1:s.r+4]) == "..." {
	// 		s.nextCh()
	// 		s.nextCh()
	// 		s.nextCh()
	// 		str = "..."
	// 		str += s.OnType(params...)
	// 	}
	// default:

	// 	nextTokenIdx := s.nextTokenIdx()
	// 	endTokenIdx := s.nextTargetTokenIdx(endTag)
	// 	if nextTokenIdx == endTokenIdx {
	// 		s.nextToken()
	// 		nextRangeStr := string(s.rangeStr())
	// 		str = nextRangeStr
	// 	} else if endTokenIdx < nextTokenIdx {
	// 		nextRangeStr := string(s.buf[s.r+1 : endTokenIdx])
	// 		nextRangeStr = strings.TrimSpace(nextRangeStr)
	// 		switch nextRangeStr {
	// 		case _interface:
	// 			str = s.OnInterfaceType()
	// 		case _func:
	// 			str = s.OnFuncType(endTag)
	// 		case _struct:
	// 			str = s.OnStructType()
	// 		default:
	// 			if s.buf[nextTokenIdx] == '.' {
	// 				str = s.OnQualifiedIdentType()
	// 			} else {
	// 				s.nextTargetToken(rune(endTag))
	// 				str = nextRangeStr
	// 			}
	// 		}
	// 	} else {
	// 		nextRangeStr := string(s.buf[s.r+1 : nextTokenIdx])
	// 		nextRangeStr = strings.TrimSpace(nextRangeStr)
	// 		switch nextRangeStr {
	// 		case _interface:
	// 			str = s.OnInterfaceType()
	// 		case _func:
	// 			str = s.OnFuncType(endTag)
	// 		case _struct:
	// 			str = s.OnStructType()
	// 		case _map:
	// 			str = s.OnMapType()
	// 		default:
	// 			if s.buf[nextTokenIdx] == '.' {
	// 				str = s.OnQualifiedIdentType()
	// 			} else {
	// 				s.nextToken()
	// 				str = nextRangeStr
	// 			}
	// 		}
	// 	}

	// }

	return
}
