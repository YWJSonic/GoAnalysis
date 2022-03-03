package dao

import "fmt"

type FuncInfo struct {
	TypeBase
	Receiver       FuncParams   // 篩選器
	ParamsInPoint  []FuncParams // 輸入參數
	ParamsOutPoint []FuncParams // 輸出參數
	Body           string       `json:"-"` // 方法內文 *尚未解析
}

func (self *FuncInfo) GetTypeBase() TypeBase {
	return self.TypeBase
}

func (self *FuncInfo) Print() string {
	str := "func("
	if self.Receiver.ContentTypeInfo != nil {
		str += self.Receiver.GetName() + " " + self.Receiver.ContentTypeInfo.GetTypeName()
	}
	str += ") " + self.GetName() + "("

	for idx, input := range self.ParamsInPoint {
		str += input.GetName() + " " + input.ContentTypeInfo.GetTypeName()
		if len(self.ParamsInPoint) > idx+1 {
			str += ", "
		}
	}
	str += ") ("

	for idx, output := range self.ParamsOutPoint {
		str += output.GetName() + " " + output.ContentTypeInfo.GetName()
		if len(self.ParamsOutPoint) > idx+1 {
			str += ", "
		}
	}
	str += ")"

	return str
}

func (self *FuncInfo) GetName() string {
	if self.Receiver.ContentTypeInfo != nil {
		rev := self.Receiver.GetName()
		ret := self.Receiver.ContentTypeInfo.GetTypeName()
		name := self.TypeBase.GetName()
		return fmt.Sprintf("(%s %s)%s", rev, ret, name)
	} else {
		return self.TypeBase.GetName()
	}
}

// Func關聯
func NewFuncInfo() *FuncInfo {
	info := &FuncInfo{
		TypeBase: NewPointBase(),
	}
	return info
}

type FuncParams struct {
	TypeBase
	ContentTypeInfo ITypeInfo
}

// Func 傳輸參數
func NewFuncParams() FuncParams {
	params := FuncParams{
		TypeBase: NewPointBase(),
	}
	return params
}
