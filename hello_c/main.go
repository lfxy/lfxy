package hello
/*
#cgo CFLAGS: -I../hello_c/foo
#cgo LDFLAGS: -L../hello_c/foo -lhello
#include <hello.h>
*/
import "C"
func Hello(str string) {
    C.hello(C.CString(str))
}
