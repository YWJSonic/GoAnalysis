package dao

import "codeanalysis/load/project/goloader"

func NewPackageInfo(name string) *PackageInfo {
	info := new(PackageInfo)
	info.Name = name
	return info
}

func NewPackageInfoByNode(node *goloader.GoFileNode) *PackageInfo {
	info := &PackageInfo{
		StructInfo: make(map[string]*StructInfo),
		FuncInfo:   make(map[string]*FuncInfo),
		ImportLink: make(map[string]*PackageLink),
	}
	info.FileNodes = node
	return info
}

// Package 節點
type PackageInfo struct {
	PointBase
	FileNodes  FileDataNode
	StructInfo map[string]*StructInfo
	FuncInfo   map[string]*FuncInfo
	ImportLink map[string]*PackageLink
}

func (self *PackageInfo) LorSStructInfo(name string) (*StructInfo, bool) {
	structInfo, isLoad := self.StructInfo[name]
	if isLoad {
		return structInfo, isLoad
	}

	structInfo = NewStructInfo(name)
	return structInfo, isLoad
}
