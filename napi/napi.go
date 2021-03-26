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
import "fmt"

// Hello1
func Hello1() {
	fmt.Println(C.napi_writable)
}

// Hello function ...
func Hello() string {
	return "Hello, world."
}
