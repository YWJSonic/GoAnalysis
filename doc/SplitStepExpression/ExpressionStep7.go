Expression = 
decimal_lit | 
binary_lit | 
octal_lit | 
hex_lit | 
decimal_float_lit | 
hex_float_lit | 
( decimal_digits | int_lit | float_lit ) "i" | 
"'" ( unicode_value | byte_value ) "'" | 
raw_string_lit | 
interpreted_string_lit | 
LiteralType LiteralValue | 
"func" Signature FunctionBody | 
identifier | 
QualifiedIdent | 
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