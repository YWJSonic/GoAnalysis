package dao

type Expressions struct {
	PointBase
	Objs []ITypeInfo
}

func NewExpressions() *Expressions {
	return &Expressions{}
}
