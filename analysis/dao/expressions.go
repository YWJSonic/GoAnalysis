package dao

type Expression struct {
	ExpressionType ITypeInfo // 表達式最終類型
	ContentStr     string

	SubExpression []Expression // 子表達式
	Operators     string       // 操作符號
	PrimaryExpr   *PrimaryExpr // 初級表達式
}

func NewExpression() *Expression {
	return &Expression{}
}

// 初始化表達式
type PrimaryExpr struct {
	Operand    *Operand    // 數值操作
	Conversion *Conversion // 型態轉換
	MethodExpr *MethodExpr // 此包內的 結構方法

	SubPrimaryExpr *PrimaryExpr   // 子物件關聯父物件
	Arguments      *Arguments     // 子物件操作參數 f(Argument, Argument)
	Selector       *Selector      // 子物件操作方法
	Index          *Index         // 子物件索引
	Slice          *Slice         // 子物件切片範圍
	TypeAssertion  *TypeAssertion // 子物件類型斷言
}

type Arguments struct {
	SubExpression *Expression
}

type Selector struct {
	TypeInfoFunc *TypeInfoQualifiedIdent
}

type Index struct {
	Expression *Expression
}

type Slice struct {
	Expression []Expression
}

type TypeAssertion struct {
	Typeinfo ITypeInfo
}

type Conversion struct {
	ConversionType ITypeInfo
	SubExpression  *Expression
}

type MethodExpr struct {
	ReceiverType *TypeInfo
	Method       *FuncInfo
}

// 操作參數
// A + B => A,B 為操作參數
type Operand struct {
	OperandStrName string
	OperandName    ITypeInfo   // 操作變數 const, var, func, package.const, package.var, package.func
	Literal        *Literal    // 操作參數文字
	SubExpression  *Expression // 表達式
}

// 操作參數 文字類型
type Literal struct {
	BasicLit     ITypeInfo     // 基礎類型
	CompositeLit *CompositeLit // 組合類型
	FunctionLit  *FunctionLit  // 函數類型
}

type FunctionLit struct {
	BasicLit       ITypeInfo    // 基礎類型
	ParamsInPoint  []FuncParams // 輸入參數
	ParamsOutPoint []FuncParams // 輸出參數
	Body           string       `json:"-"` // 方法內文 *尚未解析
}

// CompositeLit 組合類型
// []int{5,8,11,3} => LiteralType: []int, LiteralValue: {5,8,11,3}
type CompositeLit struct {
	LiteralType  ITypeInfo    // 型態宣告
	LiteralValue *ElementList // 數值
}

type ElementList struct {
	KeyedElements []KeyedElement
}

type KeyedElement struct {
	SubExpressionKey *Key
	SubElement       *Element
}

type Key struct {
	SubFieldName    string       // 字符
	SubExpression   *Expression  // 表達式
	SubLiteralValue *ElementList // 數值
}

type Element struct {
	SubExpression   *Expression  // 表達式
	SubLiteralValue *ElementList // 數值
}
