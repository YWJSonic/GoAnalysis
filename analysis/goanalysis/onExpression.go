package goanalysis

import (
	"codeanalysis/analysis/dao"
)

func (s *source) onFakeExpression(endTag rune) string {
	offset := s.r
	var scanDone bool
	endTokenQueue := []rune{endTag}

	for !scanDone {
		s.nextCh()

		if s.ch == endTag {
			if s.buf[s.r-1] == '{' {
				endTokenQueue = append(endTokenQueue, '}')
				continue
			} else if s.buf[s.r-1] == '(' {
				endTokenQueue = append(endTokenQueue, ')')
				continue
			}
			if len(endTokenQueue) > 2 && s.buf[s.r-1] == ',' {
				if rune(s.buf[s.r-2]) == endTokenQueue[len(endTokenQueue)-1] {
					endTokenQueue = endTokenQueue[:len(endTokenQueue)-1]
				}

			} else if len(endTokenQueue) > 1 {
				if rune(s.buf[s.r-1]) == endTokenQueue[len(endTokenQueue)-1] {
					endTokenQueue = endTokenQueue[:len(endTokenQueue)-1]
				}

			}

			if s.CheckCommon() {
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			}

			if len(endTokenQueue) == 1 {
				scanDone = true
			}
		}
	}
	return string(s.buf[offset:s.r])
}

func (s *source) onScanArrayLengthExpression(endTag rune) *dao.Expressions {
	var expInfo = &dao.Expressions{}

	return expInfo
}

func (s *source) onScanExpression(endTag rune) string {
	return ""
}

// func (s *source) onScanExpression(endTag rune) *dao.Expressions {
// 	var expInfo = &dao.Expressions{}
// 	var scanDone bool
// 	var identifier string
// 	for !scanDone {

// 		if util.IsDecimal(s.ch) {
// 			// 開頭為數字
// 			identifier = s.scanIdentifier()

// 		} else if util.IsLetter(s.ch) {
// 			// 開頭為字串
// 			identifier = s.scanIdentifier()

// 			qualifiedInfo := dao.NewTypeInfoQualifiedIdent()
// 			s.onQualifiedIdentifier(identifier, qualifiedInfo)

// 		} else if util.IsToken(s.ch) {
// 			// 開頭為符號

// 			switch s.ch {
// 			case '&':

// 				if util.IsLetter(rune(s.buf[s.r+1])) {
// 					// point type
// 					pointTypeInfo := dao.NewTypeInfoPointer()
// 					pointTypeInfo.ContentTypeInfo = s.OnDeclarationsType()
// 					switch ctInfo:= pointTypeInfo.ContentTypeInfo.(type){
// 					case *dao.TypeInfoQualifiedIdent:
// 						ctInfo.ContentTypeInfo

// 					}
// 					expInfo.ConstantType = pointTypeInfo
// 					expInfo.ConstantTypes = append(expInfo.ConstantTypes, pointTypeInfo)
// 					s.funcBodyBlock()
// 				}

// 			case '+':
// 			case '-':
// 			case '!':
// 			case '^':
// 			case '*':
// 			case '<':

// 			}

// 			identifier = s.scanIdentifier()

// 		}

// 		scanDone = true
// 		fmt.Println(identifier)
// 	}

// 	return expInfo
// }
