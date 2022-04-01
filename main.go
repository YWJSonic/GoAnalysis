package main

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/analysis/goanalysis"
	"codeanalysis/graphics"
	"codeanalysis/load/project/goloader"
)

func main() {
	RunAnalysis()
}

func RunAnalysis() {
	// go http.ListenAndServe("0.0.0.0:6060", nil)
	projectRootNode := goloader.LoadRoot("/home/yang/Desktop/GoAnalysis")

	goanalysis.Instants = dao.NewProjectInfo(projectRootNode)

	// 解析 go.mod
	for _, child := range projectRootNode.Childes {
		if child.Name() == "go.mod" {
			goanalysis.GoAnalysisGoMod(child)
			break
		}
	}

	// 解析專案
	goanalysis.GoAnalysisSpaceFirst(projectRootNode)
	// 解析表達式 第二階段
	// goanalysis.GoAnalysisExpression(goanalysis.Instants)
	// 輸出圖形化
	graphics.OutputPlantUmlGraphics(goanalysis.Instants)

	// data, err := json.Marshal(goanalysis.Instants)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// ioutil.WriteFile("./log.json", data, 0666)
}
