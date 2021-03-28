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

type env = C.napi_env

// Env ...
type Env struct{}

// TODO define fields and methods for Env

type value = C.napi_value

// Value ...
type Value struct{}

// TODO define fields and methods for Value

type ref = C.napi_ref

// Reference ...
type Reference struct{}

// TODO define fields and methods for Reference

type handle_scope = C.napi_handle_scope

// HandleScope
type HandleScope struct{}

// TODO define fields and methods for HandleScope

type escapable_handle_scope = C.napi_escapable_handle_scope

// EscapableHanldeScope ...
type EscapableHanldeScope struct{}

// TODO define fields and methods for EscapableHandleScope

type callback_info = C.napi_callback_info

// CallbackInfo ...
type CallbackInfo struct{}

// TODO define fields and methods for CallbackInfo

type deferred = C.napi_deferred

// Deferred ...
type Deferred struct{}

// TODO define fields and methods for Deferred

// PropertyAttribute ...
type PropertyAttribute int

const (
	DefaultProperty       PropertyAttribute = C.napi_default
	WritableProperty      PropertyAttribute = C.napi_writable
	EnumerableProperty    PropertyAttribute = C.napi_enumerable
	ConfigurableProperty  PropertyAttribute = C.napi_configurable
	StaticProperty        PropertyAttribute = C.napi_static
	DefaultMethodProperty PropertyAttribute = C.napi_default_method
	DefaultJSProperty     PropertyAttribute = C.napi_default_jsproperty
)

// ValueType ...
type ValueType int

const (
	UndefinedType ValueType = C.napi_undefined
	NullType      ValueType = C.napi_null
	BooleanType   ValueType = C.napi_boolean
	NumberType    ValueType = C.napi_number
	StringType    ValueType = C.napi_string
	SymbolType    ValueType = C.napi_symbol
	ObjectType    ValueType = C.napi_object
	FunctionType  ValueType = C.napi_function
	ExternalType  ValueType = C.napi_external
	BigIntType    ValueType = C.napi_bigint
)

// ArrayType ...
type ArrayType int

const (
	Int8ArrayType         ArrayType = C.napi_int8_array
	UInt8ArrayType        ArrayType = C.napi_uint8_array
	UInt8ClampedArrayType ArrayType = C.napi_uint8_clamped_array
	Int16ArrayType        ArrayType = C.napi_int16_array
	UInt16ArrayType       ArrayType = C.napi_uint16_array
	Int32ArrayType        ArrayType = C.napi_int32_array
	UInt32ArrayType       ArrayType = C.napi_uint32_array
	Float32ArrayType      ArrayType = C.napi_float32_array
	Float64ArrayType      ArrayType = C.napi_float64_array
	BigInt64ArrayType     ArrayType = C.napi_bigint64_array
	BigUInt64ArrayType    ArrayType = C.napi_biguint64_array
)

// Status ...
type Status int

const (
	Ok                    Status = C.napi_ok
	InvalidArg            Status = C.napi_invalid_arg
	ObjectExpected        Status = C.napi_object_expected
	StringExpected        Status = C.napi_string_expected
	NameExpected          Status = C.napi_name_expected
	FunctionExpected      Status = C.napi_function_expected
	NumberExpected        Status = C.napi_number_expected
	BooleanExpected       Status = C.napi_boolean_expected
	ArrayExpected         Status = C.napi_array_expected
	GenericFailure        Status = C.napi_generic_failure
	PendingException      Status = C.napi_pending_exception
	Cancelled             Status = C.napi_cancelled
	EscapeCalledTwice     Status = C.napi_escape_called_twice
	HandleScopeMismatch   Status = C.napi_handle_scope_mismatch
	CallbackScopeMismatch Status = C.napi_callback_scope_mismatch
	QueueFull             Status = C.napi_queue_full
	Closing               Status = C.napi_closing
	BigIntExpected        Status = C.napi_bigint_expected
	DateExpected          Status = C.napi_date_expected
)

type callback = C.napi_callback_info

// Callback ...
type Callback struct{}

type finalize = C.napi_finalize

// Finalize ...
type Finalize struct{}

type propertyDescriptor = C.napi_property_descriptor

// PropertyDescriptor ...
type PropertyDescriptor struct{}

type extendedErrorInfo = *C.napi_extended_error_info

// ExtendedErrorInfo ...
type ExtendedErrorInfo struct{}

type callbackScope = C.napi_callback_scope

// CallbackScope ...
type CallbackScope struct{}

type asyncContext = C.napi_async_context

// AsyncContext ...
type AsyncContext struct{}

type asyncWork = C.napi_async_work

// AsyncWork ...
type AsyncWork struct{}

type threadSafeFunction = C.napi_threadsafe_function

// ThereadSafeFunction ...
type ThreadSafeFunction struct{}

// ThreadSafeReleaseMode ...
type ThreadSafeReleaseMode int

const (
	TSFnRelease ThreadSafeReleaseMode = C.napi_tsfn_release
	TSFnAbort   ThreadSafeReleaseMode = C.napi_tsfn_abort
)

// ThreadSafeCallMode ...
type ThreadSafeCallMode int

const (
	TSFnNonBlocking ThreadSafeCallMode = C.napi_tsfn_nonblocking
	TSFnBlocking    ThreadSafeCallMode = C.napi_tsfn_blocking
)

type asyncExecuteCallback = C.napi_async_execute_callback

// AsyncExecuteCallback ...
type AsyncExecuteCallback struct{}

type asyncCompleteCallback = C.napi_async_complete_callback

// AsyncCompleteCallback ...
type AsyncCompleteCallback struct{}

type threadSafeFunctionCallJS = C.napi_threadsafe_function_call_js

// ThreadSafeFunctionCallJS ...
type ThreadSafeFunctionCallJS struct{}

type nodeVersion = *C.napi_node_version

// NodeVersion ...
type NodeVersion struct{}

type uvLoop = *C.struct_uv_loop_s

// UVLoop ...
type UVLoop struct{}

// Hello1
func Hello1() {
	fmt.Println(C.napi_writable)
}

// Hello function ...
func Hello() string {
	return "Hello, world."
}
