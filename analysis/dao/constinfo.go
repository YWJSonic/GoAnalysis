package dao

type ConstInfo struct {
	TypeBase
	ContentTypeInfo ITypeInfo
	Expression      string //*Expression
	Common          string
}

func NewConstInfo(name string) *ConstInfo {
	info := &ConstInfo{
		TypeBase: NewPointBase(),
	}
	info.SetName(name)
	return info
}
