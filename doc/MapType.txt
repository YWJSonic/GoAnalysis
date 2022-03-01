// Types =========================================================================
Type      = TypeName | TypeLit | "(" Type ")" .
TypeName  = identifier | QualifiedIdent . (QualifiedIdent = PackageName "." identifier .) (PackageName    = identifier .)
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
	    SliceType | MapType | ChannelType .

// Map types =========================================================================
MapType     = "map" "[" KeyType "]" ElementType .(ElementType = Type .)
KeyType     = Type .

map[string]int
map[*T]struct{ x, y float64 }
map[string]interface{}