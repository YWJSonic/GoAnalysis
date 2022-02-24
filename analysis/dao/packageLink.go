package dao

import "fmt"

func NewImportLink() *ImportInfo {
	return &ImportInfo{}
}

// Package 關聯(import 資料)
type ImportInfo struct {
	Package   *PackageInfo `json:"-"`
	Path      string
	NewName   string // 使用一般引入方式
	ImportMod string // 包引入方式定義, 一般使用： "", 隱藏式: ".", 只限初始化: "_"
}

func (self *ImportInfo) GetName() string {
	if self.ImportMod == "_" || self.ImportMod == "." {
		return fmt.Sprintf("%s%s", self.ImportMod, self.NewName)
	}
	return self.NewName
}
