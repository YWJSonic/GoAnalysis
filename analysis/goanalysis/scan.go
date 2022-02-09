package goanalysis

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/tool"
)

// 解析常數
func (s *source) scanNumber() *dao.Expressions {
	// var digsep bool // 常數存在 "_" 分隔符號
	var value string
	offset := s.r
	info := dao.NewExpressions()

	if s.nextCh() != '.' {

		if s.ch == '0' {
			s.nextCh()
			switch s.ch {
			case 'b': // 二進制
				s.nextCh()
			case 'x': // 十六進制
				s.nextCh()
			case 'o': // 八進制
				s.nextCh()
			default: // 預設 八進制

			}
		}

		for tool.IsDecimal(s.ch) || s.ch == '_' {
			s.nextCh()
		}

	}

	switch s.ch {

	case '.': // 浮點數表示
		value += string(s.ch)
		for s.nextCh(); tool.IsDecimal(s.ch); s.nextCh() {
			value += string(s.ch)
		}
		baseInfo := dao.BaseTypeInfo["float64"]
		baseInfo.SetName(value)
		info.Objs = append(info.Objs, baseInfo)
	default:
		baseInfo := dao.BaseTypeInfo["int"]
		baseInfo.SetName(string(s.buf[offset:s.r]))
		info.Objs = append(info.Objs, baseInfo)
	}

	return info
}

// 掃描字串
//
// 起始位置 _"stringcontext"
func (s *source) scanString() *dao.Expressions {
	var ch rune
	s.nextCh()
	offset := s.r

	s.nextCh()
	for {

		ch = s.ch
		if ch == '\n' || ch < 0 { // 字串格式錯誤
			panic("")
		}

		s.nextCh()
		if ch == '"' { // 字串結束
			break
		}

		if ch == '\\' { //跳脫符號
			s.nextCh()
		}
	}

	value := string(s.buf[offset:s.r])
	baseInfo := dao.BaseTypeInfo["string"]
	baseInfo.SetName(value)

	info := dao.NewExpressions()
	info.Objs = append(info.Objs, baseInfo)
	return info
}

func (s *source) scanIdentifiers() string {
	offset := s.r
	for {
		if tool.IsLetter(s.ch) || tool.IsDecimal(s.ch) || s.ch == '_' {
			s.nextCh()
			continue
		}
		break
	}

	return string(s.buf[offset:s.r])
}
