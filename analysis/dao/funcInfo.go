package dao

import "fmt"

type FuncInfo struct {
	PointBase
	Receiver       FuncParams   // 篩選器
	ParamsInPoint  []FuncParams // 輸入參數
	ParamsOutPoint []FuncParams // 輸出參數
	Body           string       // 方法內文 *尚未解析
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
