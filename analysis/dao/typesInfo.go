package dao

// 類型物件結構

type ITypeInfo interface {
	SetName(name string)
	GetTypeName() string
}

type TypeInfoArray struct{}
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

type TypeInfoPointer struct {
	PointBase
	ContentTypeInfo ITypeInfo
}

func NewTypeInfoPointer() *TypeInfoPointer {
	return &TypeInfoPointer{}
}

type TypeInfoFunction struct{}
type TypeInfoInterface struct{}
type TypeInfoSlice struct {
	PointBase
	ContentTypeInfo ITypeInfo // 資料型別
}

func NewTypeInfoSlice() *TypeInfoSlice {
	return &TypeInfoSlice{}
}

type TypeInfoMap struct {
	PointBase
	KeyType   ITypeInfo
	ValueType ITypeInfo
}

func NewTypeInfoMap() *TypeInfoMap {
	return &TypeInfoMap{}
}

type TypeInfoChannel struct {
	PointBase
	FlowType        int       // channel 類型: 雙向, 單出, 單進
	ContentTypeInfo ITypeInfo // channel 傳輸型別
}

func NewTypeInfoChannel() *TypeInfoChannel {
	return &TypeInfoChannel{}
}

type TypeInfoNumeric struct {
	PointBase
}

type TypeInfoString struct {
	PointBase
}

type TypeInfoBool struct {
	PointBase
}

type TypeInfoQualifiedIdent struct {
	PointBase
	ImportLink      *PackageLink // 指定的 import package
	ContentTypeInfo ITypeInfo    // 指定該包的 type
}

func NewTypeInfoQualifiedIdent() *TypeInfoQualifiedIdent {
	return &TypeInfoQualifiedIdent{}
}
