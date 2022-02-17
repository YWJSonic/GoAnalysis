package main

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/analysis/goanalysis"
	"codeanalysis/load/project/goloader"
	"fmt"
	"testing"

	_ "net/http/pprof"
)

func TestGolandAnalysis(t *testing.T) {
	// go http.ListenAndServe("0.0.0.0:6060", nil)
	projectRootNode := goloader.LoadRoot("/home/yang/Desktop/GameBackend/gameservice/")
	goanalysis.Instants = dao.NewProjectInfo("GameService", projectRootNode)
	goanalysis.GoAnalysisSpaceFirst(projectRootNode)
	fmt.Println(goanalysis.Instants)
}
