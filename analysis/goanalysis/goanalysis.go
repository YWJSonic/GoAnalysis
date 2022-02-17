package goanalysis

import (
	"bytes"
	"codeanalysis/analysis/dao"
	_ "codeanalysis/load/project/goloader"
	wwe "codeanalysis/load/project/goloader"
	"codeanalysis/util"
	"fmt"
	"strings"
)

var Instants *dao.ProjectInfo

//專案節點樹分析 廣度優先
func GoAnalysisSpaceFirst(node *wwe.GoFileNode) {
	var childs []*wwe.GoFileNode

	// 取得節點下的全部子節點方法
	// getSubChiles := func(node *wwe.GoFileNode) []*wwe.GoFileNode {
	// 	var childs []*wwe.GoFileNode
	// 	for _, child := range node.Childes {
	// 		if child.Childes != nil {
	// 			childs = append(childs, child.Childes...)
	// 		}
	// 	}
	// 	return childs
	// }

	// 根節點如果是 go 檔直接加入解析
	if node.FileType == "go" {
		childs = append(childs, node)
	} else {
		childs = node.Childes
	}

	// 要輪尋的節點尚未
	var currentChild *wwe.GoFileNode
	for len(childs) > 0 {
		currentChild, childs = childs[0], childs[1:]
		if currentChild.FileType != "go" {
			// go 檔以外的往下執行
			GoAnalysisSpaceFirst(currentChild)
		} else if currentChild.FileType == "go" { // 檔案節點

			// 檔案路徑與 import 路徑無法批配
			// 未讀取 go.mod 檔
			// main 檔層級路徑待確認
			if node.Path() == Instants.ProjectRoot.Path() {
				// panic("node is go file")
			}
			idx := strings.Index(node.Path(), Instants.ProjectRoot.Path())
			pwd := node.Path()
			if idx == 0 {
				pwd = pwd[len(Instants.ProjectRoot.Path()):]
				pwd = Instants.ModuleName + pwd
			} else if idx > 0 {
				panic("path root error")
			}

			// 確認是否有已生成的資料
			PackageInfo, ok := Instants.LoadOrStoryPackage(pwd, dao.NewPackageInfoByNode(currentChild))
			if ok {
				PackageInfo.CurrentFileNodes = currentChild
			}
			// 讀檔
			code := util.ReadFile(PackageInfo.CurrentFileNodes.Path())

			// 檔案分析
			AnalysisStyle2(PackageInfo, code)
			PackageInfo.CurrentFileNodes = nil
		}
	}
}

// 專案節點樹分析  深度優先
// func GoAnalysisDepthFirst(node *wwe.GoFileNode) {
// 	// 節點下還有子節點
// 	for _, child := range node.Childes {
// 		GoAnalysisDepthFirst(child)
// 	}

// 	if node.FileType == "go" { // 檔案節點
// 		// fmt.Println(node.Name())
// 		// 讀檔
// 		code := util.ReadFile(node.Path())
// 		// 檔案分析
// 		AnalysisStyle1(code)
// 	}
// }

// 語言分析
// func AnalysisStyle1(code string) {
// 	// var package map[string]*PackageInfo
// 	buf := bytes.NewBuffer([]byte(" "))
// 	buf.WriteString(code)
// 	s := source{
// 		buf: buf.Bytes(),
// 	}
// 	s.start()
// 	for {
// 		if s.r == s.e {
// 			break
// 		}

// 		s.next()
// 		switch s.rangeStr() {
// 		case "package":
// 			s.next()
// 			s.PackageInfo.Name = s.rangeStr()
// 		case "import":
// 			s.importSpec()
// 		case "type":
// 			s.TypeDeclarations()
// 		case "func":
// 			s.FunctionDeclarations()
// 		case "const":
// 			s.ConstantDeclarations()
// 		case "var":
// 			s.VariableDeclarations()
// 		default:
// 			fmt.Println("--- ", s.rangeStr())
// 		}
// 	}
// }

func AnalysisStyle2(packageInfo *dao.PackageInfo, code string) {
	// var package map[string]*PackageInfo
	buf := bytes.NewBuffer([]byte(" "))
	buf.WriteString(code)
	s := source{
		buf:         buf.Bytes(),
		PackageInfo: packageInfo,
	}

	if s.PackageInfo.CurrentFileNodes.Path() == "/home/yang/Desktop/GameBackend/gameservice/parties/party1/example/flow_trace.go" {
		fmt.Println("")
	}

	s.start()
	for {
		s.toNextCh()
		if s.r+1 == s.e {
			break
		}

		// 最外層註解
		if s.CheckCommon() {
			s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		} else {

			s.next()
			switch s.rangeStr() {
			case "package":
				s.next()
				s.PackageInfo.SetName(s.rangeStr())
			case "import":
				s.ImportDeclarations()

			case "type":
				s.TypeDeclarations()

			case "func":
				if s.buf[s.r+1] == '(' {
					s.MethodDeclarations()
				} else {
					s.FunctionDeclarations()
				}
			case "const":
				s.ConstantDeclarations()
			case "var":
				s.VariableDeclarations()

			default:
			}
		}
	}
	// for name, info := range s.PackageInfo.AllImportLink {
	// 	fmt.Println(name, info.Path)
	// }
	// for name, info := range s.PackageInfo.AllConstInfos {
	// 	fmt.Println(name, info.GetTypeName())
	// }
	// for name, info := range s.PackageInfo.AllVarInfos {
	// 	fmt.Println(name, info.GetTypeName())
	// }
	// for name, info := range s.PackageInfo.AllTypeInfos {
	// 	fmt.Println(name, info.GetTypeName())
	// }
	// for name, info := range s.PackageInfo.AllFuncInfo {
	// 	fmt.Println(name, info.GetTypeName())
	// }
}
