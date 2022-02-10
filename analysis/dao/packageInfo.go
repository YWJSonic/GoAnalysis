package dao

import "codeanalysis/load/project/goloader"

func NewPackageInfo(name string) *PackageInfo {
	info := &PackageInfo{
		AllTypeInfos:  make(map[string]ITypeInfo),
		AllVarInfos:   make(map[string]ITypeInfo),
		AllConstInfos: make(map[string]ITypeInfo),
		AllFuncInfo:   make(map[string]*FuncInfo),
		AllImportLink: make(map[string]*ImportInfo),
	}
	info.name = name
	return info
}

func NewPackageInfoByNode(node *goloader.GoFileNode) *PackageInfo {
	info := &PackageInfo{
		AllTypeInfos:  make(map[string]ITypeInfo),
		AllVarInfos:   make(map[string]ITypeInfo),
		AllConstInfos: make(map[string]ITypeInfo),
		AllFuncInfo:   make(map[string]*FuncInfo),
		AllImportLink: make(map[string]*ImportInfo),
	}
	info.FileNodes = node
	return info
}

// Package 節點
type PackageInfo struct {
	PointBase
	FileNodes     FileDataNode
	AllTypeInfos  map[string]ITypeInfo
	AllVarInfos   map[string]ITypeInfo
	AllConstInfos map[string]ITypeInfo
	AllFuncInfo   map[string]*FuncInfo
	AllImportLink map[string]*ImportInfo
}

func (self *PackageInfo) GetPackageType(packageName, typeName string) (*ImportInfo, ITypeInfo) {
	link, ok := self.AllImportLink[packageName]
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
		typeInfo = NewTypeDef()
		typeInfo.SetName(typeName)
		self.AllTypeInfos[typeName] = typeInfo
	}

	return typeInfo
}
