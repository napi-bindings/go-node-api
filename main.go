package main

import (
	"C"
)
import (
	"time"
	"unsafe"

	"github.com/napi-bindings/go-node-api/napi"
)

func unixNano(env napi.Env, info napi.CallbackInfo) napi.Value {
	now := time.Now()
	value, _ := napi.CreateInt64(env, now.UnixNano())
	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	caller := &napi.Caller{
		Cb: unixNano,
	}
	desc := napi.Property{
		Name:   "unixNano",
		Method: caller,
	}
	props := []napi.Property{desc}
	napi.DefineProperties((napi.Env)(env), (napi.Value)(exports), props)
	return exports
}

func main() {}
