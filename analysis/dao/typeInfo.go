package dao

// Type declarations 類型聲明結構

type ITypeInfo interface {
	GetName() string
	SetName(name string)
	GetTypeName() string
	SetTypeName(typeName string)
}

// 別名宣告

type TypeAliasDecl struct {
	PointBase
	ContentTypeInfo ITypeInfo
}

func NewTypeAliasDecl() *TypeAliasDecl {
	return &TypeAliasDecl{}
}

// 類型定義

type TypeDef struct {
	PointBase
	ContentTypeInfo ITypeInfo
}

func NewTypeDef() *TypeDef {
	return &TypeDef{}
}

// struct 類型
type TypeInfoStruct struct {
	PointBase
	VarInfos  map[string]*VarInfo
	FuncPoint map[string]*FuncInfo
}

func NewTypeInfoStruct() *TypeInfoStruct {
	return &TypeInfoStruct{
		VarInfos:  make(map[string]*VarInfo),
		FuncPoint: make(map[string]*FuncInfo),
	}
}

// 指標類型
type TypeInfoPointer struct {
	PointBase
	ContentTypeInfo ITypeInfo
}

func NewTypeInfoPointer() *TypeInfoPointer {
	return &TypeInfoPointer{}
}

// func 類型
type TypeInfoFunction struct {
	PointBase
	ParamsInPoint  []*FuncParams // 輸入參數
	ParamsOutPoint []*FuncParams // 輸出參數
	Common         string        // 註解
}

func NewTypeInfoFunction() *TypeInfoFunction {
	return &TypeInfoFunction{}
}

// interface 類型
type TypeInfoInterface struct {
	PointBase
	MatchInfos []*FuncInfo
}

func NewTypeInfoInterface() *TypeInfoInterface {
	return &TypeInfoInterface{}
}

// [] 類型
type TypeInfoSlice struct {
	PointBase
	ContentTypeInfo ITypeInfo // 資料型別
}

func NewTypeInfoSlice() *TypeInfoSlice {
	return &TypeInfoSlice{}
}

// [n] 類型
type TypeInfoArray struct {
	PointBase
	Size            string    // 陣列大小
	ContentTypeInfo ITypeInfo // 資料型別
}

func NewTypeInfoArray() *TypeInfoArray {
	return &TypeInfoArray{}
}

// map 類型
type TypeInfoMap struct {
	PointBase
	KeyType   ITypeInfo
	ValueType ITypeInfo
}

func NewTypeInfoMap() *TypeInfoMap {
	return &TypeInfoMap{}
}

// chan 類型
type TypeInfoChannel struct {
	PointBase
	FlowType        int       // channel 類型: 雙向, 單出, 單進
	ContentTypeInfo ITypeInfo // channel 傳輸型別
}

func NewTypeInfoChannel() *TypeInfoChannel {
	return &TypeInfoChannel{}
}

// 數值基礎類型
type TypeInfoNumeric struct {
	PointBase
}

// string 類型
type TypeInfoString struct {
	PointBase
}

// bool 類型
type TypeInfoBool struct {
	PointBase
}

// package.A 類型
type TypeInfoQualifiedIdent struct {
	PointBase
	ImportLink      *ImportInfo // 指定的 import package
	ContentTypeInfo ITypeInfo   // 指定該包的 type
}

func NewTypeInfoQualifiedIdent() *TypeInfoQualifiedIdent {
	return &TypeInfoQualifiedIdent{}
}

// 基礎類型表
var BaseTypeInfo map[string]ITypeInfo = map[string]ITypeInfo{
	"bool":       &TypeInfoBool{PointBase: PointBase{typeName: "bool"}},
	"byte":       &TypeInfoNumeric{PointBase: PointBase{typeName: "byte"}},
	"complex64":  &TypeInfoNumeric{PointBase: PointBase{typeName: "complex64"}},
	"complex128": &TypeInfoNumeric{PointBase: PointBase{typeName: "complex128"}},
	"float32":    &TypeInfoNumeric{PointBase: PointBase{typeName: "float32"}},
	"float64":    &TypeInfoNumeric{PointBase: PointBase{typeName: "float64"}},
	"int":        &TypeInfoNumeric{PointBase: PointBase{typeName: "int"}},
	"int8":       &TypeInfoNumeric{PointBase: PointBase{typeName: "int8"}},
	"int16":      &TypeInfoNumeric{PointBase: PointBase{typeName: "int16"}},
	"int32":      &TypeInfoNumeric{PointBase: PointBase{typeName: "int32"}},
	"int64":      &TypeInfoNumeric{PointBase: PointBase{typeName: "int64"}},
	"uint":       &TypeInfoNumeric{PointBase: PointBase{typeName: "uint"}},
	"uint8":      &TypeInfoNumeric{PointBase: PointBase{typeName: "uint8"}},
	"uint16":     &TypeInfoNumeric{PointBase: PointBase{typeName: "uint16"}},
	"uint32":     &TypeInfoNumeric{PointBase: PointBase{typeName: "uint32"}},
	"uint64":     &TypeInfoNumeric{PointBase: PointBase{typeName: "uint64"}},
	"uintptr":    &TypeInfoNumeric{PointBase: PointBase{typeName: "uintptr"}},
	"error":      &TypeInfoStruct{PointBase: PointBase{typeName: "error"}},
	"string":     &TypeInfoString{PointBase: PointBase{typeName: "string"}},
	"rune":       &TypeInfoString{PointBase: PointBase{typeName: "rune"}},
}
