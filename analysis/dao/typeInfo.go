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
	TypeBase
	DefType         string
	ContentTypeInfo ITypeInfo
}

func (self *TypeInfo) GetTypeName() string {
	return self.GetName()
}

// 別名宣告
func NewTypeAliasDecl() *TypeInfo {
	return &TypeInfo{
		TypeBase: NewPointBase(),
		DefType:  "Decl",
	}
}

// 類型定義
func NewTypeDef() *TypeInfo {
	return &TypeInfo{
		TypeBase: NewPointBase(),
		DefType:  "Def",
	}
}

// struct 類型
type TypeInfoStruct struct {
	TypeBase

	ImplicitlyVarInfos []*VarInfo // 隱藏式宣告參數
	VarInfos           map[string]*VarInfo
	// FuncPoint          map[string]*FuncInfo
}

func (self *TypeInfoStruct) GetTypeName() string {
	return self.GetName()
}
func NewTypeInfoStruct() *TypeInfoStruct {
	return &TypeInfoStruct{
		TypeBase: NewPointBase(),
		VarInfos: make(map[string]*VarInfo),
		// FuncPoint: make(map[string]*FuncInfo),
	}
}

// 指標類型
type TypeInfoPointer struct {
	TypeBase
	ContentTypeInfo ITypeInfo
}

func (self *TypeInfoPointer) GetTypeName() string {
	return fmt.Sprintf("*%s", self.ContentTypeInfo.GetTypeName())
}

func NewTypeInfoPointer() *TypeInfoPointer {
	return &TypeInfoPointer{
		TypeBase: NewPointBase(),
	}
}

// func 類型
type TypeInfoFunction struct {
	TypeBase
	ParamsInPoint  []FuncParams // 輸入參數
	ParamsOutPoint []FuncParams // 輸出參數
	Common         string       // 註解
}

func NewTypeInfoFunction() *TypeInfoFunction {
	return &TypeInfoFunction{
		TypeBase: NewPointBase(),
	}
}

// interface 類型
type TypeInfoInterface struct {
	TypeBase
	MatchInfos []ITypeInfo // interface 定義的方法
}

func NewTypeInfoInterface() *TypeInfoInterface {
	return &TypeInfoInterface{
		TypeBase: NewPointBase(),
	}
}

// [] 類型
type TypeInfoSlice struct {
	TypeBase
	ContentTypeInfo ITypeInfo // 資料型別
}

func (self *TypeInfoSlice) GetTypeName() string {
	return fmt.Sprintf("[]%s", self.ContentTypeInfo.GetTypeName())
}

func NewTypeInfoSlice() *TypeInfoSlice {
	return &TypeInfoSlice{
		TypeBase: NewPointBase(),
	}
}

// [n] 類型
type TypeInfoArray struct {
	TypeBase
	Size            string    // 陣列大小
	ContentTypeInfo ITypeInfo // 資料型別
}

func (self *TypeInfoArray) GetTypeName() string {
	return fmt.Sprintf("[%s]%s", self.Size, self.ContentTypeInfo.GetTypeName())
}

func NewTypeInfoArray() *TypeInfoArray {
	return &TypeInfoArray{
		TypeBase: NewPointBase(),
	}
}

// map 類型
type TypeInfoMap struct {
	TypeBase
	KeyType   ITypeInfo
	ValueType ITypeInfo
}

func NewTypeInfoMap() *TypeInfoMap {
	return &TypeInfoMap{
		TypeBase: NewPointBase(),
	}
}

// chan 類型
type TypeInfoChannel struct {
	TypeBase
	FlowType        int       // channel 類型: 雙向, 單出, 單進
	ContentTypeInfo ITypeInfo // channel 傳輸型別
}

func NewTypeInfoChannel() *TypeInfoChannel {
	return &TypeInfoChannel{
		TypeBase: NewPointBase(),
	}
}

// 數值基礎類型
type TypeInfoNumeric struct {
	TypeBase
	RefBase
}

// string 類型
type TypeInfoString struct {
	TypeBase
	RefBase
}

// bool 類型
type TypeInfoBool struct {
	TypeBase
	RefBase
}

// package.A 類型
type TypeInfoQualifiedIdent struct {
	TypeBase
	RefBase
	ImportLink      *ImportInfo // 指定的 import package
	ContentTypeInfo ITypeInfo   // 指定該包的 type
}

func (self *TypeInfoQualifiedIdent) GetTypeName() string {
	return self.GetName()
}

func NewTypeInfoQualifiedIdent() *TypeInfoQualifiedIdent {
	return &TypeInfoQualifiedIdent{
		TypeBase: NewPointBase(),
	}
}

// 基礎類型表
var BaseTypeInfo map[string]ITypeInfo = map[string]ITypeInfo{
	"bool":       &TypeInfoBool{TypeBase: &PointBase{TypeName: "bool"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"byte":       &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "byte"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"complex64":  &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "complex64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"complex128": &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "complex128"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"float32":    &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "float32"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"float64":    &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "float64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int":        &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int8":       &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int8"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"iota":       &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int16":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int16"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int32":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int32"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"int64":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint":       &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint8":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint8"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint16":     &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint16"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint32":     &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint32"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uint64":     &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint64"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"uintptr":    &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uintptr"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	// "error":      &TypeInfoStruct{TypeBase: TypeBase{typeName: "error"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"string": &TypeInfoString{TypeBase: &PointBase{TypeName: "string"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
	"rune":   &TypeInfoString{TypeBase: &PointBase{TypeName: "rune"}, RefBase: RefBase{TypeFrom: constant.From_Golang}},
}
