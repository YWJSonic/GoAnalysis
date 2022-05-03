Expression = 
"0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ] | 
"0" ( "b" | "B" ) [ "_" ] binary_digits | 
"0" [ "o" | "O" ] [ "_" ] octal_digits | 
"0" ( "x" | "X" ) [ "_" ] hex_digits | 
decimal_digits "." [ decimal_digits ] [ decimal_exponent ] | decimal_digits decimal_exponent | "." decimal_digits [ decimal_exponent ] | 
"0" ( "x" | "X" ) hex_mantissa hex_exponent | 
( decimal_digits | decimal_lit | binary_lit | octal_lit | hex_lit | decimal_float_lit | hex_float_lit ) "i" | 
"'" ( unicode_char | little_u_value | big_u_value | escaped_char | octal_byte_value | hex_byte_value ) "'" | 
"`" { unicode_char | newline } "`" | 
`"` { unicode_value | byte_value } `"` | 
( StructType | ArrayType | "[" "..." "]" Expression | SliceType | MapType | TypeName ) LiteralValue | 
"func" Signature FunctionBody | 
identifier | 
PackageName "." identifier | 
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