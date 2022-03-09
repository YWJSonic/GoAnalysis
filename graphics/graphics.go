package graphics

import (
	"codeanalysis/analysis/dao"
	ansdao "codeanalysis/analysis/dao"
	"fmt"
	"io/ioutil"
	"strings"
)

func OutputPlantUmlGraphics(data *dao.ProjectInfo) {
	uml := PlaneUml{
		data: data,
		line: make(map[string]struct{}),
	}
	uml.Start(2)

	path := fmt.Sprintf("./%s.plantuml", data.ModuleInfo.ModuleName)
	ioutil.WriteFile(path, []byte(uml.ToString()), 0666)
}

func lineStrs(packageName, startTypeName string, target ansdao.ITypeInfo) []string {
	lineStrList := []string{}
	if _, ok := disableLineKeyword[target.GetTypeName()]; ok {
		return lineStrList
	}

	switch info := target.(type) {
	case *ansdao.TypeInfoQualifiedIdent:
		targetNameSapce := ReplaceName(info.ImportLink.Package.GoPath)
		lineStr := fmt.Sprintf("%s.%s %%s %s.%s %%s", packageName, startTypeName, targetNameSapce, info.ContentTypeInfo.GetTypeName())
		lineStrList = append(lineStrList, lineStr)

	case *ansdao.TypeInfoArray:
		lineStrList = lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoChannel:
		lineStrList = lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoSlice:
		lineStrList = lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoPointer:
		lineStrList = lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoStruct:
		if len(info.VarInfos) == 0 && len(info.ImplicitlyVarInfos) == 0 {
			lineStr := packageName + "." + startTypeName + " %s struct %s"
			lineStrList = append(lineStrList, lineStr)
		} else {
			lineStr := packageName + "." + startTypeName + " %s " + info.GetName() + " %s"
			lineStrList = append(lineStrList, lineStr)
		}

	case *ansdao.TypeInfo:
		lineStr := packageName + "." + startTypeName + " %s " + packageName + "." + info.GetName() + " %s"
		lineStrList = append(lineStrList, lineStr)

	case *ansdao.TypeInfoMap:
		lineStrList = append(lineStrList, lineStrs(packageName, startTypeName, info.KeyType)...)
		lineStrList = append(lineStrList, lineStrs(packageName, startTypeName, info.ValueType)...)

	case *ansdao.TypeInfoFunction:

		for _, paramsInfo := range info.ParamsInPoint {
			lineStrList = append(lineStrList, lineStrs(packageName, startTypeName, paramsInfo.ContentTypeInfo)...)
		}

		for _, paramsInfo := range info.ParamsOutPoint {
			lineStrList = append(lineStrList, lineStrs(packageName, startTypeName, paramsInfo.ContentTypeInfo)...)
		}

	case *ansdao.TypeInfoInterface:
		if len(info.MatchInfos) == 0 {
			lineStr := fmt.Sprintf("%s.%s %s interface %s", packageName, startTypeName, "%s", "%s")
			lineStrList = append(lineStrList, lineStr)

		} else {
			lineStr := packageName + "." + startTypeName + " %s " + packageName + "." + info.GetName() + " %s"
			lineStrList = append(lineStrList, lineStr)

		}

	case *ansdao.TypeInfoNumeric, *ansdao.TypeInfoString, *ansdao.TypeInfoBool:
		lineStr := packageName + "." + startTypeName + " %s " + info.GetTypeName() + " %s"
		lineStrList = append(lineStrList, lineStr)
	}

	return lineStrList
}

var nameReplacer = strings.NewReplacer(".", "_", "-", "_", "/", ".")

func ReplaceName(name string) string {
	return nameReplacer.Replace(name)
}

func GetContentAllTypeName(target ansdao.ITypeInfo) []string {
	lineStrList := []string{}
	switch info := target.(type) {
	case *ansdao.TypeInfoQualifiedIdent:
		targetNameSapce := ReplaceName(info.ImportLink.Package.GoPath)
		lineStrList = append(lineStrList, fmt.Sprintf("%s.%s", targetNameSapce, info.ContentTypeInfo.GetTypeName()))

	case *ansdao.TypeInfoArray:
		lineStrList = GetContentAllTypeName(info.ContentTypeInfo)

	case *ansdao.TypeInfoChannel:
		lineStrList = GetContentAllTypeName(info.ContentTypeInfo)

	case *ansdao.TypeInfoSlice:
		lineStrList = GetContentAllTypeName(info.ContentTypeInfo)

	case *ansdao.TypeInfoPointer:
		lineStrList = GetContentAllTypeName(info.ContentTypeInfo)

	case *ansdao.TypeInfoStruct:
		if len(info.VarInfos) == 0 && len(info.ImplicitlyVarInfos) == 0 {
			lineStrList = append(lineStrList, "struct")
		} else {
			lineStrList = append(lineStrList, info.GetName())
		}

	case *ansdao.TypeInfo:
		lineStrList = append(lineStrList, info.GetName())

	case *ansdao.TypeInfoMap:
		lineStrList = append(lineStrList, GetContentAllTypeName(info.KeyType)...)
		lineStrList = append(lineStrList, GetContentAllTypeName(info.ValueType)...)
	case *ansdao.TypeInfoFunction:
		for _, paramsInfo := range info.ParamsInPoint {
			lineStrList = append(lineStrList, GetContentAllTypeName(paramsInfo.ContentTypeInfo)...)
		}

		for _, paramsInfo := range info.ParamsOutPoint {
			lineStrList = append(lineStrList, GetContentAllTypeName(paramsInfo.ContentTypeInfo)...)
		}

	case *ansdao.TypeInfoInterface:
		if len(info.MatchInfos) == 0 {
			lineStrList = append(lineStrList, "interface")

		} else {
			lineStrList = append(lineStrList, info.GetName())

		}

	case *ansdao.TypeInfoNumeric, *ansdao.TypeInfoString, *ansdao.TypeInfoBool:
		lineStrList = append(lineStrList, info.GetTypeName())
	}
	return lineStrList
}
