package dao

// 通用檔案節點
type FileDataNode interface {
	Path() string
	Name() string
}

// 全指標共用基底
type PointBase struct {
	Name string
}
