package goanalysis

import "codeanalysis/load/project/goloader"

var Instants *ProjectInfo

func NewProjectInfo(node *goloader.GoFileNode) *ProjectInfo {
	info := &ProjectInfo{
		AllPackageMap:    make(map[string]*PackageInfo),
		AllStructMap:     make(map[string]*StructInfo),
		AllPublicFuncMap: make(map[string]*FuncInfo),
	}
	info.Name = node.Path()
	info.ProjectRoot = node
	return info
}

type ProjectInfo struct {
	PointBase
	ProjectRoot      FileDataNode
	AllPackageMap    map[string]*PackageInfo
	AllStructMap     map[string]*StructInfo
	AllPublicFuncMap map[string]*FuncInfo
	// AllVendorMap     map[string]string
}

func (self *ProjectInfo) LoadOrStoryPackage(info *PackageInfo) *PackageInfo {

	if packageInfo, ok := self.AllPackageMap[info.FileNodes.Path()]; ok {
		return packageInfo
	}

	self.AllPackageMap[info.FileNodes.Path()] = info
	return info
}
