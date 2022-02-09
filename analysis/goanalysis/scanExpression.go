package goanalysis

func (s *source) scanExpression() string {
	offset := s.r
	endToken := map[rune]struct{}{
		']':  {},
		')':  {},
		'\n': {},
	}
	for {
		if _, ok := endToken[s.ch]; ok {
			break
		}
		s.nextCh()
	}

	return string(s.buf[offset:s.r])
}
