package dao

type Expression struct {
	ConstantType ITypeInfo // 表達式最終類型
	IExpressionContant
}

func NewExpression() *Expression {
	return &Expression{}
}

type IExpressionContant interface{}
type ExpressionArray struct {
}

// type LiteralValue struct {
// 	ElementMap map[string]ElementList //  <name, ElementList>
// }

// type ElementList struct {
// 	Key   *Expression
// 	Value *Expression
// }
