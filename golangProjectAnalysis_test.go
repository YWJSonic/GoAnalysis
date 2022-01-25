package main

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/analysis/goanalysis"
	"codeanalysis/load/project/goloader"
	"testing"
)

func TestGolandAnalysis(t *testing.T) {
	projectRootNode := goloader.LoadRoot("./structTest/")
	goanalysis.Instants = dao.NewProjectInfo(projectRootNode)
	goanalysis.GoAnalysisSpaceFirst(projectRootNode)
}
