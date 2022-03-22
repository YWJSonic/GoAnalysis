package goanalysis

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/util"
)

// _'\t'
func (s *source) scanRuneLit() string {
	if s.buf[s.r+1] == '\\' {
		// unicode 编码 or byte 字元 or 跳脱字元
		s.nextCh()
		offset := s.r

		switch s.buf[s.r+1] {
		case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '\'', '"':
			// 跳脱字元
			s.nextCh()
			s.nextCh()
			return string(s.buf[offset:s.r])
		case 'x', 'u', 'U':
			//unicode 编码
			s.nextCh()
			for {
				if !util.IsHex(s.ch) {
					break
				}
			}

			str := string(s.buf[offset:s.r])
			s.nextCh()
			return str
		default:

			if util.IsOctal(rune(s.buf[s.r+1])) {
				// byte 字元
				s.nextCh()
				for {
					if util.IsLetter(s.ch) || util.IsDecimal(s.ch) || s.ch == '_' {
						s.nextCh()
						continue
					}
					break
				}
				return string(s.buf[offset:s.r])
			} else {

			}

			str := string(s.buf[offset:s.r])
			s.nextCh()
			return str

		}

	} else {
		// 单一字元
		s.nextCh()
		str := string(s.ch)
		s.nextCh()
		return str
	}
}

//  _"string"
func (s *source) scanStringLit(strToken rune) string {
	var ch rune
	offset := s.r
	s.nextCh()
	for {

		ch = s.ch
		if (ch == '\n' && strToken == '"') || ch < 0 { // 字串格式錯誤
			panic("")
		}

		s.nextCh()
		if ch == strToken { // 字串結束
			break
		}

		if ch == '\\' { //跳脫符號
			s.nextCh()
		}
	}
	return string(s.buf[offset:s.r])
}

func (s *source) scanIdentifiers() []string {
	var names []string
	for {
		names = append(names, s.scanIdentifier())
		if s.ch != ',' {
			break
		}
		s.nextCh()
		s.nextCh()
	}
	return names
}

func (s *source) scanIdentifier() string {
	offset := s.r
	for {
		if util.IsLetter(s.ch) || util.IsDecimal(s.ch) || s.ch == '_' {
			s.nextCh()
			continue
		}
		break
	}

	return string(s.buf[offset:s.r])
}

// 解析隱藏宣告
func (s *source) scanEmbeddedField() *dao.VarInfo {
	info := dao.NewVarInfo("_")
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}
