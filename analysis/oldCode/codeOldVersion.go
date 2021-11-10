package oldCode

// func (s *source) OnType(params ...byte) (str string) {
// 	// s.nextCh()
// 	var endTag byte
// 	if params != nil && len(params) > 0 {
// 		endTag = params[0]
// 	}
// 	switch s.buf[s.r+1] {
// 	case '*': // PointerType
// 		s.nextCh()
// 		str = string(s.ch)
// 		str += s.OnType(params...)
// 	case '[': // ArrayType, SliceType
// 		if s.buf[s.r+2] == ']' { // slice type
// 			str = s.OnSliceType()
// 		} else { // array type
// 			str = s.OnArrayType()
// 		}
// 	case '<': // OutPutChanelType
// 		str = s.OnChannelType()
// 	case '.': // ArranType
// 		if string(s.buf[s.r+1:s.r+4]) == "..." {
// 			s.nextCh()
// 			s.nextCh()
// 			s.nextCh()
// 			str = "..."
// 			str += s.OnType(params...)
// 		}
// 	default:
// 		tmpntidx := s.nextTokenIdx()
// 		fmt.Println("------ " + string(s.buf[s.r:tmpntidx]))
// 		fmt.Println(string(s.buf[s.r+1 : s.r+5]))
// 		if string(s.buf[s.r+1:s.r+5]) == "map[" {
// 			str = s.OnMapType()
// 		} else {
// 			if endTag != 0 {
// 				ntidx := s.nextTokenIdx()
// 				nttidx := s.nextTargetTokenIdx(endTag)
// 				if ntidx == nttidx {
// 					s.nextToken()
// 					str = s.rangeStr()
// 				} else {
// 					rangeStr := string(s.buf[s.r+1 : ntidx])
// 					rangeStr = strings.TrimSpace(rangeStr)
// 					switch rangeStr {
// 					case "func":
// 						str = s.OnFuncType()
// 					case "interface":
// 						s.nextTargetToken(rune(endTag))
// 						str = s.rangeStr()
// 					case "struct":
// 						str = s.OnStructType()
// 					default:
// 						switch s.buf[ntidx] {
// 						case '.':
// 							str = s.OnQualifiedIdentType()
// 						}
// 					}
// 				}
// 				// tmp := string(s.buf[s.r:])
// 				// tmp = strings.TrimSpace(tmp)
// 			} else {
// 				tmp := string(s.buf[s.r:s.nextIdx()])
// 				tmp = strings.TrimSpace(tmp)
// 				switch tmp {
// 				case "":
// 					s.next()
// 					str = s.OnType(params...)
// 				case "struct":
// 					str = s.OnStructType()
// 				case "interface":
// 					str = s.OnInterfaceType()
// 				default:
// 					s.next()
// 					str = strings.TrimSpace(s.rangeStr())
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// func (s *source) OnParams(isParameters bool) [][2]string {
// 	strLit := [][2]string{}
// 	var knowData [2]string
// 	unknowData := []string{}
// 	isUnknow := true
// 	for isNext := true; isNext; {
// 		s.next()
// 		switch s.buf[s.r-1] {
// 		case ')':
// 			if isUnknow {
// 				for _, v := range unknowData {
// 					strLit = append(strLit, [2]string{"", v})
// 				}
// 				// strLit = append(strLit, [2]string{"", string(s.buf[s.b : s.r-1])})
// 			}
// 			isNext = false
// 		case ',':
// 			if isUnknow {
// 				data := string(s.buf[s.b : s.r-1])
// 				unknowData = append(unknowData, data)
// 			}
// 			isUnknow = true
// 		default:
// 			isUnknow = false
// 			knowData = [2]string{}
// 			knowData[0] = s.rangeStr()
// 			endTagIdx := s.findMostCloseEndIdx([]byte{')'})
// 			knowData[1] = s.OnType(s.buf[endTagIdx])
// 			if len(unknowData) > 0 {
// 				for _, v := range unknowData {
// 					strLit = append(strLit, [2]string{v, knowData[1]})
// 				}
// 				unknowData = make([]string, 0)
// 			}
// 			strLit = append(strLit, knowData)
// 		}
// 	}
// 	return strLit
// }

// func (s *source) OnInterfaceType() string {
// 	strLit := [][2]string{}
// 	str := ""
// 	if s.buf[s.r+1] == '}' {
// 		s.next()
// 		str = "interface{}"
// 	} else if s.buf[s.r] == '{' {
// 		str = "interface"
// 		s.next()
// 		for {
// 			// name := strings.TrimSpace(s.rangeStr())
// 			// strLit = append(strLit, [2]string{name, s.OnFuncType()})
// 			name := s.OnMethodName()
// 			Signatures := s.OnParams(true)
// 			Results := s.OnResult()
// 			strLit = append(strLit, [2]string{name, fmt.Sprint(Signatures, Results)})
// 			// 清除格式
// 			if s.ch == '\n' {
// 				for {
// 					if s.buf[s.r+1] == '\t' {
// 						s.nextCh()
// 					} else {
// 						break
// 					}
// 				}
// 			}
// 			if s.buf[s.r] == '}' {
// 				s.next()
// 				break
// 			}
// 		}
// 		str += fmt.Sprint(strLit)
// 	} else {
// 		fmt.Println("OnInterfaceType Error")
// 	}
// 	return str
// }
