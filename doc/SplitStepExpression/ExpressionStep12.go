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

cs{
	A: 5,
	struct{	B int}:{B: 7},
}

LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = identifier | Expression | LiteralValue .
Element       = Expression | LiteralValue .

decimal_digits		= ("0" … "9") { [ "_" ] ("0" … "9") } .
binary_digits		= ( "0" | "1" ) { [ "_" ] ( "0" | "1" ) } .
octal_digits		= ( "0" … "7" ) { [ "_" ] ( "0" … "7" ) } .
hex_digits			= ( "0" … "9" | "A" … "F" | "a" … "f" ) { [ "_" ] ( "0" … "9" | "A" … "F" | "a" … "f" ) } .
octal_byte_value	= `\` ( "0" … "7" ) ( "0" … "7" ) ( "0" … "7" ) .
hex_byte_value		= `\` "x" ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) .
little_u_value		= `\` "u" ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) .
big_u_value			= `\` "U" ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) ( "0" … "9" | "A" … "F" | "a" … "f" ) .
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