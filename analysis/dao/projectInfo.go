package dao

import "codeanalysis/load/project/goloader"

func NewProjectInfo(node *goloader.GoFileNode) *ProjectInfo {
	info := &ProjectInfo{
		AllPackageMap: make(map[string]*PackageInfo),
		// AllStructMap:     make(map[string]*StructInfo),
		// AllPublicFuncMap: make(map[string]*FuncInfo),
	}
	info.name = node.Path()
	info.ProjectRoot = node
	return info
}

type ProjectInfo struct {
	PointBase
	ProjectRoot   FileDataNode
	AllPackageMap map[string]*PackageInfo // <packagePath, *PackageInfo>
	// AllStructMap     map[string]*StructInfo
	// AllPublicFuncMap map[string]*FuncInfo
	// AllVendorMap     map[string]string
}

func (self *ProjectInfo) LoadOrStoryPackage(info *PackageInfo) *PackageInfo {

	if packageInfo, ok := self.AllPackageMap[info.FileNodes.Path()]; ok {
		return packageInfo
	}

	self.AllPackageMap[info.FileNodes.Path()] = info
	return info
}

func (self *ProjectInfo) LorSPackageInfo(name string) (*PackageInfo, bool) {
	packageInfo, isLoad := self.AllPackageMap[name]
	if isLoad {
		return packageInfo, isLoad
	}
	packageInfo = NewPackageInfo(name)
	return packageInfo, isLoad
}
