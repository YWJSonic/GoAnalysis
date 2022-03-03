package dao

import "fmt"

// 抽象物件
type Abstract struct{}

// 註解物件
type Annotaion struct{}

type PackageSpace struct {
	Name      string
	ClassList []Class
}

// 類型物件 # golang as struct
type Class struct {
	Name   string
	Field  []string
	Method []string
	Line   []string
}

func (self *Class) ToString() string {

	fieldStr := ""
	for _, field := range self.Field {
		fieldStr += fmt.Sprintf("%v\n", field)
	}

	methodStr := ""
	for _, method := range self.Method {
		methodStr += fmt.Sprintf("%v\n", method)
	}

	lineStr := ""
	for _, line := range self.Line {
		lineStr += fmt.Sprintf("%v\n", line)
	}

	return fmt.Sprintf("class %s {\n%s%s}\n%s", self.Name, fieldStr, methodStr, lineStr)

}

// 實例物件
type Entity struct{}

// 列舉物件
type Enum struct{}

// 接口物件
type Interface struct{}
