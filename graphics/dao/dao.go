package dao

import (
	"fmt"
)

// 抽象物件
type Abstract struct{}

// 註解物件
type Annotaion struct{}

type PackageSpace struct {
	Name      string
	ClassList []Class
	Interface []Interface
	Color     string
}

func (self *PackageSpace) ToString() string {
	classStr := ""

	if len(self.ClassList) == 0 {
		return classStr
	}

	for _, v := range self.ClassList {
		classStr += v.ToString()
	}

	interfaceStr := ""
	for _, v := range self.Interface {
		interfaceStr += v.ToString()
	}

	return fmt.Sprintf("package %s %s{\n%s\n%s}\n", self.Name, self.Color, interfaceStr, classStr)
}

// 類型物件 # golang as struct
type Class struct {
	Name   string
	Field  []string
	Method []string
}

func (self *Class) ToString() string {

	fieldStr := ""
	for _, field := range self.Field {
		fieldStr += fmt.Sprintf("\t%v\n", field)
	}

	methodStr := ""
	for _, method := range self.Method {
		methodStr += fmt.Sprintf("\t{method} %v\n", method)
	}

	return fmt.Sprintf("\tclass %s {\n%s%s\t}\n", self.Name, fieldStr, methodStr)
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

// 列舉物件
type Enum struct{}

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
