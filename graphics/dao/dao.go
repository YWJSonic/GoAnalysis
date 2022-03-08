package dao

import (
	"codeanalysis/graphics/constant"
	"fmt"
)

// 抽象物件
type Abstract struct{}

// 註解物件
type Annotaion struct{}

type PackageSpace struct {
	Name      string
	VarList   []*UserClass
	TypeList  []*UserClass
	Interface []Interface
	Color     string
}

func (self *PackageSpace) ToString() string {
	classStr := ""

	for _, v := range self.TypeList {
		classStr += v.ToString()
	}

	interfaceStr := ""
	for _, v := range self.Interface {
		interfaceStr += v.ToString()
	}

	varListStr := ""
	for _, v := range self.VarList {
		varListStr += v.ToString()
	}

	return fmt.Sprintf("package %s %s{\n%s\n%s\n%s}\n", self.Name, self.Color, interfaceStr, classStr, varListStr)
}

// 類型物件 # golang as struct
type Class struct {
	Name  string
	Field []string
}

func (self *Class) ToString() string {

	fieldStr := ""
	for _, field := range self.Field {
		fieldStr += fmt.Sprintf("\t%v\n", field)
	}

	return fmt.Sprintf("\tclass %s {\n%s\t}\n", self.Name, fieldStr)
}

// 實例物件
type Entity struct {
	Name  string
	Field []string
}

func (self *Entity) ToString() string {

	fieldStr := ""
	for _, field := range self.Field {
		fieldStr += fmt.Sprintf("\t%v\n", field)
	}

	return fmt.Sprintf("\tentity %s {\n%s\t}\n", self.Name, fieldStr)
}

// 使用者字定義 spot 類型
type UserClass struct {
	Name      string
	Field     []string
	SpotWord  rune
	SpotColor string
}

func (self *UserClass) ToString() string {

	fieldStr := ""
	for _, field := range self.Field {
		fieldStr += fmt.Sprintf("\t%v\n", field)
	}

	return fmt.Sprintf("\tclass %s <<(%s,%s)>> {\n%s\t}\n", self.Name, string(self.SpotWord), self.SpotColor, fieldStr)
}

func NewTypeClass() *UserClass {
	return &UserClass{
		SpotWord:  'T',
		SpotColor: constant.SportCss['T'],
	}
}

func NewVarClass() *UserClass {
	return &UserClass{
		SpotWord:  'V',
		SpotColor: constant.SportCss['V'],
	}
}

func NewConstClass() *UserClass {
	return &UserClass{
		SpotWord:  'C',
		SpotColor: constant.SportCss['C'],
	}
}

// 接口物件
type Interface struct {
	Name  string
	Field []string
}

func (self *Interface) ToString() string {

	fieldStr := ""
	for _, field := range self.Field {
		fieldStr += fmt.Sprintf("\t\t%v\n", field)
	}

	return fmt.Sprintf("\tinterface %s {\n%s\t}\n", self.Name, fieldStr)
}

type Together struct {
	Obj []string
}

func (self *Together) ToString() string {
	str := "together {\n"
	for _, className := range self.Obj {
		str += fmt.Sprintf("\tclass %s\n", className)
	}

	str += "}\n"
	return str
}
