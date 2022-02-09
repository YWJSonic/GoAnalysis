package dao

type ConstInfo struct {
	PointBase
	TypeInfo    ITypeInfo
	Expressions *Expressions
	Common      string
}

func NewConstInfo(name string) *ConstInfo {
	info := &ConstInfo{}
	info.name = name
	return info
}
