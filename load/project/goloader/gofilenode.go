package goloader

import "os"

type ANAT1 struct {
	M string
}
type GoFileNode struct {
	path     string
	FileType string
	file     os.FileInfo
	parent   *GoFileNode
	Childes  []*GoFileNode
}

func (self *GoFileNode) Path() string {
	return self.path
}
func (self *GoFileNode) Name() string {
	return self.file.Name()
}
