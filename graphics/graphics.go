package graphics

import (
	"codeanalysis/analysis/dao"
	"fmt"
	"io/ioutil"
)

func OutputPlantUmlGraphics(data *dao.ProjectInfo) {
	uml := PlaneUml{data: data}
	uml.Start()

	// packangeNames := []string{}
	// for name := range data.LocalPackageMap {
	// 	packangeNames = append(packangeNames, name)
	// }
	// sort.Strings(packangeNames)
	// fmt.Println(packangeNames)

	path := fmt.Sprintf("./%s.plantuml", data.ModuleInfo.ModuleName)
	ioutil.WriteFile(path, []byte(uml.ToString()), 0666)
}
