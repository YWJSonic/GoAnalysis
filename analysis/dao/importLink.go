package dao

import (
	"encoding/json"
	"fmt"
)

func NewImportLink() *ImportInfo {
	return &ImportInfo{
		TypeBase: NewPointBase(),
	}
}

// Package 關聯(import 資料)
type ImportInfo struct {
	TypeBase
	Package   *PackageInfo `json:"-"`
	NewName   string       // 使用一般引入方式
	ImportMod string       // 包引入方式定義, 一般使用： "", 隱藏式: ".", 只限初始化: "_"
}

func (self *ImportInfo) MarshalJSON() ([]byte, error) {
	tmp := struct {
		TypeBase
		NewName     string
		PackagePath string
		ImportMod   string
	}{
		TypeBase:    self.TypeBase,
		PackagePath: self.Package.GoPath,
		NewName:     self.NewName,
		ImportMod:   self.ImportMod,
	}

	return json.Marshal(tmp)
}

func (self *ImportInfo) GetName() string {
	if self.ImportMod == "_" || self.ImportMod == "." {
		return fmt.Sprintf("%s%s", self.ImportMod, self.NewName)
	}
	return self.NewName
}
