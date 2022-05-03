Expression = 
int_lit | 
float_lit | 
imaginary_lit | 
rune_lit | 
string_lit | 
CompositeLit | 
FunctionLit | 
OperandName | 
"(" Expression ")" |
Type "(" Expression [ "," ] ")" |
ReceiverType "." MethodName |
PrimaryExpr "." identifier |
PrimaryExpr "[" Expression "]" |
PrimaryExpr ( "[" [ Expression ] ":" [ Expression ] "]" | "[" [ Expression ] ":" Expression ":" Expression "]" ) |
PrimaryExpr "." "(" Type ")" |
PrimaryExpr "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")" | 
( "+" | "-" | "!" | "^" | "*" | "&" | "<-" ) ( PrimaryExpr | ( "+" | "-" | "!" | "^" | "*" | "&" | "<-" ) UnaryExpr ) | 
Expression ( "||" | "&&" | "==" | "!=" | "<" | "<=" | ">" | ">=" | "+" | "-" | "|" | "^" | "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" ) Expression .