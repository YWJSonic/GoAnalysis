package goanalysis

import (
	"codeanalysis/analysis/constant"
	"codeanalysis/util"
)

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

// 判斷是否為數字, 並回傳進制單位
//
// @return bool		使否為數字
// @return uint8	進制單位
func (s *source) isInt_lit() (bool, constant.IntLiteType) {
	if !util.IsDecimal(s.ch) {
		return false, constant.IntLiteType_None
	}

	nextCh := rune(s.buf[s.r+1])
	switch nextCh {
	case 'b', 'B':
		return true, constant.IntLiteType_Binary

	case 'o', 'O':
		return true, constant.IntLiteType_Octal

	case 'x', 'X':
		return true, constant.IntLiteType_Hex

	default:
		return true, constant.IntLiteType_Decimal
	}
}

// 確認後續是否為註解
// 指標指定位置 _//
// r="_"
func (s *source) checkCommon() bool {
	if s.r+1 == s.e {
		return false
	}

	if s.buf[s.r+1] != '/' {
		return false
	}

	previrwCh := string(s.buf[s.r+1 : s.r+3])
	if previrwCh != "//" && previrwCh != "/*" {
		return false
	}

	return true
}

// 判斷運算子
// 起始位置: "+" , "&"&
//
// @return Operator	// 運算子類型
// @return int 		// 運算子長度
// @return bool		// 是否為運算子
func (s *source) checkBinary_op() (constant.Operator, int, bool) {
	// "||", "&&", "==", "!=", "<", "<=", ">", ">=", "+", "-", ",", "^", "*", "/", "%", "<<", ">>", "&", "&^"
	previrwCh1 := s.buf[s.r]
	previrwCh2 := s.buf[s.r+1]
	switch previrwCh1 {
	case '+':
		return constant.Add, 1, true
	case '-':
		return constant.Sub, 1, true
	case '^':
		return constant.Xor, 1, true
	case '*':
		return constant.Mul, 1, true
	case '/':
		return constant.Div, 1, true
	case '%':
		return constant.Rem, 1, true
	case '!':
		if previrwCh2 == '=' {
			return constant.Neq, 2, true
		}
	case '=':
		if previrwCh2 == '=' {
			return constant.Eql, 2, true
		}
	case '|':
		if previrwCh2 == '|' {
			return constant.OrOr, 2, true
		} else {
			return constant.Or, 1, true
		}
	case '<':
		switch previrwCh2 {
		case '=':
			return constant.Leq, 2, true
		case '<':
			return constant.Shl, 2, true
		}
		return constant.Lss, 1, true
	case '>':
		switch previrwCh2 {
		case '=':
			return constant.Geq, 2, true
		case '>':
			return constant.Shr, 2, true
		}
		return constant.Gtr, 1, true

	case '&':
		switch previrwCh2 {
		case '&':
			return constant.AndAnd, 2, true
		case '^':
			return constant.AndNot, 2, true
		}
		return constant.And, 1, true
	}

	return 0, 0, false
}

func (s *source) checkRel_op() (constant.Operator, int, bool) {
	// "==" | "!=" | "<" | "<=" | ">" | ">=" .
	previrwCh1 := s.buf[s.r]
	previrwCh2 := s.buf[s.r+1]
	switch previrwCh1 {
	case '!':
		if previrwCh2 == '=' {
			return constant.Neq, 2, true
		}
	case '=':
		if previrwCh2 == '=' {
			return constant.Eql, 2, true
		}
	case '<':
		switch previrwCh2 {
		case '=':
			return constant.Leq, 2, true
		}
		return constant.Lss, 1, true
	case '>':
		switch previrwCh2 {
		case '=':
			return constant.Geq, 2, true
		}
		return constant.Gtr, 1, true
	}

	return 0, 0, false
}

func (s *source) checAdd_op() (constant.Operator, int, bool) {
	//  "+" | "-" | "|" | "^"
	previrwCh1 := s.buf[s.r]
	switch previrwCh1 {
	case '+':
		return constant.Add, 1, true
	case '-':
		return constant.Sub, 1, true
	case '|':
		return constant.Or, 1, true
	case '^':
		return constant.Xor, 1, true
	}
	return 0, 0, false
}

func (s *source) checMul_op() (constant.Operator, int, bool) {
	// "*" | "/" | "%" | "<<" | ">>" | "&" | "&^"
	previrwCh1 := s.buf[s.r]
	previrwCh2 := s.buf[s.r+1]
	switch previrwCh1 {
	case '*':
		return constant.Mul, 1, true
	case '/':
		return constant.Div, 1, true
	case '%':
		return constant.Rem, 1, true
	case '<':
		if previrwCh2 == '<' {
			return constant.Shl, 2, true
		}
	case '>':
		if previrwCh2 == '>' {
			return constant.Shr, 2, true
		}
	case '&':
		switch previrwCh2 {
		case '&':
			return constant.AndAnd, 2, true
		case '^':
			return constant.AndNot, 2, true
		}
	}

	return 0, 0, false
}

func (s *source) checUnary_op() (constant.Operator, int, bool) {
	// "+" | "-" | "!" | "^" | "*" | "&" | "<-"
	previrwCh1 := s.buf[s.r]
	previrwCh2 := s.buf[s.r+1]
	switch previrwCh1 {
	case '+':
		return constant.Add, 1, true
	case '-':
		return constant.Sub, 1, true
	case '!':
		return constant.Not, 1, true
	case '^':
		return constant.Xor, 1, true
	case '*':
		return constant.Mul, 1, true
	case '<':
		if previrwCh2 == '-' {
			return constant.Recv, 2, true
		}
	case '&':
		return constant.And, 1, true
	}

	return 0, 0, false
}
