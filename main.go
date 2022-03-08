package main

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/analysis/goanalysis"
	"codeanalysis/graphics"
	"codeanalysis/load/project/goloader"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	RunAnalysis()
}

func RunAnalysis() {
	// go http.ListenAndServe("0.0.0.0:6060", nil)
	projectRootNode := goloader.LoadRoot("/home/yang/Desktop/GameBackend/gamemaster")

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
	graphics.OutputPlantUmlGraphics(goanalysis.Instants)

	data, err := json.Marshal(goanalysis.Instants)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile("./log.json", data, 0666)
}
