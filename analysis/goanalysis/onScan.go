package goanalysis

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/util"
)

//  _"string"
func (s *source) scanStringLit(strToken rune) string {
	var ch rune
	offset := s.r
	s.nextCh()
	for {

		ch = s.ch
		if ch == '\n' || ch < 0 { // 字串格式錯誤
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
