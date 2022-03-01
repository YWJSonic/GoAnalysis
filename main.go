package main

import (
	"codeanalysis/analysis/dao"
	"codeanalysis/analysis/goanalysis"
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

	data, err := json.Marshal(goanalysis.Instants)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile("./log.json", data, 0666)
	fmt.Println(string(data))

	// 測試輸出
	rebuildCode()
}

func rebuildCode() {
	for key, packageInfo := range goanalysis.Instants.LocalPackageMap {
		fmt.Println("== Package", packageInfo.GetName(), " ==========================", key)

		// fmt.Println("========= Import Link =============")
		// for importKey, importInfo := range packageInfo.AllImportLink {
		// 	fmt.Println("import", importKey, importInfo)
		// }
		// fmt.Println("========= Const =============")
		// for constKey, constInfo := range packageInfo.AllConstInfos {
		// 	fmt.Println("const", constKey, constInfo.(*dao.ConstInfo).Expression)
		// }
		// fmt.Println("========= Var =============")
		// for varKey, varInfo := range packageInfo.AllVarInfos {
		// 	fmt.Println("var", varKey, varInfo.(*dao.VarInfo).Expression)
		// }
		fmt.Println("========= Type =============")
		for _, typeInfo := range packageInfo.AllTypeInfos {
			PrintConstant(0, []dao.ITypeInfo{typeInfo})
		}
		// fmt.Println("========= Func =============")
		// for _, funcInfo := range packageInfo.AllFuncInfo {
		// 	funcInfo.Print()
		// }
	}
}

func PrintConstant(depth int, infos []dao.ITypeInfo) {
	for _, typrInfo := range infos {

		switch info := typrInfo.(type) {
		case *dao.TypeInfo:
			PrintType(depth+1, info.GetName(), info.ContentTypeInfo)
		case *dao.VarInfo:
			PrintType(depth+1, info.GetName(), info.ContentTypeInfo)
		case *dao.ConstInfo:
			PrintType(depth+1, info.GetName(), info.ContentTypeInfo)
		default:
			PrintType(depth+1, info.GetName(), info)
		}
	}
}

func PrintType(depth int, name string, iInfo dao.ITypeInfo) {
	if name == "TypeInfoNumeric" {
		fmt.Println("")
	}

	space := ""
	for i := 0; i < depth; i++ {
		space += "  "
	}
	switch info := iInfo.(type) {
	case *dao.TypeInfoArray:
		var childs []dao.ITypeInfo
		childs = append(childs, info.ContentTypeInfo)
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, childs)

	case *dao.TypeInfoChannel:
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, []dao.ITypeInfo{info.ContentTypeInfo})

	case *dao.TypeInfoFunction:
		var childs []dao.ITypeInfo
		for _, child := range info.ParamsInPoint {
			childs = append(childs, child.ContentTypeInfo)
		}
		for _, child := range info.ParamsOutPoint {
			childs = append(childs, child.ContentTypeInfo)
		}

		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, childs)

	case *dao.TypeInfoInterface:
		var childs []dao.ITypeInfo
		childs = append(childs, info.MatchInfos...)
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, childs)

	case *dao.TypeInfoMap:
		var childs []dao.ITypeInfo
		childs = append(childs, info.KeyType)
		childs = append(childs, info.ValueType)
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, childs)

	case *dao.TypeInfoPointer:
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, []dao.ITypeInfo{info.ContentTypeInfo})

	case *dao.TypeInfoSlice:
		var childs []dao.ITypeInfo
		childs = append(childs, info.ContentTypeInfo)
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, childs)

	case *dao.TypeInfoQualifiedIdent:
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, []dao.ITypeInfo{info.ContentTypeInfo})

	case *dao.TypeInfoStruct:
		var childs []dao.ITypeInfo

		for _, child := range info.VarInfos {
			if child.Tag != "`json:\"-\"`" {
				childs = append(childs, child)
			}
		}
		for _, child := range info.ImplicitlyVarInfos {
			if child.Tag != "`json:\"-\"`" {
				childs = append(childs, child)
			}
		}
		fmt.Println(space, name, info.GetTypeName())
		PrintConstant(depth+1, childs)

	case *dao.TypeInfoBool:
		fmt.Println(space, name, info.GetTypeName())
	case *dao.TypeInfoNumeric:
		fmt.Println(space, name, info.GetTypeName())
	case *dao.TypeInfoString:
		fmt.Println(space, name, info.GetTypeName())

	}
}
