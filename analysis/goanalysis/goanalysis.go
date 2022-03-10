package goanalysis

import (
	"bytes"
	"codeanalysis/analysis/constant"
	"codeanalysis/analysis/dao"
	_ "codeanalysis/load/project/goloader"
	wwe "codeanalysis/load/project/goloader"
	"codeanalysis/util"
	"strings"
)

var Instants *dao.ProjectInfo

//專案節點樹分析 深度優先
func GoAnalysisSpaceFirst(node *wwe.GoFileNode) {
	var childs []*wwe.GoFileNode

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
				pwd = Instants.ModuleInfo.ModuleName + pwd
			} else if idx > 0 {
				panic("path root error")
			}

			// 確認是否有已生成的資料
			PackageInfo, ok := Instants.LoadOrStoryPackage(constant.From_Local, pwd, dao.NewPackageInfoByNode(currentChild))
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

func AnalysisStyle2(packageInfo *dao.PackageInfo, code string) {
	// var package map[string]*PackageInfo
	buf := bytes.NewBuffer([]byte(" "))
	buf.WriteString(code)
	s := source{
		buf:         buf.Bytes(),
		PackageInfo: packageInfo,
	}

	s.start()
	for {
		s.toNextCh()

		if s.r+1 == s.e {
			break
		}

		// 最外層註解
		if s.checkCommon() {
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
}
