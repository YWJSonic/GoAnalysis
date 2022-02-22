package dao

import "strings"

type ModuleInfo struct {
	ModuleName       string // module 名稱
	GoVersion        string // go 最低版本
	CurrentFileNodes FileDataNode

	VendorMap map[string]struct{} // 外部參考 package <packagePath, *PackageInfo>
}

func (self *ModuleInfo) IsLocalPackage(path string) bool {
	return strings.Index(path, self.ModuleName) == 0
}

func (self *ModuleInfo) IsVendorPackage(path string) bool {
	for vendorPath := range self.VendorMap {
		if strings.Index(path, vendorPath) == 0 {
			return true
		}
	}
	return false
}

func NewModuleInfo() *ModuleInfo {
	return &ModuleInfo{
		VendorMap: make(map[string]struct{}),
	}
}
