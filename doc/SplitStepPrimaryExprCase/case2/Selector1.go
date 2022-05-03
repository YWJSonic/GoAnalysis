time.Now().Unix()

// 從左邊開始解析
time.Now
// PackageName "." identifier
// => QualifiedIdent
// => PrimaryExpr

time.Now()
// PrimaryExpr "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")"
// => PrimaryExpr Arguments
// => PrimaryExpr

time.Now().Unix
// PrimaryExpr "." identifier
// => PrimaryExpr Selector
// => PrimaryExpr

time.Now().Unix()
// PrimaryExpr "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")"
// => PrimaryExpr Arguments
// => PrimaryExpr