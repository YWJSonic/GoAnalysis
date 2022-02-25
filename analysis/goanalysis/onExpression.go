package goanalysis

import (
	"codeanalysis/analysis/dao"
)

func (s *source) scanExpressionList() []string {
	var expressions []string
	s.nextCh()
	for {
		expressions = append(expressions, s.onFakeExpression('\n'))
		if s.ch != ',' {
			break
		}
		s.nextCh()
		s.nextCh()
	}

	return expressions
}
func (s *source) scanConstExpressionList() []*dao.Expression {
	var expressions []*dao.Expression
	s.nextCh()
	for {
		expressions = append(expressions, s.onScanConstExpression('\n'))
		if s.ch != ',' {
			break
		}
		s.nextCh()
		s.nextCh()
	}

	return expressions
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

func (s *source) onScanArrayLengthExpression(endTag rune) *dao.Expression {
	var expInfo = &dao.Expression{}
	if isInt, _ := s.isInt_lit(); isInt {
		s.scanIdentifier()
		expInfo.ConstantType = dao.BaseTypeInfo["int"]

	} else {
		identifierOrType := s.scanIdentifier()
		if s.ch == endTag {
			expInfo.ConstantType = s.PackageInfo.GetType(identifierOrType)

		} else if s.ch == '.' {
			qualifiedInfo := dao.NewTypeInfoQualifiedIdent()
			s.onQualifiedIdentifier(identifierOrType, qualifiedInfo)
			expInfo.ConstantType = qualifiedInfo
		}
	}
	return expInfo
}

func (s *source) onScanConstExpression(endTag rune) *dao.Expression {
	var expInfo = &dao.Expression{}
	if isInt, _ := s.isInt_lit(); isInt {
		s.scanIdentifier()
		expInfo.ConstantType = dao.BaseTypeInfo["int"]

	} else {
		identifierOrType := s.scanIdentifier()
		if s.ch == endTag {
			expInfo.ConstantType = s.PackageInfo.GetType(identifierOrType)

		} else if s.ch == '.' {
			qualifiedInfo := dao.NewTypeInfoQualifiedIdent()
			s.onQualifiedIdentifier(identifierOrType, qualifiedInfo)
			expInfo.ConstantType = qualifiedInfo
		}
	}
	return expInfo
}
