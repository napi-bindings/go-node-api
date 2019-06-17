package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo CXXFLAGS:  -I./include/
#cgo CFLAGS: -I./include/
#cgo LDFLAGS: -L./lib/ -lnode_api
#include <stdlib.h>
#include "gonapi.h"
#include <node_api.h>
*/
import "C"
import "unsafe"

// JavaScript API

// Env represents the opaque data structure containing the environment in which
// the request is being run.
type Env struct{}

// Value is a the base data type upon which other JavaScript values such as
// Number, Boolean, String, and Object are based.
type Value struct{}

// Boolean class is a representation of the JavaScript Boolean object.
type Boolean struct{}

// Number is a representation of the JavaScript Number object.
type Number struct{}

// String is a representation of the JavaScript String object.
type String struct{}

// Name is a representation of the JavaScript String or Symbol object.
type Name struct{}

// Symbol is a representation of the JavaScript Symbol object.
type Symbol struct{}

// Object is
type Object struct{}

// Array iss
type Array struct{}

// Function is
type Function struct{}

// Buffer is
type Buffer struct{}

// Error is
type Error struct{}

// PropertyDescriptor is
type PropertyDescriptor struct{}

// ClassPropertyDescriptor is
type ClassPropertyDescriptor struct{}

// ObjectPropertyDescriptor is
type ObjectPropertyDescriptor struct{}

// CallbackInfo is
type CallbackInfo struct{}

// Reference is
type Reference struct{}

// TypedArray is
type TypedArray struct{}

// ArrayBuffer is
type ArrayBuffer struct{}

// DataView is
type DataView struct{}

// Promise is
type Promise struct{}

// ObjectReference is
type ObjectReference struct{}

// FunctionReference is
type FunctionReference struct{}

// ObjectWrap iss
type ObjectWrap struct{}

// HandleScope is
type HandleScope struct{}

// EscapableHandleScope iss
type EscapableHandleScope struct{}

// CallbackScope is
type CallbackScope struct{}

// MemoryManagement is
type MemoryManagement struct{}

// Runtime API

// AsyncContext is
type AsyncContext struct{}

// AsyncWorker is
type AsyncWorker struct{}

// Caller contains a callback to call
type Caller struct{}

func (s *Caller) cb(env C.napi_env, info C.napi_callback_info) C.napi_value {
	var res C.napi_value
	C.napi_create_int32(env, C.int(5), &res)
	return res
}

//export ExecuteCallback
func ExecuteCallback(data unsafe.Pointer, env C.napi_env, info C.napi_callback_info) C.napi_value {
	caller := (*Caller)(data)
	return caller.cb(env, info)
}

//export Initialize
func Initialize(env C.napi_env, exports C.napi_value) C.napi_value {
	name := C.CString("createInt32")
	defer C.free(unsafe.Pointer(name))
	caller := &Caller{}
	desc := C.napi_property_descriptor{
		utf8name:   name,
		name:       nil,
		method:     (C.napi_callback)(C.CallbackMethod(unsafe.Pointer(&caller))), //nil,
		getter:     nil,
		setter:     nil,
		value:      nil,
		attributes: C.napi_default,
		data:       nil,
	}
	C.napi_define_properties(env, exports, 1, &desc)
	return exports
}

func main() {}
