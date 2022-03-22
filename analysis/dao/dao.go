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

	// 此物件名稱, 宣告物件 為自定義名稱
	GetName() string
	SetName(name string)

	// 此 Type Info 類型名稱
	GetTypeName() string
	SetTypeName(typeName string)

	// go mod 使用的路徑
	SetGoPath(goPath string)

	// 用來判斷 package 來源(目前不准)
	SetTypeFrom(from types.TypeFrom)

	// 旗標 判斷是否已被解析
	GetIsAnalysis() bool
	SetIsAnalysis()
}

// 全指標共用基底

// 名稱 基礎物件
type PointBase struct {
	Name       string
	TypeName   string
	TypeFrom   types.TypeFrom
	GoPath     string
	IsAnalysis bool
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

func (self *PointBase) SetGoPath(goPath string) {
	self.GoPath = goPath
}

func (self *PointBase) SetTypeFrom(from types.TypeFrom) {
	self.TypeFrom = from
}

func (self *PointBase) GetIsAnalysis() bool {
	return self.IsAnalysis
}

func (self *PointBase) SetIsAnalysis() {
	self.IsAnalysis = true
}

func NewPointBase() *PointBase {
	return &PointBase{}
}
