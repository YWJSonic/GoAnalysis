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

func (self *FuncInfo) GetReceiver() string {
	var receiver string
	if self.Receiver.ContentTypeInfo == nil {
		return ""
	} else {
		rev := self.Receiver.GetName()
		ret := ""
		if info, ok := self.Receiver.ContentTypeInfo.(*TypeInfoPointer); ok {
			ret = "*" + info.ContentTypeInfo.GetTypeName()
		}
		receiver = fmt.Sprintf("(%s %s)", rev, ret)
	}

	return receiver
}

func (self *FuncInfo) GetNameKey() string {
	if self.Receiver.ContentTypeInfo != nil {
		return fmt.Sprintf("%s %s", self.GetReceiver(), self.TypeBase.GetName())
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
