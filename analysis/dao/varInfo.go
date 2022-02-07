package dao

type VarInfo struct {
	PointBase
	TypeInfo    ITypeInfo
	Expressions *Expressions
	Tag         string
	Common      string
}

func NewVarInfo(name string) *VarInfo {
	info := &VarInfo{}
	info.Name = name
	return info
}
