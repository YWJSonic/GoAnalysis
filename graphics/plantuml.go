package graphics

import (
	anscon "codeanalysis/analysis/constant"
	ansdao "codeanalysis/analysis/dao"
	"codeanalysis/graphics/constant"
	"codeanalysis/graphics/dao"
	"fmt"
	"sort"
)

var disableLineKeyword = map[string]struct{}{
	"bool":       {},
	"byte":       {},
	"complex64":  {},
	"complex128": {},
	"float32":    {},
	"float64":    {},
	"int":        {},
	"int8":       {},
	"int16":      {},
	"int32":      {},
	"int64":      {},
	"uint":       {},
	"uint8":      {},
	"uint16":     {},
	"uint32":     {},
	"uint64":     {},
	"uintptr":    {},
	"string":     {},
	"rune":       {},
	"error":      {},
	"struct":     {},
	"iota":       {},
	"interface":  {},
}

type PlaneUml struct {
	context  []string
	packages []*dao.PackageSpace
	line     map[string]struct{}
	data     *ansdao.ProjectInfo
}

// 開始轉換資料
//
// @params uint 關聯層級,數值越大越詳細 1：package, import 層, 2:  type, func, var 層, 3: struct, array, map, method
func (self *PlaneUml) Start(OutputLevel uint) {
	// 整合全部 package 資料
	allPackageMap := map[string]*ansdao.PackageInfo{}
	for k, v := range self.data.LocalPackageMap {
		allPackageMap[k] = v
	}

	for k, v := range self.data.VendorPackageMap {
		allPackageMap[k] = v
	}

	for k, v := range self.data.GolangPackageMap {
		allPackageMap[k] = v
	}

	// 初始化圖形資料
	inc := 0xffffff / (len(allPackageMap) * 2)
	inccolor := 0x000000

	// package 資料轉換
	for packageName, packageInfo := range allPackageMap {

		if packageName == "demeter/match" {
			fmt.Println("")
		}
		// 設定 line 圖形資料
		inccolor += inc
		color := fmt.Sprintf("#%06x", inccolor)

		inccolor += inc
		lineColor := fmt.Sprintf("#%06x", inccolor)
		lineCss := fmt.Sprintf("%s;line.%s;", lineColor, constant.LineStyle_Dashed)
		lineType := constant.Line_Normal + ">"

		// package 格式資料
		packageSpace := &dao.PackageSpace{
			Name:  packageName,
			Color: color,
		}

		if OutputLevel > 1 {
			for typeName, typeInfo := range packageInfo.AllTypeInfos {
				// 初始化 type 圖形資料
				// var lines = map[string]struct{}{}
				// type 關聯 line 圖形

				if typeInfo.DefType == anscon.DefType_Def {
					if typeInfo.ContentTypeInfo == nil {
						class := dao.NewTypeClass()
						class.Name = packageInfo.GetName() + "_" + typeName
						packageSpace.TypeList = append(packageSpace.TypeList, class)

					} else {
						switch info := typeInfo.ContentTypeInfo.(type) {
						case *ansdao.TypeInfoStruct:
							// struct 內物件圖形
							class := dao.NewTypeClass()
							class.Name = packageInfo.GetName() + "_" + typeName
							for varName, varInfo := range info.VarInfos {
								varStr := fmt.Sprintf("\t%v %v", varInfo.ContentTypeInfo.GetTypeName(), varName)
								class.Field = append(class.Field, varStr)
								// struct 內物件關聯 line 圖形
								for _, lineStr := range self.lineStrs(packageInfo.GetName(), typeName, varInfo.ContentTypeInfo) {
									if lineStr != "" {
										self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
									}
								}
							}
							packageSpace.TypeList = append(packageSpace.TypeList, class)

						case *ansdao.TypeInfoInterface:

							interfaceInfo := dao.Interface{}
							interfaceInfo.Name = packageInfo.GetName() + "_" + typeName
							for _, ivarInfo := range info.MatchInfos {
								// 接口類型
								varInfo := ivarInfo.(*ansdao.FuncInfo)
								varStr := fmt.Sprintf("%s(", varInfo.GetName())
								for idx, count := 0, len(varInfo.ParamsInPoint); idx < count; idx++ {
									if idx != 0 {
										varStr += ","
									}

									varStr += varInfo.ParamsInPoint[idx].ContentTypeInfo.GetTypeName()
								}
								varStr += ") "
								for idx, count := 0, len(varInfo.ParamsOutPoint); idx < count; idx++ {
									if idx != 0 {
										varStr += ","
									} else if count > 1 {
										varStr += "("
									}

									varStr += varInfo.ParamsOutPoint[idx].ContentTypeInfo.GetTypeName()
								}

								interfaceInfo.Field = append(interfaceInfo.Field, varStr)
							}

							packageSpace.Interface = append(packageSpace.Interface, interfaceInfo)
						default:
							// panic("")

						}
					}

				} else if typeInfo.DefType == anscon.DefType_Decl {
					// 替換名稱 or 定義類型
					class := dao.NewTypeClass()
					class.Name = packageInfo.GetName() + "_" + typeName
					if typeInfo.ContentTypeInfo != nil {
						for _, lineStr := range self.lineStrs(packageInfo.GetName(), typeName, typeInfo.ContentTypeInfo) {
							if lineStr != "" {
								self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
							}
						}
					}
					packageSpace.TypeList = append(packageSpace.TypeList, class)

				} else {
					panic("")
				}
			}

			for varname, varInfo := range packageInfo.AllVarInfos {
				class := dao.NewVarClass()
				class.Name = packageInfo.GetName() + "_" + varname

				if varInfo.ContentTypeInfo != nil {
					for _, lineStr := range self.lineStrs(packageInfo.GetName(), varname, varInfo.ContentTypeInfo) {
						if lineStr != "" {
							self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
						}
					}
				} else {
					fmt.Println("var ContentTypeInfo empty ")
				}

				packageSpace.VarList = append(packageSpace.VarList, class)

			}

		} else {

			for _, ImportLink := range packageInfo.AllImportLink {
				lineStr := fmt.Sprintf("\"%s\" %s \"%s\" %s", packageInfo.GoPath, lineType, ImportLink.Package.GoPath, lineCss)
				self.line[lineStr] = struct{}{}
			}
		}

		self.packages = append(self.packages, packageSpace)
	}
}

func (self *PlaneUml) ToString() string {

	self.context = append(self.context, "@startuml")

	for _, packageSpace := range self.packages {
		self.context = append(self.context, packageSpace.ToString())
	}

	lineStrs := []string{}
	for line := range self.line {
		lineStrs = append(lineStrs, line)
	}
	sort.Strings(lineStrs)

	for _, line := range lineStrs {
		self.context = append(self.context, fmt.Sprintf("%v", line))
	}

	self.context = append(self.context, "@enduml")

	contextStr := ""
	for _, contextLine := range self.context {
		contextStr += fmt.Sprintf("%s\n", contextLine)
	}
	return contextStr
}

func (self *PlaneUml) lineStrs(packageName, startTypeName string, target ansdao.ITypeInfo) []string {
	lineStrList := []string{}
	if _, ok := disableLineKeyword[target.GetTypeName()]; ok {
		return lineStrList
	}

	switch info := target.(type) {
	case *ansdao.TypeInfoQualifiedIdent:
		targetPackageName := info.ImportLink.Package.GetName()
		if targetPackageName == "" {
			targetPackageName = info.ImportLink.Package.GoPath
		}

		targetStr := fmt.Sprintf("%s_%s", targetPackageName, info.ContentTypeInfo.GetTypeName())
		lineStr := packageName + "_" + startTypeName + " %s " + targetStr + " %s"
		lineStrList = append(lineStrList, lineStr)

	case *ansdao.TypeInfoArray:
		lineStrList = self.lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoChannel:
		lineStrList = self.lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoSlice:
		lineStrList = self.lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoPointer:
		lineStrList = self.lineStrs(packageName, startTypeName, info.ContentTypeInfo)

	case *ansdao.TypeInfoStruct:
		if len(info.VarInfos) == 0 && len(info.ImplicitlyVarInfos) == 0 {
			lineStr := packageName + "_" + startTypeName + " %s struct %s"
			lineStrList = append(lineStrList, lineStr)
		} else {
			lineStr := packageName + "_" + startTypeName + " %s " + info.GetName() + " %s"
			lineStrList = append(lineStrList, lineStr)
		}

	case *ansdao.TypeInfo:
		lineStr := packageName + "_" + startTypeName + " %s " + packageName + "_" + info.GetName() + " %s"
		lineStrList = append(lineStrList, lineStr)

	case *ansdao.TypeInfoMap:
		lineStrList = append(lineStrList, self.lineStrs(packageName, startTypeName, info.KeyType)...)
		lineStrList = append(lineStrList, self.lineStrs(packageName, startTypeName, info.ValueType)...)

	case *ansdao.TypeInfoFunction:

		for _, paramsInfo := range info.ParamsInPoint {
			lineStrList = append(lineStrList, self.lineStrs(packageName, startTypeName, paramsInfo.ContentTypeInfo)...)
		}

		for _, paramsInfo := range info.ParamsOutPoint {
			lineStrList = append(lineStrList, self.lineStrs(packageName, startTypeName, paramsInfo.ContentTypeInfo)...)
		}

	case *ansdao.TypeInfoInterface:
		if len(info.MatchInfos) == 0 {
			lineStr := fmt.Sprintf("%s_%s %s interface %s", packageName, startTypeName, "%s", "%s")
			lineStrList = append(lineStrList, lineStr)

		} else {
			lineStr := packageName + "_" + startTypeName + " %s " + packageName + "_" + info.GetName() + " %s"
			lineStrList = append(lineStrList, lineStr)

		}

	case *ansdao.TypeInfoNumeric, *ansdao.TypeInfoString, *ansdao.TypeInfoBool:
		lineStr := packageName + "_" + startTypeName + " %s " + info.GetTypeName() + " %s"
		lineStrList = append(lineStrList, lineStr)
	}

	return lineStrList
}
