package dao

type ModuleInfo struct {
	ModuleName       string // module 名稱
	GoVersion        string // go 最低版本
	CurrentFileNodes FileDataNode

	VendorMap map[string]*PackageInfo // 外部參考 package <packagePath, *PackageInfo>
}

func NewModuleInfo() *ModuleInfo {
	return &ModuleInfo{
		VendorMap: make(map[string]*PackageInfo),
	}
}
