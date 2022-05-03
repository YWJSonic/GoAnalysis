package goloader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 建立專案節點樹
func LoadRoot(rootPath string) *GoFileNode {
	rootNode := newNode(rootPath, nil)

	option := Option{
		MaxDepth: 0,
		SkipFile: []string{".vscode", ".git", "vendor", "doc"},
	}
	rootNode.Childes = loopBuildStruct(rootNode, 0, option)
	return rootNode
}

// 建立參數
type Option struct {
	MaxDepth int      // 專案讀取深度
	SkipFile []string // 指定隱藏檔案
}

// 建立目標節點的子節點樹
func loopBuildStruct(node *GoFileNode, depth int, option Option) []*GoFileNode {
	var nodeChilds []*GoFileNode

	// 檔案深度遞增
	depth++
	files, _ := ioutil.ReadDir(node.path)

	// 建立子檔案節點
	for _, file := range files {
		isSkip := false
		for _, skipFile := range option.SkipFile {
			if skipFile == file.Name() {
				isSkip = true
			}
		}
		if isSkip {
			continue
		}

		// 建立檔案節點
		newPath := filepath.Join(node.path, file.Name())
		newNode := newNode(newPath, file)
		nodeChilds = append(nodeChilds, newNode)

		// 取得子節點的 子節點
		newNode.Childes = loopBuildStruct(newNode, depth, option)
	}

	return nodeChilds
}

// 建立新 go 檔案節點
func newNode(path string, file os.FileInfo) *GoFileNode {
	return &GoFileNode{
		path:     path,
		FileType: fileType(path),
		file:     file,
	}
}

// 檢查檔案類型
func fileType(fileName string) string {
	splitStr := strings.Split(fileName, ".")
	if len(splitStr) < 2 { // 無副檔名
		return ""
	}
	return splitStr[len(splitStr)-1]
}
