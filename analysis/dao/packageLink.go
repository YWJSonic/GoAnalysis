package dao

func NewPackageLink(name, path string, targetPackage *PackageInfo) *PackageLink {
	link := &PackageLink{
		NewName: name,
		Path:    path,
		Package: targetPackage,
	}
	return link
}

// Package 關聯(import 資料)
type PackageLink struct {
	Package *PackageInfo
	Path    string
	NewName string
}

func (self *PackageLink) Name() string {
	if self.NewName == "" || self.NewName == "_" { // 未明新命名或隱藏式 import
		return self.Package.Name
	}
	return self.NewName
}
