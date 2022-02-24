package dao

import (
	"codeanalysis/types"
)

// 通用檔案節點
type FileDataNode interface {
	Path() string
	Name() string
}
type TypeBase interface {
	GetName() string
	SetName(name string)
	GetTypeName() string
	SetTypeName(typeName string)
}

// 全指標共用基底

// 名稱 基礎物件
type PointBase struct {
	Name     string
	TypeName string
}

func (self *PointBase) GetName() string {
	return self.Name
}

func (self *PointBase) SetName(name string) {
	self.Name = name
}

func (self *PointBase) GetTypeName() string {
	return self.TypeName
}

func (self *PointBase) SetTypeName(typeName string) {
	self.TypeName = typeName
}

func NewPointBase() *PointBase {
	return &PointBase{}
}

// 資料來源 基礎物件
type RefBase struct {
	TypeFrom types.TypeFrom
}
