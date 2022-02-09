package dao

func NewImportLink(name, path string, targetPackage *PackageInfo) *ImportInfo {
	return &ImportInfo{
		NewName: name,
		Path:    path,
		Package: targetPackage,
	}
}

// Package 關聯(import 資料)
type ImportInfo struct {
	Package *PackageInfo
	Path    string
	NewName string
}

func (self *ImportInfo) GetName() string {
	if self.NewName == "" || self.NewName == "_" { // 未明新命名或隱藏式 import
		return self.Package.name
	}
	return self.NewName
}
