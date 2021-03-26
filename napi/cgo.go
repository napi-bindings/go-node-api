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

// These imports forces `go mod vendor` to pull in all the folders that
// contain Node-API libraries and headers which otherwise would be ignored.
// DO NOT REMOVE
import (
	_ "github.com/napi-bindings/go-node-api/napi/deps/include"
	_ "github.com/napi-bindings/go-node-api/napi/deps/lib/darwin"
	_ "github.com/napi-bindings/go-node-api/napi/deps/lib/linux"
	_ "github.com/napi-bindings/go-node-api/napi/deps/lib/windows"
)
