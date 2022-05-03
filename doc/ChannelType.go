// Types =========================================================================
Type      = TypeName | TypeLit | "(" Type ")" .
TypeName  = identifier | QualifiedIdent . (QualifiedIdent = PackageName "." identifier .) (PackageName    = identifier .)
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
	    SliceType | MapType | ChannelType .

// Channel types =========================================================================
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .(ElementType = Type .)

chan T          // can be used to send and receive values of type T
chan<- float64  // can only be used to send float64s
<-chan int      // can only be used to receive ints
chan<- chan int    // same as chan<- (chan int)
chan<- <-chan int  // same as chan<- (<-chan int)
<-chan <-chan int  // same as <-chan (<-chan int)
chan (<-chan int)
