Type = 
identifier | 
identifier "." identifier | 
"[" Expression "]" Type | 
"struct" "{" { (identifier { "," identifier } Type | [ "*" ] ( identifier | QualifiedIdent )) [ Tag ] ";" } "}" | 
"*" Type | 
"func" Signature | 
"interface" "{" { ( identifier Signature | ( identifier | QualifiedIdent ) ) ";" } "}" |
"[" "]" Type | 
"map" "[" Type "]" Type | 
( "chan" | "chan" "<-" | "<-" "chan" ) Type | 
"(" Type ")" .