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
var disableNamsSpace = map[string]struct{}{
	"GameService.games": {},
	"gitlab.geax.io":    {},
}

type PlaneUml struct {
	context  []string
	packages []*dao.NameSpace
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
	inc := 0xffffff / (len(allPackageMap) * 3)
	inccolor := 0xffffff

	// package 資料轉換
	for packageName, packageInfo := range allPackageMap {
		// 設定 line 圖形資料
		inccolor -= inc
		color := fmt.Sprintf("#%06x", inccolor)

		inccolor -= inc
		lineColor := fmt.Sprintf("#%06x", inccolor)
		lineCss := fmt.Sprintf("%s;line.%s;", lineColor, constant.LineStyle_Dashed)
		lineType := constant.Line_Normal + ">"

		// package 格式資料
		nameSpace := &dao.NameSpace{
			Color: color,
		}
		nameSpace.SetName(ReplaceName(packageName))

		if OutputLevel > 1 {
			for typeName, typeInfo := range packageInfo.AllTypeInfos {
				// 初始化 type 圖形資料
				// var lines = map[string]struct{}{}
				// type 關聯 line 圖形

				if typeInfo.DefType == anscon.DefType_Def {
					if typeInfo.ContentTypeInfo == nil {
						class := dao.NewTypeClass()
						class.Name = typeName
						nameSpace.TypeList = append(nameSpace.TypeList, class)

					} else {
						switch info := typeInfo.ContentTypeInfo.(type) {
						case *ansdao.TypeInfoStruct:
							// struct 內物件圖形
							class := dao.NewTypeClass()
							class.Name = typeName
							for varName, varInfo := range info.VarInfos {
								varStr := fmt.Sprintf("\t%v %v", varInfo.ContentTypeInfo.GetTypeName(), varName)
								class.Field = append(class.Field, varStr)
								// struct 內物件關聯 line 圖形
								for _, lineStr := range lineStrs(nameSpace.GetName(), typeName, varInfo.ContentTypeInfo) {
									if lineStr != "" {
										self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
									}
								}
							}
							nameSpace.TypeList = append(nameSpace.TypeList, class)

						case *ansdao.TypeInfoInterface:

							interfaceInfo := dao.Interface{}
							interfaceInfo.Name = typeName
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
								count := len(varInfo.ParamsOutPoint)
								if count > 1 {
									varStr += "("
								}
								for idx := 0; idx < count; idx++ {
									if idx != 0 {
										varStr += ","
									}

									varStr += varInfo.ParamsOutPoint[idx].ContentTypeInfo.GetTypeName()
								}
								if count > 1 {
									varStr += ")"
								}

								interfaceInfo.Field = append(interfaceInfo.Field, varStr)
							}

							nameSpace.Interface = append(nameSpace.Interface, interfaceInfo)
						default:
							class := dao.NewTypeClass()
							class.Name = typeName
							for _, lineStr := range lineStrs(nameSpace.GetName(), typeName, typeInfo.ContentTypeInfo) {
								if lineStr != "" {
									self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
								}
							}
							nameSpace.TypeList = append(nameSpace.TypeList, class)
							// fmt.Println(typeName, GetContentAllTypeName(typeInfo.ContentTypeInfo))

						}
					}

				} else if typeInfo.DefType == anscon.DefType_Decl {
					// 替換名稱 or 定義類型
					class := dao.NewTypeClass()
					class.Name = typeName
					if typeInfo.ContentTypeInfo != nil {
						for _, lineStr := range lineStrs(nameSpace.GetName(), typeName, typeInfo.ContentTypeInfo) {
							if lineStr != "" {
								self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
							}
						}
					}
					nameSpace.TypeList = append(nameSpace.TypeList, class)

				} else {
					panic("")
				}
			}

			for varname, varInfo := range packageInfo.AllVarInfos {
				class := dao.NewVarClass()
				// class.Name = packageInfo.GetName() + "_" + varname
				class.Name = varname

				if varInfo.ContentTypeInfo != nil {
					for _, lineStr := range lineStrs(nameSpace.GetName(), varname, varInfo.ContentTypeInfo) {
						if lineStr != "" {
							self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
						}
					}
				} else {
					fmt.Println("var ContentTypeInfo empty ")
				}

				nameSpace.VarList = append(nameSpace.VarList, class)

			}

			for constname, constInfo := range packageInfo.AllConstInfos {
				class := dao.NewConstClass()
				// class.Name = packageInfo.GetName() + "_" + varname
				class.Name = constname

				if constInfo.ContentTypeInfo != nil {
					for _, lineStr := range lineStrs(nameSpace.GetName(), constname, constInfo.ContentTypeInfo) {
						if lineStr != "" {
							self.line[fmt.Sprintf(lineStr, lineType, lineCss)] = struct{}{}
						}
					}
				} else {
					fmt.Println("const ContentTypeInfo empty ")
				}

				nameSpace.ConstList = append(nameSpace.ConstList, class)

			}

		} else {

			for _, ImportLink := range packageInfo.AllImportLink {
				lineStr := fmt.Sprintf("\"%s\" %s \"%s\" %s", packageInfo.GoPath, lineType, ImportLink.Package.GoPath, lineCss)
				self.line[lineStr] = struct{}{}
			}
		}

		self.packages = append(self.packages, nameSpace)
	}
	self.goModCheck()
}

func (self *PlaneUml) ToString() string {

	self.context = append(self.context, "@startuml")

	sort.Slice(self.packages, func(i, j int) bool {
		return self.packages[i].GetName() < self.packages[j].GetName()
	})

	for _, packageSpace := range self.packages {
		if !IsDisableStr(packageSpace.GetName(), disableNamsSpace) {
			self.context = append(self.context, packageSpace.ToString())
		}
	}

	lineStrs := []string{}
	for line := range self.line {
		if !IsDisableStr(line, disableNamsSpace) {
			lineStrs = append(lineStrs, line)
		}
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

func (self *PlaneUml) goModCheck() {
	var isExist bool
	for name := range self.data.ModuleInfo.VendorMap {
		spacename := ReplaceName(name)

		isExist = false
		for _, spaceInfo := range self.packages {
			if spaceInfo.GetName() == spacename {
				isExist = true
				break
			}
		}

		if !isExist {
			nameSpaceInfo := &dao.NameSpace{}
			nameSpaceInfo.SetName(spacename)
			self.packages = append(self.packages, nameSpaceInfo)
		}
	}
}
