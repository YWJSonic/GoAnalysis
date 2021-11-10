package goanalysis

import "codeanalysis/load/project/goloader"

func NewPackageInfo(name string) *PackageInfo {
	info := new(PackageInfo)
	info.Name = name
	return info
}

func NewPackageInfoByNode(node *goloader.GoFileNode) *PackageInfo {
	info := new(PackageInfo)
	info.FileNodes = node
	return info
}

// Package 節點
type PackageInfo struct {
	PointBase
	FileNodes   FileDataNode
	StructPoint map[string]*StructInfo
	FuncPoint   map[string]*FuncInfo
	ImportLink  map[string]*PackageLink
}

// 通用檔案節點
type FileDataNode interface {
	Path() string
	Name() string
}

func NewPackageLink(name string, targetPackage *PackageInfo) *PackageLink {
	link := &PackageLink{
		NewName: name,
		Package: targetPackage,
	}
	return link
}

// Package 關聯(import 資料)
type PackageLink struct {
	Package *PackageInfo
	NewName string
}

func (self *PackageLink) Name() string {
	if self.NewName == "" || self.NewName == "_" { // 未明新命名或隱藏式 import
		return self.Package.Name
	}
	return self.NewName
}

func NewStructInfo(name string) *StructInfo {
	info := &StructInfo{}
	info.Name = name
	return info
}

// Struct 關聯
type StructInfo struct {
	PointBase
	Package   *PackageInfo
	FuncPoint *FuncInfo
}

func NewMethodExpression(name string, structPoint *StructInfo) *MethodExpression {
	params := &MethodExpression{}
	params.Name = name
	params.StructPoint = structPoint
	return params
}

type MethodExpression struct {
	PointBase
	StructPoint *StructInfo
}

func NewFuncInfo(name string) *FuncInfo {
	info := &FuncInfo{}
	info.Name = name
	return info
}

// Func關聯
type FuncInfo struct {
	PointBase
	Method         *MethodExpression
	ParamsInPoint  []*FuncParams
	ParamsOutPoint []*FuncParams
}

func NewFuncParams(name string, structPoint *StructInfo) *FuncParams {
	params := &FuncParams{}
	params.Name = name
	params.StructPoint = structPoint
	return params
}

// Func 傳輸參數
type FuncParams struct {
	PointBase
	StructPoint *StructInfo
}

// 全指標共用基底
type PointBase struct {
	Name string
}
