package graphics

import (
	ansdao "codeanalysis/analysis/dao"
	"codeanalysis/graphics/constant"
	"codeanalysis/graphics/dao"
	"fmt"
)

type PlaneUml struct {
	context []string
	data    *ansdao.ProjectInfo
}

func (self *PlaneUml) Start() {
	self.context = append(self.context, "@startuml")
	inc := 0xffffff / (len(self.data.LocalPackageMap) * 2)
	inccolor := 0x000000
	// package 轉換
	for packageName, packageInfo := range self.data.LocalPackageMap {
		inccolor += inc
		color := fmt.Sprintf("#%06x", inccolor)

		inccolor += inc
		lineColor := fmt.Sprintf("#%06x", inccolor)
		lineCss := fmt.Sprintf("%s;line.%s", lineColor, constant.LineStyle_Bold)

		// package 寫入開始
		self.context = append(self.context, fmt.Sprintf("package %s %s{\n", packageName, color))

		// type 轉換
		for typeName, typeInfo := range packageInfo.AllTypeInfos {
			class := dao.Class{}
			class.Name = packageInfo.GetName() + "_" + typeName

			switch info := typeInfo.ContentTypeInfo.(type) {
			case *ansdao.TypeInfoStruct:
				for varName, varInfo := range info.VarInfos {
					varStr := fmt.Sprintf("\t%v %v", varInfo.ContentTypeInfo.GetTypeName(), varName)
					class.Field = append(class.Field, varStr)

					switch info := varInfo.ContentTypeInfo.(type) {
					case *ansdao.TypeInfoQualifiedIdent:

						// line 寫入
						lineStr := fmt.Sprintf("%s %s %s %s", class.Name, constant.Line_Normal, info.ContentTypeInfo.GetTypeName(), lineCss)
						class.Line = append(class.Line, lineStr)
					}
				}

				// type 寫入
				self.context = append(self.context, class.ToString())
			case *ansdao.TypeInfoInterface:

			}

		}

		// package 寫入結束
		self.context = append(self.context, "}\n")
	}

	self.context = append(self.context, "@enduml")
}

func (self *PlaneUml) ToString() string {
	contextStr := ""
	for _, contextLine := range self.context {
		contextStr += fmt.Sprintf("%s\n", contextLine)
	}
	return contextStr
}
