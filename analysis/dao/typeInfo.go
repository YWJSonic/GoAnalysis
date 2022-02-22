package dao

import (
	"codeanalysis/analysis/constant"
	"fmt"
)

// Type declarations 類型聲明結構

type ITypeInfo interface {
	GetName() string
	SetName(name string)
	GetTypeName() string
	SetTypeName(typeName string)
}

type TypeInfo struct {
	PointBase
	DefType         string
	ContentTypeInfo ITypeInfo
}

func (self *TypeInfo) GetTypeName() string {
	return self.name
}

// 別名宣告
func NewTypeAliasDecl() *TypeInfo {
	return &TypeInfo{
		DefType: "Decl",
	}
}

// 類型定義
func NewTypeDef() *TypeInfo {
	return &TypeInfo{
		DefType: "Def",
	}
}

// struct 類型
type TypeInfoStruct struct {
	PointBase

	ImplicitlyVarInfos []*VarInfo // 隱藏式宣告參數
	VarInfos           map[string]*VarInfo
	// FuncPoint          map[string]*FuncInfo
}

func (self *TypeInfoStruct) GetTypeName() string {
	return self.name
}
func NewTypeInfoStruct() *TypeInfoStruct {
	return &TypeInfoStruct{
		VarInfos: make(map[string]*VarInfo),
		// FuncPoint: make(map[string]*FuncInfo),
	}
}

// 指標類型
type TypeInfoPointer struct {
	PointBase
	ContentTypeInfo ITypeInfo
}

func (self *TypeInfoPointer) GetTypeName() string {
	return fmt.Sprintf("*%s", self.ContentTypeInfo.GetTypeName())
}

func NewTypeInfoPointer() *TypeInfoPointer {
	return &TypeInfoPointer{}
}

// func 類型
type TypeInfoFunction struct {
	PointBase
	ParamsInPoint  []FuncParams // 輸入參數
	ParamsOutPoint []FuncParams // 輸出參數
	Common         string       // 註解
}

func NewTypeInfoFunction() *TypeInfoFunction {
	return &TypeInfoFunction{}
}

// interface 類型
type TypeInfoInterface struct {
	PointBase
	MatchInfos []ITypeInfo // interface 定義的方法
}

func NewTypeInfoInterface() *TypeInfoInterface {
	return &TypeInfoInterface{}
}

// [] 類型
type TypeInfoSlice struct {
	PointBase
	ContentTypeInfo ITypeInfo // 資料型別
}

func (self *TypeInfoSlice) GetTypeName() string {
	return fmt.Sprintf("[]%s", self.ContentTypeInfo.GetTypeName())
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

func (self *TypeInfoArray) GetTypeName() string {
	return fmt.Sprintf("[%s]%s", self.Size, self.ContentTypeInfo.GetTypeName())
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
	RefBase
}

// string 類型
type TypeInfoString struct {
	PointBase
	RefBase
}

// bool 類型
type TypeInfoBool struct {
	PointBase
	RefBase
}

// package.A 類型
type TypeInfoQualifiedIdent struct {
	PointBase
	RefBase
	ImportLink      *ImportInfo // 指定的 import package
	ContentTypeInfo ITypeInfo   // 指定該包的 type
}

func (self *TypeInfoQualifiedIdent) GetTypeName() string {
	return self.name
}

func NewTypeInfoQualifiedIdent() *TypeInfoQualifiedIdent {
	return &TypeInfoQualifiedIdent{}
}

// 基礎類型表
var BaseTypeInfo map[string]ITypeInfo = map[string]ITypeInfo{
	"bool":       &TypeInfoBool{PointBase: PointBase{typeName: "bool"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"byte":       &TypeInfoNumeric{PointBase: PointBase{typeName: "byte"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"complex64":  &TypeInfoNumeric{PointBase: PointBase{typeName: "complex64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"complex128": &TypeInfoNumeric{PointBase: PointBase{typeName: "complex128"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"float32":    &TypeInfoNumeric{PointBase: PointBase{typeName: "float32"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"float64":    &TypeInfoNumeric{PointBase: PointBase{typeName: "float64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int":        &TypeInfoNumeric{PointBase: PointBase{typeName: "int"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int8":       &TypeInfoNumeric{PointBase: PointBase{typeName: "int8"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int16":      &TypeInfoNumeric{PointBase: PointBase{typeName: "int16"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int32":      &TypeInfoNumeric{PointBase: PointBase{typeName: "int32"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int64":      &TypeInfoNumeric{PointBase: PointBase{typeName: "int64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint":       &TypeInfoNumeric{PointBase: PointBase{typeName: "uint"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint8":      &TypeInfoNumeric{PointBase: PointBase{typeName: "uint8"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint16":     &TypeInfoNumeric{PointBase: PointBase{typeName: "uint16"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint32":     &TypeInfoNumeric{PointBase: PointBase{typeName: "uint32"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint64":     &TypeInfoNumeric{PointBase: PointBase{typeName: "uint64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uintptr":    &TypeInfoNumeric{PointBase: PointBase{typeName: "uintptr"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	// "error":      &TypeInfoStruct{PointBase: PointBase{typeName: "error"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"string": &TypeInfoString{PointBase: PointBase{typeName: "string"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"rune":   &TypeInfoString{PointBase: PointBase{typeName: "rune"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
}
