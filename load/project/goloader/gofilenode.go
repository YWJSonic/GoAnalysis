package goloader

import "os"

type ANAT1 struct {
	M string
}
type GoFileNode struct {
	path     string
	FileType string
	file     os.FileInfo   `json:"-"`
	parent   *GoFileNode   `json:"-"`
	Childes  []*GoFileNode `json:"-"`
}

func (self *GoFileNode) Path() string {
	return self.path
}
func (self *GoFileNode) Name() string {
	return self.file.Name()
}
