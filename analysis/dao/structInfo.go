package dao

type IStructInfo interface {
	SetName(name string)
}

func NewStructInfo(name string) *StructInfo {
	structInfo := &StructInfo{}
	structInfo.Name = name
	return structInfo
}

// Struct 關聯
type StructInfo struct {
	PointBase
	Package     *PackageInfo
	StructInfos map[string]*StructInfo
	FuncPoint   *FuncInfo
}

func (self *StructInfo) SetName(name string) {
	self.Name = name
}

type StructChannelInfo struct {
	PointBase
	Package    *PackageInfo
	StructInfo *StructInfo
	FuncPoint  *FuncInfo
	ChanFlow   string
}
