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

func (s *source) scanExpressionList() []string {
	var expressions []string
	s.nextCh()
	for {
		expressions = append(expressions, s.onScanExpression('\n'))
		if s.ch != ',' {
			break
		}
		s.nextCh()
		s.nextCh()
	}

	return expressions
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
		primaryExpr = s.onScanExpression('\n')
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
