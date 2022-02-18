package goanalysis

import "codeanalysis/analysis/constant"

// 索引開始
func (s *source) start() {
	s.b, s.r = 0, 0
	s.e = len(s.buf)
}

func (s *source) checkEnd() bool {
	if s.r == s.e {
		s.isEnd = true
		return true
	}
	return false
}

// 取得掃描範圍的字串
func (s *source) rangeStr() string {
	return string(s.buf[s.b:s.r])
}

// 移動一個文字
func (s *source) nextCh() rune {
	s.r++
	s.b = s.r
	s.ch = rune(s.buf[s.r])
	s.checkEnd()
	return s.ch
}

// 跳至下一個文字存在"前"的位置
// func (s *source) toNextChOld() rune {
// 	var rp1 byte = s.buf[s.r+1]
// 	for rp1 == ' ' || rp1 == '\t' || rp1 == '\n' || rp1 == '\r' {
// 		s.nextCh()
// 		rp1 = s.buf[s.r+1]
// 	}
// 	return s.ch
// }

// 跳至下一個文字存在"前"的位置
func (s *source) toNextCh() rune {
	var rp1 byte
	var beSkip = true
	for nextR := s.r + 1; ; nextR++ {
		if s.e == nextR {
			break
		}
		rp1 = s.buf[nextR]
		if rp1 == '\\' && beSkip { // 遇到第一次跳脫字元
			beSkip = false
		} else if rp1 != ' ' && rp1 != '\t' && rp1 != '\n' && rp1 != '\r' {
			break
		}
		s.nextCh()
	}
	return s.ch
}

func (s *source) toNextToken() rune {
	var rp1 byte = s.buf[s.r+1]
	var isRun bool = true
	for isRun {
		for _, token := range constant.TokenLit {
			if token == rp1 {
				isRun = false
			}
		}

		if isRun {
			s.nextCh()
			rp1 = s.buf[s.r+1]
		}

		if s.checkEnd() {
			isRun = false
		}
	}
	return s.ch
}

// 跳至下一個 ' ' or '\n'
// return isLineEnd
func (s *source) next() {
	s.b = s.r + 1
	for {
		s.r++
		if s.checkEnd() {
			return
		}

		if s.buf[s.r] == ' ' || s.isOnNewlineSymbol() {
			s.ch = rune(s.buf[s.r])
			return
		}
	}
}

// 取得下一個 ' ' or '\n' 位置
func (s *source) nextIdx() (tmpr int) {
	tmpr = s.nextIdxByPosition(s.r)
	return
}

// 取得從指定位置開始至下一個 ' ' or '\n'
func (s *source) nextIdxByPosition(idx int) (tmpr int) {
	tmpr = idx
	for {
		tmpr++
		if tmpr == s.e {
			return
		}

		if s.buf[tmpr] == ' ' || s.buf[tmpr] == '\n' {
			return
		}
	}
}

// 跳至下一個 token 位置
func (s *source) nextToken() {
	s.b = s.r + 1
	defer func() {
		if s.checkEnd() {
			return
		}

		s.ch = rune(s.buf[s.r])
	}()
	for {
		s.r++
		if s.checkEnd() {
			return
		}

		for _, token := range constant.TokenLit {
			if token == s.buf[s.r] {
				return
			}
		}
	}
}

// 取得下一個 token 所在位置
func (s *source) nextTokenIdx() (tmpr int) {
	tmpr = s.nextTokenIdxByPosition(s.r)
	return
}

// 取得從指定位置開始往下一遇到的第一個 token 位置
func (s *source) nextTokenIdxByPosition(idx int) (tmpr int) {
	tmpr = idx
	for {
		tmpr++
		if tmpr >= s.e {
			return
		}

		for _, token := range constant.TokenLit {
			if token == s.buf[tmpr] {
				return
			}
		}
	}
}

// 跳至下一個目標 token 位置
func (s *source) nextTargetToken(ch rune) {
	s.b = s.r + 1
	for {
		s.r++
		if s.checkEnd() {
			return
		}

		s.ch = rune(s.buf[s.r])
		if s.ch == ch {
			return
		}
	}
}

// 跳至下一個目標 token 位置
func (s *source) nextTargetTokenIdx(ch byte) (tmpr int) {
	tmpr = s.r
	for {
		tmpr++
		if tmpr == s.e {
			return
		}

		if s.buf[tmpr] == ch {
			return
		}
	}
}

// 取得下一個結束位置
// @params 指定結束位元 0=只使用預設 '\n'
// @return 結束位置
func (s *source) nextEndIdx(endTarget byte) (tmpr int) {
	tmpr = s.r
	for {
		tmpr++
		if tmpr == s.e {
			return
		}

		if s.buf[tmpr] == endTarget || s.buf[tmpr] == '\n' {
			return
		}
	}
}

// 跳至第一個符合陣列內元素的物件
func (s *source) toFirstNextTarget(targets []byte) {
	s.b = s.r + 1
	for {
		s.r++
		if s.checkEnd() {
			return
		}

		for _, target := range targets {
			if s.buf[s.r] == target {
				s.ch = rune(s.buf[s.r])
				return
			}
		}
	}
}

func (s *source) isOnNewlineSymbol() bool {
	if s.buf[s.r] == '\n' || (s.buf[s.r] == '\r' && s.buf[s.r+1] == '\n') {
		return true
	}
	return false
}
