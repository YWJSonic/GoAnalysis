package goanalysis

import (
	_ "codeanalysis/load/project/goloader"
	wwe "codeanalysis/load/project/goloader"
	"codeanalysis/tool"
	"fmt"
)

//專案節點樹分析 廣度優先
func GoAnalysisSpaceFirst(node *wwe.GoFileNode) {
	var childs []*wwe.GoFileNode

	//取得節點下的全部子節點
	getSubChiles := func(node *wwe.GoFileNode) []*wwe.GoFileNode {
		var childs []*wwe.GoFileNode
		for _, child := range node.Childes {
			if child.Childes != nil {
				childs = append(childs, child.Childes...)
			}
		}
		return childs
	}

	childs = node.Childes
	var currentChild *wwe.GoFileNode
	// 要輪尋的節點尚未
	for len(childs) > 0 {

		currentChild, childs = childs[0], childs[1:]

		if currentChild.Childes != nil { // 節點下還有子節點
			childs = append(childs, getSubChiles(currentChild)...)

		} else if currentChild.FileType == "go" { // 檔案節點
			PackageInfo := Instants.LoadOrStoryPackage(NewPackageInfoByNode(currentChild))
			// 讀檔
			code := tool.ReadFile(PackageInfo.FileNodes.Path())
			// 檔案分析
			fmt.Println("------------------------")
			AnalysisStyle2(PackageInfo, code)
			fmt.Println("------------------------")
		}
	}
}

// 專案節點樹分析  深度優先
func GoAnalysisDepthFirst(node *wwe.GoFileNode) {
	// 節點下還有子節點
	for _, child := range node.Childes {
		GoAnalysisDepthFirst(child)
	}

	if node.FileType == "go" { // 檔案節點
		// fmt.Println(node.Name())
		// 讀檔
		code := tool.ReadFile(node.Path())
		// 檔案分析
		AnalysisStyle1(code)
	}
}
