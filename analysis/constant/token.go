package constant

import "codeanalysis/types"

const (
	From_Golang        types.TypeFrom = 0 // go 系統自帶
	From_Local         types.TypeFrom = 1 // 專案實做
	From_Local_Package types.TypeFrom = 2 // 專案其他 package 實做
	From_Vendor        types.TypeFrom = 3 // 外部引用
)
const (
	DefType_Def  string = "Def"
	DefType_Decl string = "Decl"
)

const (
	_    token = iota
	_EOF       // EOF

	// names and literals
	_Name    // nametokenLit
	_Literal // literal

	// operators and operations
	// _Operator is excluding '*' (_Star)
	_Operator // op
	_AssignOp // op=
	_IncOp    // opop
	_Assign   // =
	_Define   // :=
	_Arrow    // <-
	_Star     // *

	// delimiters
	_Lparen    // (
	_Lbrack    // [
	_Lbrace    // {
	_Rparen    // )
	_Rbrack    // ]
	_Rbrace    // }
	_Comma     // ,
	_Semi      // ;
	_Colon     // :
	_Dot       // .
	_DotDotDot // ...

	// keywords
	_Break       // break
	_Case        // case
	_Chan        // chan
	_Const       // const
	_Continue    // continue
	_Default     // default
	_Defer       // defer
	_Else        // else
	_Fallthrough // fallthrough
	_For         // for
	_Func        // func
	_Go          // go
	_Goto        // goto
	_If          // if
	_Import      // import
	_Interface   // interface
	_Map         // map
	_Package     // package
	_Range       // range
	_Return      // return
	_Select      // select
	_Struct      // struct
	_Switch      // switch
	_Type        // type
	_Var         // var

	// empty line comment to exclude it from .String
	tokenCount //
)

const (
	// for BranchStmt
	Break       = _Break
	Continue    = _Continue
	Fallthrough = _Fallthrough
	Goto        = _Goto

	// for CallStmt
	Go    = _Go
	Defer = _Defer
)

type LitKind = uint8

// TODO(gri) With the 'i' (imaginary) suffix now permitted on integer
//           and floating-point numbers, having a single ImagLit does
//           not represent the literal kind well anymore. Remove it?
const (
	IntLit LitKind = iota
	FloatLit
	ImagLit
	RuneLit
	StringLit
)

type IntLiteType = uint8

const (
	IntLiteType_None    IntLiteType = iota // 不是數字
	IntLiteType_Binary                     // 2進制
	IntLiteType_Octal                      // 8進制
	IntLiteType_Decimal                    // 10進制
	IntLiteType_Hex                        // 16進制
)

type token = uint

//go:generate stringer -type token -linecomment tokens.go

// Make sure we have at most 64 tokens so we can use them in a set.
const _ uint64 = 1 << (tokenCount - 1)

// contains reports whether tok is in tokset.
func contains(tokset uint64, tok token) bool {
	return tokset&(1<<tok) != 0
}

type Operator uint

//go:generate stringer -type Operator -linecomment tokens.go

const (
	_ Operator = iota

	// Def is the : in :=
	Def   // :
	Not   // !
	Recv  // <-
	Tilde // ~

	// precOrOr
	OrOr // ||

	// precAndAnd
	AndAnd // &&

	// precCmp
	Eql // ==
	Neq // !=
	Lss // <
	Leq // <=
	Gtr // >
	Geq // >=

	// precAdd
	Add // +
	Sub // -
	Or  // |
	Xor // ^

	// precMul
	Mul    // *
	Div    // /
	Rem    // %
	And    // &
	AndNot // &^
	Shl    // <<
	Shr    // >>
)

// Operator precedences
const (
	_ = iota
	precOrOr
	precAndAnd
	precCmp
	precAdd
	precMul
)

type TypeString = string

var typeStringLit = []TypeString{
	_bool,
	_byte,
	_complex64,
	_complex128,
	_error,
	_float32,
	_float64,
	_int,
	_int8,
	_int16,
	_int32,
	_int64,
	_rune,
	_string,
	_uint,
	_uint8,
	_uint16,
	_uint32,
	_uint64,
	_uintptr,
}

const (
	_bool       TypeString = "bool"
	_byte       TypeString = "byte"
	_complex64  TypeString = "complex64"
	_complex128 TypeString = "complex128"
	_error      TypeString = "error"
	_float32    TypeString = "float32"
	_float64    TypeString = "float64"
	_int        TypeString = "int"
	_int8       TypeString = "int8"
	_int16      TypeString = "int16"
	_int32      TypeString = "int32"
	_int64      TypeString = "int64"
	_rune       TypeString = "rune"
	_string     TypeString = "string"
	_uint       TypeString = "uint"
	_uint8      TypeString = "uint8"
	_uint16     TypeString = "uint16"
	_uint32     TypeString = "uint32"
	_uint64     TypeString = "uint64"
	_uintptr    TypeString = "uintptr"
)

type KeyWordString = string

var keyWordStringLit = []string{
	_break,
	_default,
	_func,
	_interface,
	_select,
	_case,
	_defer,
	_go,
	_map,
	_struct,
	_chan,
	_else,
	_goto,
	_package,
	_switch,
	_const,
	_fallthrough,
	_if,
	_range,
	_type,
	_continue,
	_for,
	_import,
	_return,
	_var,
}

const (
	_break       KeyWordString = "break"
	_default     KeyWordString = "default"
	_func        KeyWordString = "func"
	_interface   KeyWordString = "interface"
	_select      KeyWordString = "select"
	_case        KeyWordString = "case"
	_defer       KeyWordString = "defer"
	_go          KeyWordString = "go"
	_map         KeyWordString = "map"
	_struct      KeyWordString = "struct"
	_chan        KeyWordString = "chan"
	_else        KeyWordString = "else"
	_goto        KeyWordString = "goto"
	_package     KeyWordString = "package"
	_switch      KeyWordString = "switch"
	_const       KeyWordString = "const"
	_fallthrough KeyWordString = "fallthrough"
	_if          KeyWordString = "if"
	_range       KeyWordString = "range"
	_type        KeyWordString = "type"
	_continue    KeyWordString = "continue"
	_for         KeyWordString = "for"
	_import      KeyWordString = "import"
	_return      KeyWordString = "return"
	_var         KeyWordString = "var"
)

var TokenLit = []byte{'"', '`', '\'', '(', '{', '[', ';', ',', ']', ')', ':', '}', '+', '.', '*', '-', '%', '/', '|', '&', '<', '^', '=', '>', '~', '!'}
