package goanalysis

import (
	"codeanalysis/analysis/dao"
)

// func (s *source) scanExpressionList() []string {
// 	var expressions []string
// 	s.nextCh()
// 	for {
// 		expressions = append(expressions, s.onFakeExpression('\n'))
// 		if s.ch != ',' {
// 			break
// 		}
// 		s.nextCh()
// 		s.nextCh()
// 	}

// 	return expressions
// }

func (s *source) scanExpressionList(infos []*dao.ConstInfo) []string {
	var expressions []string
	s.nextCh()
	otherEndTag := ','
	for i, count := 0, len(infos); i < count; i++ {
		if i == count-1 {
			otherEndTag = 0
		}

		expressions = append(expressions, s.onFirstScanExpression(otherEndTag))
		if s.ch == ',' {
			s.toNextCh()
			s.nextCh()
		}
	}

	return expressions
	// var expressions []string
	// s.nextCh()
	// for {
	// 	expressions = append(expressions, s.onScanConstExpression())
	// 	if s.ch != ',' {
	// 		break
	// 	}
	// 	s.nextCh()
	// 	s.nextCh()
	// }

	// return expressions
}

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

			if s.checkCommon() {
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			}

			if len(endTokenQueue) == 1 {
				scanDone = true
			}
		}
	}
	return string(s.buf[offset:s.r])
}

func (s *source) onScanConstExpression() string {
	var expInfo string
	var endTag = []rune{',', '\n', ' ', '\t'}
	var scanDone = false
	var offset = s.r

	for !scanDone {
		s.nextCh()

		for _, tag := range endTag {
			if tag == s.ch {

				if s.ch == ' ' {

					if s.checkCommon() {
						scanDone = true
						break
					}
					s.nextCh()

					if s.buf[s.r] == ' ' {
						scanDone = true
						break
					}
					if _, oplen, isop := s.checkBinary_op(); isop {
						for i := 0; i < oplen; i++ {
							s.nextCh()
						}
					}

				} else {
					scanDone = true
				}
			}
		}
	}

	expInfo = string(s.buf[offset:s.r])
	return expInfo
}

func (s *source) onFirstScanExpression(otherEndTag rune) string {
	var expInfo string
	// var scanDone = false
	var offset = s.r
	var zoneTagMap = map[rune]rune{'{': '}', '(': ')', '[': ']'}
	var strZoneTagMap = map[rune]struct{}{'\'': {}, '`': {}, '"': {}}
	var zoneQueue = []rune{}
	var strZone = false

	for {
		if s.checkEnd() {
			break
		}
		s.nextCh()

		if !strZone {
			if zoneEndTag, ok := zoneTagMap[s.ch]; ok {
				zoneQueue = append(zoneQueue, zoneEndTag)
			} else if len(zoneQueue) > 0 && s.ch == zoneQueue[len(zoneQueue)-1] {
				zoneQueue = zoneQueue[:len(zoneQueue)-1]
			} else if _, ok := strZoneTagMap[s.ch]; ok {
				strZone = true
			}
		} else {
			if s.ch == '\\' {
				s.nextCh()
			} else if _, ok := strZoneTagMap[s.ch]; ok {
				strZone = false
			}
		}

		if s.isOnNewlineSymbol() {
			if len(zoneQueue) == 0 {
				break
			}
		} else if s.ch == otherEndTag {
			break
		} else if s.checkCommon() {
			if len(zoneQueue) == 0 {
				break
			}
		}
	}

	expInfo = string(s.buf[offset:s.r])
	return expInfo
}
