package goanalysis

import (
	"codeanalysis/analysis/dao"
)

// var 第一階段表達式解析接口
func (s *source) scanVarExpressionList(infos []*dao.VarInfo) []string {
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
}

// const 第一階段表達式解析接口
func (s *source) scanConstExpressionList(infos []*dao.ConstInfo) []string {
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
}

//  第一階段表達式解析
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
