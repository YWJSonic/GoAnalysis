f.p[i].x()

// 從左邊開始解析
f.p
// PackageName "." identifier
// => QualifiedIdent
// => PrimaryExpr

f.p[i]
// PrimaryExpr "[" Expression "]"
// => PrimaryExpr Index
// => PrimaryExpr

f.p[i].x
// PrimaryExpr "." identifier
// => PrimaryExpr Selector
// => PrimaryExpr

f.p[i].x()
// PrimaryExpr "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")"
// => PrimaryExpr Arguments
// => PrimaryExpr