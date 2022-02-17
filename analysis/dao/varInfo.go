package dao

type VarInfo struct {
	PointBase
	TypeInfo ITypeInfo
	// Expressions *Expressions
	Expressions string // 暫代
	Tag         string
	Common      string
}

func NewVarInfo(name string) *VarInfo {
	info := &VarInfo{}
	info.name = name
	return info
}
