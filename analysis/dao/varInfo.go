package dao

type VarInfo struct {
	TypeBase
	ContentTypeInfo ITypeInfo
	// Expression *Expression
	Expression string // 暫代
	Tag        string
	Common     string
}

func (self *VarInfo) GetTypeBase() TypeBase {
	return self.TypeBase
}

func NewVarInfo(name string) *VarInfo {
	info := &VarInfo{
		TypeBase: NewPointBase(),
	}
	info.SetName(name)
	return info
}
