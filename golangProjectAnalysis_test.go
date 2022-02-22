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

	fmt.Println(goanalysis.Instants.Output())

	// 測試輸出
	// rebuildCode()
}

func rebuildCode() {
	for key, packageInfo := range goanalysis.Instants.LocalPackageMap {
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
					fmt.Println("\t", varKey, "\t\t", varInfo.ContentTypeInfo.GetTypeName())
				}
				// for funcKey, funcInfo := range info.FuncPoint {
				// 	fmt.Println("\t", funcKey, "\t\t", funcInfo.GetName())
				// }
			}
		}
		fmt.Println("========= Func =============")
		for _, funcInfo := range packageInfo.AllFuncInfo {
			funcInfo.Print()
		}
	}
}
