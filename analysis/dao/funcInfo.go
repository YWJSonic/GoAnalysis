package dao

type FuncInfo struct {
	PointBase
	Method         *MethodExpression // 篩選器
	ParamsInPoint  []*FuncParams     // 輸入參數
	ParamsOutPoint []*FuncParams     // 輸出參數
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
func NewFuncParams() *FuncParams {
	params := &FuncParams{}
	return params
}

type MethodExpression struct {
	PointBase
	StructPoint *StructInfo
}

// 篩選器
func NewMethodExpression(name string, structPoint *StructInfo) *MethodExpression {
	params := &MethodExpression{}
	params.Name = name
	params.StructPoint = structPoint
	return params
}
