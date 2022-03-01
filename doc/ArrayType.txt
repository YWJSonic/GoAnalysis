// Types =========================================================================
Type      = TypeName | TypeLit | "(" Type ")" .
TypeName  = identifier | QualifiedIdent . (QualifiedIdent = PackageName "." identifier .) (PackageName    = identifier .)
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
	    SliceType | MapType | ChannelType .

// ArrayType declarations =========================================================================
ArrayType   = "[" ArrayLength "]" ElementType .
ArrayLength = Expression .
ElementType = Type .

[32]byte
[2*N] struct { x, y int32 }
[1000]*float64
[3][5]int
[2][2][2]float64  // same as [2]([2]([2]float64))