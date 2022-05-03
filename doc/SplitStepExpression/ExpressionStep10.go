Expression = 
"0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ] | 
"0" ( "b" | "B" ) [ "_" ] binary_digits | 
"0" [ "o" | "O" ] [ "_" ] octal_digits | 
"0" ( "x" | "X" ) [ "_" ] hex_digits | 
decimal_digits "." [ decimal_digits ] [ ( ( "e" | "E" ) [ "+" | "-" ] decimal_digits ) ] | decimal_digits ( ( "e" | "E" ) [ "+" | "-" ] decimal_digits ) | "." decimal_digits [ ( ( "e" | "E" ) [ "+" | "-" ] decimal_digits ) ] | 
"0" ( "x" | "X" ) ( [ "_" ] hex_digits "." [ hex_digits ] | [ "_" ] hex_digits | "." hex_digits ) ( ( "p" | "P" ) [ "+" | "-" ] decimal_digits ) | 
( decimal_digits | decimal_lit | binary_lit | octal_lit | hex_lit | decimal_float_lit | hex_float_lit ) "i" | 
"'" ( unicode_char | little_u_value | big_u_value | ( `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) ) | octal_byte_value | hex_byte_value ) "'" | 
"`" { unicode_char | newline } "`" | 
`"` { unicode_char | little_u_value | big_u_value | ( `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) ) | octal_byte_value | hex_byte_value } `"` | 
( "struct" "{" { FieldDecl ";" } "}" | "[" Expression "]" Type | "[" "..." "]" Type | "[" "]" Type | "map" "[" KeyType "]" Type | ( identifier | PackageName "." identifier ) ) LiteralValue | 
"func" Signature FunctionBody | 
identifier | 
identifier "." identifier | 
"(" Expression ")" |
Type "(" Expression [ "," ] ")" |
Type "." identifier |
PrimaryExpr "." identifier |
PrimaryExpr "[" Expression "]" |
PrimaryExpr ( "[" [ Expression ] ":" [ Expression ] "]" | "[" [ Expression ] ":" Expression ":" Expression "]" ) |
PrimaryExpr "." "(" Type ")" |
PrimaryExpr "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")" | 
( "+" | "-" | "!" | "^" | "*" | "&" | "<-" ) ( PrimaryExpr | ( "+" | "-" | "!" | "^" | "*" | "&" | "<-" ) UnaryExpr ) | 
Expression ( "||" | "&&" | "==" | "!=" | "<" | "<=" | ">" | ">=" | "+" | "-" | "|" | "^" | "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" ) Expression .