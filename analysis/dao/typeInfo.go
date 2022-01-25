package dao

func NewTypeInfo(name string) *TypeInfo {
	info := &TypeInfo{}
	info.Name = name
	return info
}

type TypeInfo struct {
	PointBase
	StructInfo *StructInfo
	Tag        string
	Common     string
}
