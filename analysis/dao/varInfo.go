package dao

type VarInfo struct {
	TypeBase
	ContentTypeInfo ITypeInfo
	// Expressions *Expressions
	Expressions string // 暫代
	Tag         string
	Common      string
}

func NewVarInfo(name string) *VarInfo {
	info := &VarInfo{
		TypeBase: NewPointBase(),
	}
	info.SetName(name)
	return info
}
