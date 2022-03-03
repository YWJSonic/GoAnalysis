package dao

import (
	"codeanalysis/analysis/constant"
	"codeanalysis/load/project/goloader"
	"encoding/json"
	"fmt"
)

func NewPackageInfo() *PackageInfo {
	info := &PackageInfo{
		TypeBase:      NewPointBase(),
		AllTypeInfos:  make(map[string]*TypeInfo),
		AllVarInfos:   make(map[string]*VarInfo),
		AllConstInfos: make(map[string]*ConstInfo),
		AllFuncInfo:   make(map[string]*FuncInfo),
		AllImportLink: make(map[string]*ImportInfo),
	}
	return info
}

func NewPackageInfoByNode(node *goloader.GoFileNode) *PackageInfo {
	info := &PackageInfo{
		TypeBase:                  NewPointBase(),
		ImplicitlyVarOrConstInfos: make([]ITypeInfo, 0),
		AllTypeInfos:              make(map[string]*TypeInfo),
		AllVarInfos:               make(map[string]*VarInfo),
		AllConstInfos:             make(map[string]*ConstInfo),
		AllFuncInfo:               make(map[string]*FuncInfo),
		AllImportLink:             make(map[string]*ImportInfo),
	}
	info.CurrentFileNodes = node
	return info
}

// Package 節點
type PackageInfo struct {
	TypeBase
	PackageVersion   string
	GoPath           string
	CurrentFileNodes FileDataNode `json:"-"`

	ImplicitlyVarOrConstInfos []ITypeInfo // 隱藏式宣告參數 var or const
	AllTypeInfos              map[string]*TypeInfo
	AllVarInfos               map[string]*VarInfo
	AllConstInfos             map[string]*ConstInfo
	AllFuncInfo               map[string]*FuncInfo
	AllImportLink             map[string]*ImportInfo // <path, *ImportInfo>
}

func (self *PackageInfo) UnmarshalJSON(b []byte) error {
	tmp := struct {
		TypeBase
		PackageVersion string
		AllTypeInfos   map[string]*TypeInfo
	}{}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	self.TypeBase = tmp.TypeBase
	self.PackageVersion = tmp.PackageVersion
	self.AllTypeInfos = tmp.AllTypeInfos

	return nil
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
	// importLink.SetGoPath(self.GoPath)
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
	iTypeInfo, ok := self.AllTypeInfos[typeName]
	if !ok {
		typeinfo := NewTypeDef()
		typeinfo.SetGoPath(self.GoPath)
		typeinfo.SetTypeFrom(constant.From_Local)
		typeinfo.SetName(typeName)
		self.AllTypeInfos[typeName] = typeinfo
		iTypeInfo = typeinfo
	}

	return iTypeInfo
}
