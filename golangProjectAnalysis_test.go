package main

import (
	"codeanalysis/analysis/goanalysis"
	"codeanalysis/load/project/goloader"
	"testing"
)

func TestGolandAnalysis(t *testing.T) {
	projectRootNode := goloader.LoadRoot("./structTest/")
	goanalysis.Instants = goanalysis.NewProjectInfo(projectRootNode)
	goanalysis.GoAnalysisSpaceFirst(projectRootNode)
}
