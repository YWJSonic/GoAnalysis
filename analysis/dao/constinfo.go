package dao

type ConstInfo struct {
	TypeBase
	ContentTypeInfo ITypeInfo
	Expression      *Expression
	Common          string
}

func (self *ConstInfo) GetTypeBase() TypeBase {
	return self.TypeBase
}

func NewConstInfo(name string) *ConstInfo {
	info := &ConstInfo{
		TypeBase: NewPointBase(),
	}
	info.SetName(name)
	return info
}
