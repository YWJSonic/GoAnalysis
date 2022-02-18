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
	projectRootNode := goloader.LoadRoot("/home/yang/Desktop/GameBackend/scheduler")
	goanalysis.Instants = dao.NewProjectInfo("gitlab.geax.io/demeter/backend/scheduler", projectRootNode)
	goanalysis.GoAnalysisSpaceFirst(projectRootNode)
	rebuildCode()
}

func rebuildCode() {
	for key, packageInfo := range goanalysis.Instants.AllPackageMap {
		fmt.Println("== Package", packageInfo.GetName(), " ==========================", key)

		fmt.Println("========= Import Link =============")
		for importKey, importInfo := range packageInfo.AllImportLink {
			fmt.Println("import", importKey, importInfo)
		}
		fmt.Println("========= Const =============")
		for constKey, constInfo := range packageInfo.AllConstInfos {
			fmt.Println("const", constKey, constInfo.(*dao.ConstInfo).Expressions)
		}
		fmt.Println("========= Var =============")
		for varKey, varInfo := range packageInfo.AllVarInfos {
			fmt.Println("var", varKey, varInfo.(*dao.VarInfo).Expressions)
		}
		fmt.Println("========= Type =============")
		for typeKey, typeInfo := range packageInfo.AllTypeInfos {
			fmt.Println("type", typeKey)
			constInfo := typeInfo.(*dao.TypeInfo)
			switch info := constInfo.ContentTypeInfo.(type) {
			case *dao.TypeInfoStruct:
				for varKey, varInfo := range info.VarInfos {
					fmt.Println("\t", varKey, "\t\t", varInfo.TypeInfo.GetTypeName())
				}
				for funcKey, funcInfo := range info.FuncPoint {
					fmt.Println("\t", funcKey, "\t\t", funcInfo.GetName())
				}
			}
		}
		fmt.Println("========= Func =============")
		for _, funcInfo := range packageInfo.AllFuncInfo {
			funcInfo.Print()
		}
	}
}
