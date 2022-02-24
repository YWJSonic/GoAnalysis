package dao

type Expressions struct {
	TypeBase
	ConstantType ITypeInfo // 表達式最終類型
	IExpressionContant
}

func NewExpressions() *Expressions {
	return &Expressions{
		TypeBase: NewPointBase(),
	}
}

type IExpressionContant interface{}
type ExpressionsArray struct {
}

// type LiteralValue struct {
// 	ElementMap map[string]ElementList //  <name, ElementList>
// }

// type ElementList struct {
// 	Key   *Expressions
// 	Value *Expressions
// }
