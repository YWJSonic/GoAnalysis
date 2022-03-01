package dao

import (
	"codeanalysis/analysis/constant"
	"codeanalysis/load/project/goloader"
	"codeanalysis/types"
)

func NewProjectInfo(node *goloader.GoFileNode) *ProjectInfo {
	info := &ProjectInfo{
		TypeBase:         NewPointBase(),
		LocalPackageMap:  make(map[string]*PackageInfo),
		VendorPackageMap: make(map[string]*PackageInfo),
		GolangPackageMap: make(map[string]*PackageInfo),
	}
	info.SetName(node.Path())
	info.ProjectRoot = node
	return info
}

type ProjectInfo struct {
	TypeBase
	ModuleInfo       *ModuleInfo
	ProjectRoot      FileDataNode            `json:"-"`
	LocalPackageMap  map[string]*PackageInfo // 內部實做 package <packagePath, *PackageInfo>
	VendorPackageMap map[string]*PackageInfo // 外部引用 package <packagePath, *PackageInfo>
	GolangPackageMap map[string]*PackageInfo // 系統自帶 package <packagePath, *PackageInfo>
}

// 讀寫此 package 關連資料
//
// @params string		package 路徑
// @params *PackageInfo	預設關聯結構
//
// @return *PackageInfo	回傳的關聯資料
// @return bool			是否已存在資料
func (self *ProjectInfo) LoadOrStoryPackage(packageType types.TypeFrom, pwd string, info *PackageInfo) (*PackageInfo, bool) {

	switch packageType {
	case constant.From_Golang:
		if packageInfo, ok := self.GolangPackageMap[pwd]; ok {
			return packageInfo, true
		}
		self.GolangPackageMap[pwd] = info

	case constant.From_Local:
		if packageInfo, ok := self.LocalPackageMap[pwd]; ok {
			return packageInfo, true
		}
		self.LocalPackageMap[pwd] = info

	case constant.From_Vendor:
		if packageInfo, ok := self.VendorPackageMap[pwd]; ok {
			return packageInfo, true
		}
		self.VendorPackageMap[pwd] = info
	}

	return info, false
}
