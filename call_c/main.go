package main

// #cgo CFLAGS: -I./foo
// #cgo LDFLAGS: -L./foo -lfoo
// #include "foo.h"
import "C"
import "fmt"

func main() {
    value := C.add(C.int(1), C.int(2))
    fmt.Printf("%v\n", value)
}
