package dao

import "codeanalysis/load/project/goloader"

func NewPackageInfo(name string) *PackageInfo {
	info := &PackageInfo{
		AllTypeInfos: make(map[string]ITypeInfo),
		// StructInfo:   make(map[string]*StructInfo),
		FuncInfo:   make(map[string]*FuncInfo),
		ImportLink: make(map[string]*PackageLink),
	}
	info.Name = name
	return info
}

func NewPackageInfoByNode(node *goloader.GoFileNode) *PackageInfo {
	info := &PackageInfo{
		AllTypeInfos: make(map[string]ITypeInfo),
		// StructInfo:   make(map[string]*StructInfo),
		FuncInfo:   make(map[string]*FuncInfo),
		ImportLink: make(map[string]*PackageLink),
	}
	info.FileNodes = node
	return info
}

// Package 節點
type PackageInfo struct {
	PointBase
	FileNodes    FileDataNode
	AllTypeInfos map[string]ITypeInfo
	// StructInfo   map[string]*StructInfo
	FuncInfo   map[string]*FuncInfo
	ImportLink map[string]*PackageLink
}

func (self *PackageInfo) GetPackageType(packageName, typeName string) (*PackageLink, ITypeInfo) {
	link, ok := self.ImportLink[packageName]
	if !ok {
		panic("")
	}

	typeInfo := link.Package.GetType(typeName)
	return link, typeInfo
}

func (self *PackageInfo) ExistType(typeName string) bool {
	_, ok := self.AllTypeInfos[typeName]
	return ok
}

func (self *PackageInfo) GetType(typeName string) ITypeInfo {
	typeInfo, ok := self.AllTypeInfos[typeName]
	if !ok {
		structInfo := NewTypeInfoStruct()
		structInfo.SetName(typeName)
		typeInfo = structInfo
		self.AllTypeInfos[typeName] = structInfo
	}

	return typeInfo
}
