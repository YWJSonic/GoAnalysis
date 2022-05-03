// Instantiations =========================================================================
// Making slices, maps and channels
// Call          Core type    Result
make(T, n)       slice        slice of type T with length n and capacity n
make(T, n, m)    slice        slice of type T with length n and capacity m
make(T)          map          map of type T
make(T, n)       map          map of type T with initial space for approximately n elements
make(T)          channel      unbuffered channel of type T
make(T, n)       channel      buffered channel of type T, buffer size n

s := make([]int, 10, 100)       // slice with len(s) == 10, cap(s) == 100
s := make([]int, 1e3)           // slice with len(s) == cap(s) == 1000
s := make([]int, 1<<63)         // illegal: len(s) is not representable by a value of type int
s := make([]int, 10, 0)         // illegal: len(s) > cap(s)
c := make(chan int, 10)         // channel with a buffer size of 10
m := make(map[string]int, 100)  // map with initial space for approximately 100 elements