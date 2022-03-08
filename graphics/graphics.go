package graphics

import (
	"codeanalysis/analysis/dao"
	"fmt"
	"io/ioutil"
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
