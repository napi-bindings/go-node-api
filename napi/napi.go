package napi

/*
#cgo CXXFLAGS: -std=c++11
#cgo CXXFLAGS:  -I${SRCDIR}/deps/include
#cgo CFLAGS: -I${SRCDIR}/deps/include -DNAPI_EXPERIMENTAL
#cgo darwin LDFLAGS: -L${SRCDIR}/deps/lib/darwin
#cgo linux LDFLAGS: -L${SRCDIR}/deps/lib/linux
#cgo windows LDFLAGS: -L${SRCDIR}/deps/lib/windows
#cgo linux LDFLAGS: -Wl,-unresolved-symbols=ignore-all
#cgo darwin LDFLAGS: -Wl,-undefined,dynamic_lookup
#cgo LDFLAGS: -lnode_api
#include <stdlib.h>
#include "gonapi.h"
#include <node_api.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ccstring function transforms a Go string into a C string (array of characters)
// and returns the pointer to the first element.
func cstring(s string) unsafe.Pointer {
	p := make([]byte, len(s)+1)
	copy(p, s)
	return unsafe.Pointer(&p[0])
}

type napi_env = C.napi_env

// Env is used to represent a context that the underlying Node-API
// implementation can use to persist VM-specific state. This structure is
// passed to native functions when they're invoked, and it must be passed back
// when making Node-API calls. Specifically, the same Env that was passed in
// when the initial native function was called must be passed to any subsequent
// nested Node-API calls.
type Env struct{}

// TODO define fields and methods for Env

// Hello1
func Hello1() {
	fmt.Println(C.napi_writable)
}

// Hello function ...
func Hello() string {
	return "Hello, world."
}
