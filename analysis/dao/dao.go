package dao

import "codeanalysis/types"

// 通用檔案節點
type FileDataNode interface {
	Path() string
	Name() string
}

// 全指標共用基底

// 名稱 基礎物件
type PointBase struct {
	name     string
	typeName string
}

func (self *PointBase) GetName() string {
	return self.name
}

func (self *PointBase) SetName(name string) {
	self.name = name
}

func (self *PointBase) GetTypeName() string {
	return self.typeName
}

func (self *PointBase) SetTypeName(typeName string) {
	self.typeName = typeName
}

// 資料來源 基礎物件
type RefBase struct {
	TypeFrom types.TypeFrom
}
