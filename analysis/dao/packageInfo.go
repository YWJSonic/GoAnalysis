package dao

import (
	"codeanalysis/analysis/constant"
	"codeanalysis/load/project/goloader"
	"encoding/json"
	"fmt"
)

func NewPackageInfo() *PackageInfo {
	info := &PackageInfo{
		TypeBase:        NewPointBase(),
		UndefVarOrConst: make(map[string]ITypeInfo),
		AllTypeInfos:    make(map[string]*TypeInfo),
		AllVarInfos:     make(map[string]*VarInfo),
		AllConstInfos:   make(map[string]*ConstInfo),
		AllFuncInfo:     make(map[string]*FuncInfo),
		AllImportLink:   make(map[string]*ImportInfo),
	}
	return info
}

func NewPackageInfoByNode(node *goloader.GoFileNode) *PackageInfo {
	info := &PackageInfo{
		TypeBase:                  NewPointBase(),
		ImplicitlyVarOrConstInfos: make([]ITypeInfo, 0),
		UndefVarOrConst:           make(map[string]ITypeInfo),
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

	// var, const 解析完成前儲存於此
	// 當該 package 解析完成最後在分類
	// 或是就不分了
	AllVarAndConst []ITypeInfo

	UndefVarOrConst           map[string]ITypeInfo
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

	fmt.Printf("GetPackage Error packageNAme: %s, GoPath: %s", packageName, self.GoPath)
	panic("")
}

func (self *PackageInfo) GetPackageType(packageName, typeName string) (*ImportInfo, ITypeInfo) {
	link := self.GetPackage(packageName)

	// golang 內建 package 處理
	if link.Package == nil {
		return link, nil
	}

	// 一般 package
	typeInfo := link.Package.GetType(typeName)
	return link, typeInfo
}
func (self *PackageInfo) GetPackageFunc(packageName, funcName string) (*ImportInfo, ITypeInfo) {
	link := self.GetPackage(packageName)

	// golang 內建 package 處理
	if link.Package == nil {
		return link, nil
	}

	// 一般 package
	typeInfo := link.Package.GetFunc(funcName)
	return link, typeInfo
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
		typeInfo := NewTypeDef()
		typeInfo.SetGoPath(self.GoPath)
		typeInfo.SetTypeFrom(constant.From_Local)
		typeInfo.SetName(typeName)
		self.AllTypeInfos[typeName] = typeInfo
		iTypeInfo = typeInfo
	}

	return iTypeInfo
}
func (self *PackageInfo) GetFunc(funcName string) ITypeInfo {
	if funcName == "" {
		panic("typeName empty")
	}
	if self == nil {
		panic("type not find")
	}
	iTypeInfo, ok := self.AllFuncInfo[funcName]
	if !ok {
		funcInfo := NewFuncInfo()
		funcInfo.SetGoPath(self.GoPath)
		funcInfo.SetTypeFrom(constant.From_Local)
		funcInfo.SetName(funcName)
		self.AllFuncInfo[funcName] = funcInfo
		iTypeInfo = funcInfo
	}

	return iTypeInfo
}
