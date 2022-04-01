package goanalysis

import (
	"codeanalysis/analysis/dao"
	"fmt"
)

// var 第一階段表達式解析接口
func (s *source) scanVarExpressionList(infos []*dao.VarInfo) []*dao.Expression {
	var expressions []*dao.Expression
	s.nextCh()
	otherEndTag := ','
	for i, count := 0, len(infos); i < count; i++ {
		if i == count-1 {
			otherEndTag = 0
		}

		exp := dao.NewExpression()
		exp.ContentStr = s.onFirstScanExpression(otherEndTag)
		expressions = append(expressions, exp)
		fmt.Println(exp.ContentStr)

		if s.ch == ',' {
			s.toNextCh()
			s.nextCh()
		}
	}
	return expressions
}

// const 第一階段表達式解析接口
func (s *source) scanConstExpressionList(infos []*dao.ConstInfo) []*dao.Expression {
	var expressions []*dao.Expression
	s.nextCh()
	otherEndTag := ','
	for i, count := 0, len(infos); i < count; i++ {
		if i == count-1 {
			otherEndTag = 0
		}

		exp := dao.NewExpression()
		exp.ContentStr = s.onFirstScanExpression(otherEndTag)
		expressions = append(expressions, exp)
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

		// 判斷是否在 文字區塊內
		if !strZone {
			if zoneEndTag, ok := zoneTagMap[s.ch]; ok {
				// 進入子區塊
				zoneQueue = append(zoneQueue, zoneEndTag)
			} else if len(zoneQueue) > 0 && s.ch == zoneQueue[len(zoneQueue)-1] {
				// 離開子區塊
				zoneQueue = zoneQueue[:len(zoneQueue)-1]
			} else if _, ok := strZoneTagMap[s.ch]; ok {
				// 進入文字區塊
				zoneQueue = append(zoneQueue, s.ch)
				strZone = true
			}
		} else {
			// 文字區塊內
			if s.ch == '\\' {
				// 跳脫字元處理
				s.nextCh()
			} else if zoneQueue[len(zoneQueue)-1] == s.ch {
				// 離開文字區塊
				zoneQueue = zoneQueue[:len(zoneQueue)-1]
				strZone = false
			}
		}

		if s.ch == ' ' && s.buf[s.r+1] == ' ' {
			s.toNextCh()
		}

		if !strZone {
			// 不再區域內才能結束
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

		if s.checkEnd() {
			break
		}
		s.nextCh()
	}

	expInfo = string(s.buf[offset:s.r])
	return expInfo
}
