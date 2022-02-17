package dao

import (
	"codeanalysis/load/project/goloader"
)

func NewProjectInfo(moduleName string, node *goloader.GoFileNode) *ProjectInfo {
	info := &ProjectInfo{
		AllPackageMap: make(map[string]*PackageInfo),
		// AllStructMap:     make(map[string]*StructInfo),
		// AllPublicFuncMap: make(map[string]*FuncInfo),
	}
	info.ModuleName = moduleName
	info.name = node.Path()
	info.ProjectRoot = node
	return info
}

type ProjectInfo struct {
	PointBase
	ModuleName    string
	ProjectRoot   FileDataNode
	AllPackageMap map[string]*PackageInfo // <packagePath, *PackageInfo>

}

// 讀寫此 package 關連資料
//
// @params string		package 路徑
// @params *PackageInfo	預設關聯結構
//
// @return *PackageInfo	回傳的關聯資料
// @return bool			是否已存在資料
func (self *ProjectInfo) LoadOrStoryPackage(pwd string, info *PackageInfo) (*PackageInfo, bool) {
	if packageInfo, ok := self.AllPackageMap[pwd]; ok {
		return packageInfo, true
	}

	self.AllPackageMap[pwd] = info
	return info, false
}
