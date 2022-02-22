package goanalysis

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/util"
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

		for util.IsDecimal(s.ch) || s.ch == '_' {
			s.nextCh()
		}

	}

	switch s.ch {

	case '.': // 浮點數表示
		value += string(s.ch)
		for s.nextCh(); util.IsDecimal(s.ch); s.nextCh() {
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

// 掃描字串
//
// 起始位置 _"stringcontext"
// 起始位置 _`stringcontext`
func (s *source) scanString() *dao.Expressions {

	s.nextCh()
	baseInfo := dao.BaseTypeInfo["string"]
	baseInfo.SetName(s.scanStringLit(s.ch))

	info := dao.NewExpressions()
	info.Objs = append(info.Objs, baseInfo)
	return info
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

func (s *source) scanExpressionList() []string {
	var expressions []string
	for {
		expressions = append(expressions, s.scanExpression('\n'))
		if s.ch != ',' {
			break
		}
		s.nextCh()
		s.nextCh()
	}

	return expressions
}

func (s *source) scanExpression(endTag rune) string {
	offset := s.r
	var scanDone bool
	endTokenQueue := []rune{endTag}

	for !scanDone {
		// s.next()
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

			if len(endTokenQueue) == 1 {
				scanDone = true
			}
		}
	}
	return string(s.buf[offset:s.r])
}

// PrimaryExpr =
// 	Operand |
// 	Conversion |
// 	MethodExpr |
// 	PrimaryExpr Selector |
// 	PrimaryExpr Index |
// 	PrimaryExpr Slice |
// 	PrimaryExpr TypeAssertion |
// 	PrimaryExpr Arguments .
//
// Operand     = Literal | OperandName | "(" Expression ")" .
// Literal     = BasicLit | CompositeLit | FunctionLit .
// BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
// OperandName = identifier | QualifiedIdent .
//
// QualifiedIdent = PackageName "." identifier .
//
// Conversion = Type "(" Expression [ "," ] ")" .
// scan PrimaryExpr
func (s *source) scanPrimaryExpr() string {
	var primaryExpr string
	switch s.ch {
	case '.':
		// selector, typeAssertion
		primaryExpr = s.scanPrimaryExpr_Selector()
		primaryExpr = s.scanPrimaryExpr_TypeAssertion()
	case '[':
		// index, slice
		primaryExpr = s.scanPrimaryExpr_Index()
		primaryExpr = s.scanPrimaryExpr_Slice()
	case '(':
		// Operand->Expression,
		primaryExpr = s.scanExpression('\n')
		primaryExpr = s.scanPrimaryExpr_Arguments()
	default:
		primaryExpr = s.scanIdentifier()

		switch s.ch {
		case '.':
			// Operand->OperandName->QualifiedIdent
		}

	}
	return primaryExpr
}

func (s *source) scanPrimaryExpr_Selector() string {
	var primaryExpr string
	return primaryExpr
}

func (s *source) scanPrimaryExpr_Index() string {
	var primaryExpr string
	return primaryExpr
}

func (s *source) scanPrimaryExpr_Slice() string {
	var primaryExpr string
	return primaryExpr
}

func (s *source) scanPrimaryExpr_TypeAssertion() string {
	var primaryExpr string
	return primaryExpr
}

func (s *source) scanPrimaryExpr_Arguments() string {
	var primaryExpr string
	return primaryExpr
}

// 解析隱藏宣告
func (s *source) scanEmbeddedField() *dao.VarInfo {
	info := dao.NewVarInfo("_")
	info.ContentTypeInfo = s.OnDeclarationsType()
	return info
}
