Expression = 
("0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ]) | 
"0" ( "b" | "B" ) [ "_" ] binary_digits | 
"0" [ "o" | "O" ] [ "_" ] octal_digits | 
"0" ( "x" | "X" ) [ "_" ] hex_digits | 
decimal_digits "." [ decimal_digits ] [ ( ( "e" | "E" ) [ "+" | "-" ] decimal_digits ) ] | decimal_digits ( ( "e" | "E" ) [ "+" | "-" ] decimal_digits ) | "." decimal_digits [ ( ( "e" | "E" ) [ "+" | "-" ] decimal_digits ) ] | 
"0" ( "x" | "X" ) ( [ "_" ] hex_digits "." [ hex_digits ] | [ "_" ] hex_digits | "." hex_digits ) ( ( "p" | "P" ) [ "+" | "-" ] decimal_digits ) | 
( decimal_digits | decimal_lit | binary_lit | octal_lit | hex_lit | decimal_float_lit | hex_float_lit ) "i" | 
"'" ( unicode_char | little_u_value | big_u_value | ( `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) ) | octal_byte_value | hex_byte_value ) "'" | 
"`" { unicode_char | newline } "`" | 
`"` { unicode_char | little_u_value | big_u_value | ( `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) ) | octal_byte_value | hex_byte_value } `"` | 
( "struct" "{" { FieldDecl ";" } "}" | "[" Expression "]" Type | "[" "..." "]" Type | "[" "]" Type | "map" "[" KeyType "]" Type | ( identifier | identifier "." identifier ) ) "{" [ ElementList [ "," ] ] "}" | 
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


LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .

decimal_digits		= decimal_digit { [ "_" ] decimal_digit } .
binary_digits		= binary_digit { [ "_" ] binary_digit } .
octal_digits		= octal_digit { [ "_" ] octal_digit } .
hex_digits			= hex_digit { [ "_" ] hex_digit } .
octal_byte_value	= `\` octal_digit octal_digit octal_digit .
hex_byte_value		= `\` "x" hex_digit hex_digit .
little_u_value		= `\` "u" hex_digit hex_digit hex_digit hex_digit .
big_u_value			= `\` "U" hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit .
decimal_lit			= "0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ] .
binary_lit			= "0" ( "b" | "B" ) [ "_" ] binary_digits .
octal_lit			= "0" [ "o" | "O" ] [ "_" ] octal_digits .
hex_lit				= "0" ( "x" | "X" ) [ "_" ] hex_digits .
newline        = /* the Unicode code point U+000A */ .
unicode_char   = /* an arbitrary Unicode code point except newline */ .
unicode_letter = /* a Unicode code point classified as "Letter" */ .
unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .

identifier = letter { letter | unicode_digit } .
letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .