package dao

import (
	"codeanalysis/analysis/constant"
	"encoding/json"
	"fmt"
)

// Type declarations 類型聲明結構

type ITypeInfo interface {
	GetTypeBase() TypeBase
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

func (self *TypeInfo) MarshalJSON() ([]byte, error) {
	tmp := struct {
		TypeBase
		ConstantTypeBase TypeBase
	}{
		TypeBase: self.TypeBase,
	}

	if self.ContentTypeInfo != nil {
		tmp.ConstantTypeBase = self.ContentTypeInfo.GetTypeBase()
	}

	return json.Marshal(tmp)
}

func (self *TypeInfo) GetTypeBase() TypeBase {
	return self.TypeBase
}

func (self *TypeInfo) GetTypeName() string {
	return self.GetName()
}

// 別名宣告
func NewTypeAliasDecl() *TypeInfo {
	return &TypeInfo{
		TypeBase: NewPointBase(),
		DefType:  constant.DefType_Decl,
	}
}

// 類型定義
func NewTypeDef() *TypeInfo {
	return &TypeInfo{
		TypeBase: NewPointBase(),
		DefType:  constant.DefType_Def,
	}
}

type TypeInfoRef struct {
	TypeBase
}

func (self *TypeInfoRef) GetTypeBase() TypeBase {
	return self.TypeBase
}

func NewTypeInfoRef() *TypeInfoRef {
	return &TypeInfoRef{}
}

// struct 類型
type TypeInfoStruct struct {
	TypeBase

	ImplicitlyVarInfos []*VarInfo // 隱藏式宣告參數
	VarInfos           map[string]*VarInfo
	// FuncPoint          map[string]*FuncInfo
}

func (self *TypeInfoStruct) GetTypeBase() TypeBase {
	return self.TypeBase
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

func (self *TypeInfoPointer) GetTypeBase() TypeBase {
	return self.TypeBase
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

func (self *TypeInfoFunction) GetTypeBase() TypeBase {
	return self.TypeBase
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

func (self *TypeInfoInterface) GetTypeBase() TypeBase {
	return self.TypeBase
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

func (self *TypeInfoSlice) GetTypeBase() TypeBase {
	return self.TypeBase
}

func (self *TypeInfoSlice) GetTypeName() string {
	return "[]"
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

func (self *TypeInfoArray) GetTypeBase() TypeBase {
	return self.TypeBase
}

func (self *TypeInfoArray) GetTypeName() string {
	return fmt.Sprintf("[%s]", self.Size)
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

func (self *TypeInfoMap) GetTypeBase() TypeBase {
	return self.TypeBase
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

func (self *TypeInfoChannel) GetTypeBase() TypeBase {
	return self.TypeBase
}

func NewTypeInfoChannel() *TypeInfoChannel {
	return &TypeInfoChannel{
		TypeBase: NewPointBase(),
	}
}

// 數值基礎類型
type TypeInfoNumeric struct {
	TypeBase
}

func (self *TypeInfoNumeric) GetTypeBase() TypeBase {
	return self.TypeBase
}

// string 類型
type TypeInfoString struct {
	TypeBase
}

func (self *TypeInfoString) GetTypeBase() TypeBase {
	return self.TypeBase
}

// bool 類型
type TypeInfoBool struct {
	TypeBase
}

func (self *TypeInfoBool) GetTypeBase() TypeBase {
	return self.TypeBase
}

// package.A 類型
type TypeInfoQualifiedIdent struct {
	TypeBase
	ImportLink      *ImportInfo // 指定的 import package
	ContentTypeInfo ITypeInfo   // 指定該包的 itype
}

func (self *TypeInfoQualifiedIdent) GetTypeBase() TypeBase {
	return self.TypeBase
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
	"bool":       &TypeInfoBool{TypeBase: &PointBase{TypeName: "bool", TypeFrom: constant.From_Golang, GoPath: ""}},
	"byte":       &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "byte", TypeFrom: constant.From_Golang}},
	"complex64":  &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "complex64", TypeFrom: constant.From_Golang}},
	"complex128": &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "complex128", TypeFrom: constant.From_Golang}},
	"float32":    &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "float32", TypeFrom: constant.From_Golang}},
	"float64":    &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "float64", TypeFrom: constant.From_Golang}},
	"int":        &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int", TypeFrom: constant.From_Golang}},
	"int8":       &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int8", TypeFrom: constant.From_Golang}},
	"int16":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int16", TypeFrom: constant.From_Golang}},
	"int32":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int32", TypeFrom: constant.From_Golang}},
	"int64":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int64", TypeFrom: constant.From_Golang}},
	"uint":       &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint", TypeFrom: constant.From_Golang}},
	"uint8":      &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint8", TypeFrom: constant.From_Golang}},
	"uint16":     &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint16", TypeFrom: constant.From_Golang}},
	"uint32":     &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint32", TypeFrom: constant.From_Golang}},
	"uint64":     &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uint64", TypeFrom: constant.From_Golang}},
	"uintptr":    &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "uintptr", TypeFrom: constant.From_Golang}},
	"error":      &TypeInfoStruct{TypeBase: &PointBase{TypeName: "error", TypeFrom: constant.From_Golang}},
	"string":     &TypeInfoString{TypeBase: &PointBase{TypeName: "string", TypeFrom: constant.From_Golang}},
	"rune":       &TypeInfoString{TypeBase: &PointBase{TypeName: "rune", TypeFrom: constant.From_Golang}},
}

var BaseConstInfo map[string]ITypeInfo = map[string]ITypeInfo{
	"iota":  &TypeInfoNumeric{TypeBase: &PointBase{TypeName: "int", TypeFrom: constant.From_Golang}},
	"true":  &TypeInfoBool{TypeBase: &PointBase{TypeName: "true", TypeFrom: constant.From_Golang, GoPath: ""}},
	"false": &TypeInfoBool{TypeBase: &PointBase{TypeName: "false", TypeFrom: constant.From_Golang, GoPath: ""}},
}

var BaseFuncInfo map[string]ITypeInfo = map[string]ITypeInfo{
	"append":     &TypeInfoFunction{TypeBase: &PointBase{TypeName: "append", TypeFrom: constant.From_Golang}},
	"cap":        &TypeInfoFunction{TypeBase: &PointBase{TypeName: "cap", TypeFrom: constant.From_Golang}},
	"close":      &TypeInfoFunction{TypeBase: &PointBase{TypeName: "close", TypeFrom: constant.From_Golang}},
	"comparable": &TypeInfoFunction{TypeBase: &PointBase{TypeName: "comparable", TypeFrom: constant.From_Golang}},
	"copy":       &TypeInfoFunction{TypeBase: &PointBase{TypeName: "copy", TypeFrom: constant.From_Golang}},
	"delete":     &TypeInfoFunction{TypeBase: &PointBase{TypeName: "delete", TypeFrom: constant.From_Golang}},
	"imag":       &TypeInfoFunction{TypeBase: &PointBase{TypeName: "imag", TypeFrom: constant.From_Golang}},
	"len":        &TypeInfoFunction{TypeBase: &PointBase{TypeName: "len", TypeFrom: constant.From_Golang}},
	"make":       &TypeInfoFunction{TypeBase: &PointBase{TypeName: "make", TypeFrom: constant.From_Golang}},
	"new":        &TypeInfoFunction{TypeBase: &PointBase{TypeName: "new", TypeFrom: constant.From_Golang}},
	"panic":      &TypeInfoFunction{TypeBase: &PointBase{TypeName: "panic", TypeFrom: constant.From_Golang}},
	"print":      &TypeInfoFunction{TypeBase: &PointBase{TypeName: "print", TypeFrom: constant.From_Golang}},
	"println":    &TypeInfoFunction{TypeBase: &PointBase{TypeName: "println", TypeFrom: constant.From_Golang}},
	"real":       &TypeInfoFunction{TypeBase: &PointBase{TypeName: "real", TypeFrom: constant.From_Golang}},
	"recover":    &TypeInfoFunction{TypeBase: &PointBase{TypeName: "recover", TypeFrom: constant.From_Golang}},
}
