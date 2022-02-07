package dao

// Type declarations 類型聲明結構

type ITypeDeclare interface{}

// 別名宣告

type TypeAliasDecl struct {
	PointBase
}

func NewTypeAliasDecl(name string) *TypeAliasDecl {
	info := &TypeAliasDecl{}
	info.Name = name
	return info
}

// 類型定義

type TypeDef struct {
	PointBase
}

func NewTypeDef(name string) *TypeDef {
	info := &TypeDef{}
	info.Name = name
	return info
}

var BaseTypeInfo map[string]ITypeInfo = map[string]ITypeInfo{
	"bool":       &TypeInfoBool{PointBase: PointBase{TypeName: "bool"}},
	"byte":       &TypeInfoNumeric{PointBase: PointBase{TypeName: "byte"}},
	"complex64":  &TypeInfoNumeric{PointBase: PointBase{TypeName: "complex64"}},
	"complex128": &TypeInfoNumeric{PointBase: PointBase{TypeName: "complex128"}},
	"float32":    &TypeInfoNumeric{PointBase: PointBase{TypeName: "float32"}},
	"float64":    &TypeInfoNumeric{PointBase: PointBase{TypeName: "float64"}},
	"int":        &TypeInfoNumeric{PointBase: PointBase{TypeName: "int"}},
	"int8":       &TypeInfoNumeric{PointBase: PointBase{TypeName: "int8"}},
	"int16":      &TypeInfoNumeric{PointBase: PointBase{TypeName: "int16"}},
	"int32":      &TypeInfoNumeric{PointBase: PointBase{TypeName: "int32"}},
	"int64":      &TypeInfoNumeric{PointBase: PointBase{TypeName: "int64"}},
	"uint":       &TypeInfoNumeric{PointBase: PointBase{TypeName: "uint"}},
	"uint8":      &TypeInfoNumeric{PointBase: PointBase{TypeName: "uint8"}},
	"uint16":     &TypeInfoNumeric{PointBase: PointBase{TypeName: "uint16"}},
	"uint32":     &TypeInfoNumeric{PointBase: PointBase{TypeName: "uint32"}},
	"uint64":     &TypeInfoNumeric{PointBase: PointBase{TypeName: "uint64"}},
	"uintptr":    &TypeInfoNumeric{PointBase: PointBase{TypeName: "uintptr"}},
	"error":      &TypeInfoStruct{PointBase: PointBase{TypeName: "error"}},
	"string":     &TypeInfoString{PointBase: PointBase{TypeName: "string"}},
	"rune":       &TypeInfoString{PointBase: PointBase{TypeName: "rune"}},
}
