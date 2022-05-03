Type = 
identifier | 
PackageName "." identifier | 
"[" Expression "]" Type | 
"struct" "{" { FieldDecl ";" } "}" | 
"*" Type | 
"func" Signature | 
"interface" "{" { ( MethodSpec | InterfaceTypeName ) ";" } "}" |
"[" "]" Type | 
"map" "[" Type "]" Type | 
( "chan" | "chan" "<-" | "<-" "chan" ) Type | 
"(" Type ")" .