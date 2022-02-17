package dao

import (
	"codeanalysis/load/project/goloader"
	"fmt"
)

func NewPackageInfo() *PackageInfo {
	info := &PackageInfo{
		AllTypeInfos:  make(map[string]ITypeInfo),
		AllVarInfos:   make(map[string]ITypeInfo),
		AllConstInfos: make(map[string]ITypeInfo),
		AllFuncInfo:   make(map[string]*FuncInfo),
		AllImportLink: make(map[string]*ImportInfo),
	}
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
	info.CurrentFileNodes = node
	return info
}

// Package 節點
type PackageInfo struct {
	PointBase
	CurrentFileNodes FileDataNode
	AllTypeInfos     map[string]ITypeInfo
	AllVarInfos      map[string]ITypeInfo
	AllConstInfos    map[string]ITypeInfo
	AllFuncInfo      map[string]*FuncInfo
	AllImportLink    map[string]*ImportInfo // <path, *ImportInfo>
}

func (self *PackageInfo) GetPackage(packageName string) *ImportInfo {
	for _, importLink := range self.AllImportLink {
		if importLink.NewName == packageName {
			return importLink
		}
	}

	// package 不存在
	// packageName 與 path 不同
	// import 未重新命名
	importLink := NewImportLink()
	importLink.NewName = packageName
	self.AllImportLink[fmt.Sprintf("%s_%s", "unknow", packageName)] = importLink
	return importLink
}

func (self *PackageInfo) GetPackageType(packageName, typeName string) (*ImportInfo, ITypeInfo) {
	// link, ok := self.AllImportLink[packageName]
	// if !ok {
	// 	panic("")
	// }

	link := self.GetPackage(packageName)

	// golang 內建 package 處理
	if link.Package == nil {
		return link, nil
	}

	// 一般 package
	typeInfo := link.Package.GetType(typeName)
	return link, typeInfo
}

func (self *PackageInfo) ExistType(typeName string) bool {
	_, ok := self.AllTypeInfos[typeName]
	return ok
}

func (self *PackageInfo) GetType(typeName string) ITypeInfo {
	if typeName == "" {
		panic("typeName empty")
	}
	if self == nil {
		panic("type not find")
	}
	typeInfo, ok := self.AllTypeInfos[typeName]
	if !ok {
		typeInfo = NewTypeDef()
		typeInfo.SetName(typeName)
		self.AllTypeInfos[typeName] = typeInfo
	}

	return typeInfo
}
