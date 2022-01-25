package dao

// Func關聯
func NewFuncInfo(name string) *FuncInfo {
	info := &FuncInfo{}
	info.Name = name
	return info
}

type FuncInfo struct {
	PointBase
	Method         *MethodExpression
	ParamsInPoint  []*FuncParams
	ParamsOutPoint []*FuncParams
}

// Func 傳輸參數
func NewFuncParams(name string, structPoint *StructInfo) *FuncParams {
	params := &FuncParams{}
	params.Name = name
	params.StructPoint = structPoint
	return params
}

type FuncParams struct {
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

type MethodExpression struct {
	PointBase
	StructPoint *StructInfo
}
