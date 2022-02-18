package dao

import "fmt"

type FuncInfo struct {
	PointBase
	Receiver       FuncParams   // 篩選器
	ParamsInPoint  []FuncParams // 輸入參數
	ParamsOutPoint []FuncParams // 輸出參數
	Body           string       // 方法內文 *尚未解析
}

func (self *FuncInfo) Print() string {
	str := "func("
	if self.Receiver.ContentTypeInfo != nil {
		str += self.Receiver.name + " " + self.Receiver.ContentTypeInfo.GetTypeName()
	}
	str += ") " + self.name + "("

	for idx, input := range self.ParamsInPoint {
		str += input.name + " " + input.ContentTypeInfo.GetTypeName()
		if len(self.ParamsInPoint) > idx+1 {
			str += ", "
		}
	}
	str += ") ("

	for idx, output := range self.ParamsOutPoint {
		str += output.name + " " + output.ContentTypeInfo.GetName()
		if len(self.ParamsOutPoint) > idx+1 {
			str += ", "
		}
	}
	str += ")"

	return str
}

func (self *FuncInfo) GetName() string {
	if self.Receiver.ContentTypeInfo != nil {
		return fmt.Sprintf("(%s %s)%s", self.Receiver.name, self.Receiver.ContentTypeInfo.GetTypeName(), self.name)
	} else {
		return self.name
	}
}

// Func關聯
func NewFuncInfo() *FuncInfo {
	info := &FuncInfo{}
	return info
}

type FuncParams struct {
	PointBase
	ContentTypeInfo ITypeInfo
}

// Func 傳輸參數
func NewFuncParams() FuncParams {
	params := FuncParams{}
	return params
}
