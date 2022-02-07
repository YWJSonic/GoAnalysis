package dao

// 通用檔案節點
type FileDataNode interface {
	Path() string
	Name() string
}

// 全指標共用基底
type PointBase struct {
	Name     string
	TypeName string
}

func (self *PointBase) SetName(name string) {
	self.Name = name
}

func (self *PointBase) GetTypeName() string {
	return self.TypeName
}
