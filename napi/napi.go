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
	"bytes"
	"fmt"
	"time"
	"unsafe"
)

// ccstring function transforms a Go string into a C string (array of characters)
// and returns the pointer to the first element.
func cstring(s string) unsafe.Pointer {
	p := make([]byte, len(s)+1)
	copy(p, s)
	return unsafe.Pointer(&p[0])
}

// Aliases for JavaScript types
// Basic N-API Data Types
// N-API exposes the following fundamental datatypes as abstractions that are
// consumed by the various APIs. These APIs should be treated as opaque,
// introspectable only with other N-API calls.

// Env is used to represent a context that the underlying N-API
// implementation can use to persist VM-specific state. This structure is passed
// to native functions when they're invoked, and it must be passed back when
// making N-API calls. Specifically, the same napi_env that was passed in when
// the initial native function was called must be passed to any subsequent nested
// N-API calls. Caching the napi_env for the purpose of general reuse is not
// allowed.
type Env = C.napi_env

// Value is an opaque pointer that is used to represent a JavaScript value.
type Value = C.napi_value

// Ref is an abstraction to use to reference a NapiValue. This allows for
// users to manage the lifetimes of JavaScript values, including defining their
// minimum lifetimes explicitly.
type Ref = C.napi_ref

// HandleScope is an abstraction used to control and modify the lifetime of
// objects created within a particular scope. In general, N-API values are
// created within the context of a handle scope. When a native method is called
// from JavaScript, a default handle scope will exist. If the user does not
// explicitly create a new handle scope, N-API values will be created in the
// default handle scope. For any invocations of code outside the execution of a
// native method (for instance, during a libuv callback invocation), the module
// is required to create a scope before invoking any functions that can result
// in the creation of JavaScript values.
// Handle scopes are created using NapiOnpenHandleScope and are destroyed using
// NapiCloseHandleScope. Closing the scope can indicate to the GC that all
// NapiValues created during the lifetime of the handle scope are no longer
// referenced from the current stack frame.
type HandleScope = C.napi_handle_scope

// EscapableHandleScope represents a special type of handle scope to return
// values created within a particular handle scope to a parent scope.
type EscapableHandleScope = C.napi_escapable_handle_scope

// CallbackInfo is an opaque datatype that is passed to a callback function.
// It can be used for getting additional information about the context in which
// the callback was invoked.
type CallbackInfo = C.napi_callback_info

// Deferred known sa "deferred" object is created and returned alongside a
// Promise. The deferred object is bound to the created Promise and is the only
// means to resolve or reject the Promise. The  deferred object will be
// automatically freed on rejection or on resolving the Promise.
type Deferred = C.napi_deferred

// This is a struct used as container for types of property atrributes.
// NapiPropertyAttributes represents the flags used to control the behavior of
// properties set on a JavaScript object.
// Other than napi_static they correspond to the attributes listed in
// Section 6.1.7.1 of the ECMAScript Language Specification.
// Currently they can be one or more of the following bitflags:
// napi_default - Used to indicate that no explicit attributes are set on the
// given property. By default, a property is read only, not enumerable and not
// configurable.
// napi_writable - Used to indicate that a given property is writable.
// napi_enumerable - Used to indicate that a given property is enumerable.
// napi_configurable - Used to indicate that a given property is configurable,
// as defined in Section 6.1.7.1 of the ECMAScript Language Specification.
// napi_static - Used to indicate that the property will be defined as a static
// property on a class as opposed to an instance property, which is the default.
// This is used only by NapiDefineClass. It is ignored by NapiDefineProperties.
// type PropertyAttributes = C.napi_property_attributes
type propertyAttributes struct {
	Default      int
	Writable     int
	Enumerable   int
	Configurable int
	// Used with napi_define_class to distinguish static properties
	// from instance properties. Ignored by napi_define_properties.
	Static int
}

// PropertyAttributes contains the flags to control the  behavior of properties
// set on a JavaScript object. They can be one or more of the following bitflags:
// - NapiDefault - Used to indicate that no explicit attributes are set on the
// given property. By default, a property is read only, not enumerable and not
// configurable.
// - NapiWritable - Used to indicate that a given property is writable.
// - NapiEnumerable - Used to indicate that a given property is enumerable.
// - NapiConfigurable -  Used to indicate that a given property is configurable,
// as defined in Section 6.1.7.1 of the ECMAScript Language Specification.
// - NapiStatic - Used to indicate that the property will be defined as a static
// property on a class as opposed to an instance property, which is the default.
// This is used only by NapiDefineClass. It is ignored by NapiDfineProperties.
var PropertyAttributes = &propertyAttributes{
	Default:      C.napi_default,
	Writable:     C.napi_writable,
	Enumerable:   C.napi_enumerable,
	Configurable: C.napi_configurable,
	Static:       C.napi_static,
}

// ValueType describes the type of NapiValue. This generally corresponds to
// the types described in Section 6.1 of the ECMAScript Language Specification.
// In addition to types in that section, NapiValueType can also represent
// Functions and Objects with external data.
// A JavaScript value of type napi_external appears in JavaScript as a plain
// object such that no properties can be set on it, and no prototype.
// Currently the following types are supported:
//  napi_undefined,
//  napi_null,
//  napi_boolean,
//  napi_number,
//  napi_string,
//  napi_symbol,
//  napi_object,
//  napi_function,
//  napi_external,
//  napi_bigint,
type ValueType = C.napi_valuetype
type valueTypes struct {
	// ES6 types (corresponds to typeof)
	Undefined int
	Null      int
	Boolean   int
	Number    int
	String    int
	Symbol    int
	Object    int
	Function  int
	External  int
	Bigint    int
}

// ValueTypes contains the type of a NapiValue. This generally corresponds to the
// types described in Section 6.1 of the ECMAScript Language Specification. In
// addition to types in that section, ValueType can also represent Functions and
// Objects with external data. A JavaScript value of type NapiExternal appears in
// JavaScript as a plain object such that no properties can be set on it, and no
//prototype.
var ValueTypes = &valueTypes{
	Undefined: C.napi_undefined,
	Null:      C.napi_null,
	Boolean:   C.napi_boolean,
	Number:    C.napi_number,
	String:    C.napi_string,
	Symbol:    C.napi_symbol,
	Object:    C.napi_object,
	Function:  C.napi_function,
	External:  C.napi_external,
	Bigint:    C.napi_bigint,
}

// This is a struct used as container for types used in TypedArray.
type typedArrayTypes struct {
	Int8Array         int
	UInt8Array        int
	UInt8ClampedArray int
	Int16Array        int
	UInt16Array       int
	Int32Array        int
	UInt32Array       int
	Float32Array      int
	Float64Array      int
	BigInt64Array     int
	BigUInt64Array    int
}

// TypedArrayTypes contains the underlying binary scalar datatype of the
// TypedArray defined in sectiontion 22.2 of the ECMAScript Language
// Specification.
var TypedArrayTypes = &typedArrayTypes{
	Int8Array:         C.napi_int8_array,
	UInt8Array:        C.napi_uint8_array,
	UInt8ClampedArray: C.napi_uint8_clamped_array,
	Int16Array:        C.napi_int16_array,
	UInt16Array:       C.napi_uint16_array,
	Int32Array:        C.napi_int32_array,
	UInt32Array:       C.napi_uint32_array,
	Float32Array:      C.napi_float32_array,
	Float64Array:      C.napi_float64_array,
	BigInt64Array:     C.napi_bigint64_array,
	BigUInt64Array:    C.napi_biguint64_array,
}

// TypedArrayType represents the underlying binary scalar datatype of the
// TypedArray defined in sectiontion 22.2 of the ECMAScript Language
// Specification.
type TypedArrayType C.napi_typedarray_type

// This is a struct used as container for N-API status.
type statuses struct {
	OK                    int
	InvalidArg            int
	ObjectExpected        int
	StringExpected        int
	NameExpected          int
	FunctionExpected      int
	NumberExpected        int
	BooleanExpected       int
	ArrayExpected         int
	GenericFailure        int
	PendingException      int
	Cancelled             int
	EscapeCalledTwice     int
	HandleScopeMismatch   int
	CallbackScopeMismatch int
	QueueFull             int
	Closing               int
	BigintExpected        int
	DateExpected          int
}

// Statuses contains the status code indicating the success or failure of
// a N-API call. Currently, the following status codes are supported:
//  napi_ok
//  napi_invalid_arg
//  napi_object_expected
//  napi_string_expected
//  napi_name_expected
//  napi_function_expected
//  napi_number_expected
//  napi_boolean_expected
//  napi_array_expected
//  napi_generic_failure
//  napi_pending_exception
//  napi_cancelled
//  napi_escape_called_twice
//  napi_handle_scope_mismatch
//  napi_callback_scope_mismatch
//  napi_queue_full
//  napi_closing
//  napi_bigint_expected
//  napi_date_expected
// If additional information is required upon an API returning a failed status,
// it can be obtained by calling NapiGetLastErrorInfo.
var Statuses = &statuses{
	OK:                    C.napi_ok,
	InvalidArg:            C.napi_invalid_arg,
	ObjectExpected:        C.napi_object_expected,
	StringExpected:        C.napi_string_expected,
	NameExpected:          C.napi_name_expected,
	FunctionExpected:      C.napi_function_expected,
	NumberExpected:        C.napi_number_expected,
	BooleanExpected:       C.napi_boolean_expected,
	ArrayExpected:         C.napi_array_expected,
	GenericFailure:        C.napi_generic_failure,
	PendingException:      C.napi_pending_exception,
	Cancelled:             C.napi_cancelled,
	EscapeCalledTwice:     C.napi_escape_called_twice,
	HandleScopeMismatch:   C.napi_handle_scope_mismatch,
	CallbackScopeMismatch: C.napi_callback_scope_mismatch,
	QueueFull:             C.napi_queue_full,
	Closing:               C.napi_closing,
	BigintExpected:        C.napi_bigint_expected,
	DateExpected:          C.napi_date_expected,
}

// Status represent the status code indicating the success or failure of
// a N-API call. Currently, the following status codes are supported:
//  napi_ok
//  napi_invalid_arg
//  napi_object_expected
//  napi_string_expected
//  napi_name_expected
//  napi_function_expected
//  napi_number_expected
//  napi_boolean_expected
//  napi_array_expected
//  napi_generic_failure
//  napi_pending_exception
//  napi_cancelled
//  napi_escape_called_twice
//  napi_handle_scope_mismatch
//  napi_callback_scope_mismatch
//  napi_queue_full
//  napi_closing
//  napi_bigint_expected
// If additional information is required upon an API returning a failed status,
// it can be obtained by calling NapiGetLastErrorInfo.
type Status = C.napi_status

// Callback represents a function pointer type for user-provided native
// functions which are to be exposed to JavaScript via N-API. Callback functions
// should satisfy the following signature:
// typedef napi_value (*napi_callback)(napi_env, napi_callback_info);
type Callback = C.napi_callback

// Finalize represents a function pointer type for add-on provided functions
// that allow the user to be notified when externally-owned data is ready to be
// cleaned up because the object with which it was associated with, has been
// garbage-collected. The user must provide a function satisfying the following
// signature which would get called upon the object's collection. Currently,
// `napi_finalize` can be used for finding out when objects that have external
// data are collected. Finalize functions hould satisfy the following signature:
// typedef void (*napi_finalize)(napi_env env,
//								 void* finalize_data,
//								 void* finalize_hint);
type Finalize = C.napi_finalize

// PropertyDescriptor is a data structure that used to define the properties
// of a JavaScript object.
type PropertyDescriptor = C.napi_property_descriptor

// ExtendedErrorInfo contains additional information about a failed status
// happened on an N-API call.
// The NapiStatus return value provides a VM-independent representation of the
// error which occurred. In some cases it is useful to be able to get more
// detailed information, including a string representing the error as well as
// VM (engine)-specific information.
// error_message: UTF8-encoded string containing a VM-neutral description of the
// error.
// engine_reserved: Reserved for VM-specific error details. This is currently
// not implemented for any VM.
// engine_error_code: VM-specific error code. This is currently not implemented
// for any VM.
// error_code: The N-API status code that originated with the last error.
// Do not rely on the content or format of any of the extended information as it
// is not subject to SemVer and may change at any time. It is intended only for
// logging purposes.
type ExtendedErrorInfo = *C.napi_extended_error_info

// Aliases for types strickly connected with the runtime

// CallbackScope represents
type CallbackScope = C.napi_callback_scope

// AsyncContext represents the context for the async operation that is
// invoking a callback. This should normally be a value previously obtained from
// `napi_async_init`. However `NULL` is also allowed, which indicates the current
// async context (if any) is to be used for the callback.
type AsyncContext = C.napi_async_context

// AsyncWork represents the handle for the newly created asynchronous work
// and it is used to execute logic asynchronously.
type AsyncWork = C.napi_async_work

// ThreadsafeFunction is an opaque pointer that represents a JavaScript
// function which can be called asynchronously from multiple threads.
type ThreadsafeFunction = C.napi_threadsafe_function

// This is a struct used as container for modes to release a
// NapiThreadSafeFunction.
type tsfnReleaseMode struct {
	NapiTsfnRelease int
	NapiTsfnAbort   int
}

// TsfnReleaseMode contains values to be given to NapiReleaseThreadsafeFunction()
// to indicate whether the thread-safe function is to be closed immediately
// (NapiTsfnAbort) or merely released (NapiTsfnRelease) and thus available for
// subsequent use via NapiAcquireThreadsafeFunction() and
// NapiCallThreadsafeFunction().
var TsfnReleaseMode = &tsfnReleaseMode{
	NapiTsfnRelease: C.napi_tsfn_release,
	NapiTsfnAbort:   C.napi_tsfn_abort,
}

// TheradsafeFunctionReleaseMode represents a value to be given to
// NapiReleaseThreadsafeFunction() to indicate whether the thread-safe function
// is to be closed immediately (NapiTsfnAbort) or merely released
// (NapiTsfnRelease) and thus available for subsequent use via
// NapiAcquireThreadsafeFunction() and NapiCallThreadsafeFunction().
type TheradsafeFunctionReleaseMode = C.napi_threadsafe_function_release_mode

// This is a struct used as container for types used to call a
// NapiThreadSafeFunction.
type tsfnCallMode struct {
	NapiTsfnNonBlocking int
	NapiTsfnBlocking    int
}

// TsfnCallMode contains values to be given to NapiCallThreadsafeFunction() to
// indicate whether the call should block whenever the queue associated with the
// thread-safe function is full.
var TsfnCallMode = &tsfnCallMode{
	NapiTsfnNonBlocking: C.napi_tsfn_nonblocking,
	NapiTsfnBlocking:    C.napi_tsfn_blocking,
}

// ThreadsafeFunctionCallMode contains values used to indicate whether the
// call should block whenever the queue associated with the thread-safe function
// is full.
type ThreadsafeFunctionCallMode = C.napi_threadsafe_function_call_mode

// AsyncExecuteCallback is a function pointer used with functions that
// support asynchronous operations. Callback functions must statisfy the
// following signature:
// typedef void (*napi_async_execute_callback)(napi_env env, void* data);
// Implementations of this type of function should avoid making any N-API calls
// that could result in the execution of JavaScript or interaction with
// JavaScript objects.
type AsyncExecuteCallback = C.napi_async_execute_callback

// AsyncCompleteCallback is a function pointer used with functions that
// support asynchronous operations. Callback functions must statisfy the
// following signature:
// typedef void (*napi_async_complete_callback)(napi_env env,
//												napi_status status,
//												void* data);
type AsyncCompleteCallback = C.napi_async_complete_callback

// ThreadsafeFunctionCallJS is a function pointer used with asynchronous
// thread-safe function calls. The callback will be called on the main thread.
// Its purpose is to use a data item arriving via the queue from one of the
// secondary threads to construct the parameters necessary for a call into
// JavaScript.
// Callback functions must satisfy the following signature:
// typedef void (*napi_threadsafe_function_call_js)(napi_env env,
//													napi_value js_callback,
//													void* context,
//													void* data);
type ThreadsafeFunctionCallJS = C.napi_threadsafe_function_call_js

// NodeVersion is a structure that contains informations about the version
// of Node.js instance.
// Currently, the following fields are exposed:
//  major
//  minor
//  patch
//  release
type NodeVersion = *C.napi_node_version

// UVLoop represents the current libuv event loop for a given environment
type UVLoop = *C.struct_uv_loop_s

// Error Handling
// N-API uses both return values and JavaScript exceptions for error handling.
// The following sections explain the approach for each case.
// All of the N-API functions share the same error handling pattern. The return
// type of all API functions is napi_status.
// The return value will be napi_ok if the request was successful and no uncaught
// JavaScript exception was thrown. If an error occurred AND an exception was
// thrown, the napi_status value for the error will be returned. If an exception
// was thrown, and no error occurred, napi_pending_exception will be returned.

// In cases where a return value other than napi_ok or napi_pending_exception is
// returned, napi_is_exception_pending must be called to check if an exception is
// pending. See the section on exceptions for more details.

// The napi_status return value provides a VM-independent representation of the
// error which occurred. In some cases it is useful to be able to get more
// detailed information, including a string representing the error as well as
// VM (engine)-specific information.

// Any N-API function call may result in a pending JavaScript exception. This is
// obviously the case for any function that may cause the execution of
// JavaScript, but N-API specifies that an exception may be pending on return
// from any of the API functions. If the napi_status returned by a function is
// napi_ok then no exception is pending and no additional action is required. If
// the napi_status returned is anything other than napi_ok or
// napi_pending_exception, in order to try to recover and continue instead of
// simply returning immediately, napi_is_exception_pending must be called in
// order to determine if an exception is pending or not. In many cases when an
// N-API function is called and an exception is already pending, the function
// will return immediately with a napi_status of napi_pending_exception.
// However, this is not the case for all functions. N-API allows a subset of the
// functions to be called to allow for some minimal cleanup before returning to
// JavaScript. In that case, napi_status will reflect the status for the
// function. It will not reflect previous pending exceptions. To avoid confusion,
// check the error status after every function call.

// When an exception is pending one of two approaches can be employed.:
// The first approach is to do any appropriate cleanup and then return so that
// execution will return to JavaScript. As part of the transition back to
// JavaScript the exception will be thrown at the point in the JavaScript code
// where the native method was invoked. The behavior of most N-API calls is
// unspecified while an exception is pending, and many will simply return
// napi_pending_exception, so it is important to do as little as possible and
// then return to JavaScript where the exception can be handled.
// The second approach is to try to handle the exception. There will be cases
// where the native code can catch the exception, take the appropriate action,
// and then continue. This is only recommended in specific cases where it is
// known that the exception can be safely handled.

// The Node.js project is adding error codes to all of the errors generated
// internally. The goal is for applications to use these error codes for all
// error checking. The associated error messages will remain, but will only be
// meant to be used for logging and display with the expectation that the message
// can change without SemVer applying. In order to support this model with N-API,
// both in internal functionality and for module specific functionality
// (as its good practice), the throw_ and create_ functions take an optional code
// parameter which is the string for the code to be added to the error object. If
// the optional parameter is NULL then no code will be associated with the error.

// GetLastErrorInfo function returns the information for the last N-API call
// that was made.
// [in] env: The environment that the API is invoked under.
// This API retrieves a napi_extended_error_info structure with information about
// the last error that occurred.
// The content of the napi_extended_error_info returned is only valid up until an
// n-api function is called on the same env.
// Do not rely on the content or format of any of the extended information as it
// is not subject to SemVer and may change at any time. It is intended only for
// logging purposes.
// The function can be called even if there is a pending JavaScript exception.
func GetLastErrorInfo(env Env) (ExtendedErrorInfo, Status) {
	var res *C.napi_extended_error_info
	var status = C.napi_get_last_error_info(env, &res)
	return ExtendedErrorInfo(res), Status(status)
}

// Throw function throws the JavaScript value provided.
// [in] env: The environment that the API is invoked under.
// [in] error: The JavaScript value to be thrown.
// N-API version: 1
func Throw(env Env, error Value) Status {
	return Status(C.napi_throw(env, error))
}

// ThrowError function throws a JavaScript Error with the text provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional error code to be set on the error.
// [in] msg: C string representing the text to be associated with the error.
// N-API version: 1
func ThrowError(env Env, msg string, code string) Status {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	var ccode = C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return Status(C.napi_throw_error(env, ccode, cmsg))
}

// ThrowTypeError function  throws a JavaScript TypeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional error code to be set on the error.
// [in] msg: C string representing the text to be associated with the error.
// N-API version: 1
func ThrowTypeError(env Env, msg string, code string) Status {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	var ccode = C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return Status(C.napi_throw_type_error(env, ccode, cmsg))
}

// ThrowRangError function throws a JavaScript RangeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional error code to be set on the error.
// [in] msg: C string representing the text to be associated with the error.
// N-API version: 1
func ThrowRangError(env Env, msg string, code string) Status {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	var ccode = C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return Status(C.napi_throw_range_error(env, ccode, cmsg))
}

// IsError function queries a napi_value to check if it represents an error
// object.
// [in] env: The environment that the API is invoked under.
// [in] value: The napi_value to be checked.
// Boolean value that is set to true if napi_value represents an error, false
// otherwise.
// N-API version: 1
func IsError(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_is_error(env, value, &res)
	return bool(res), Status(status)
}

// CreateError function creates a JavaScript Error with the text provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional napi_value with the string for the error code to be
// associated with the error.
// [in] msg: napi_value that references a JavaScript String to be used as the
// message for the Error.
// N-API version: 1
func CreateError(env Env, msg Value, code Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_error(env, code, msg, &res)
	return Value(res), Status(status)
}

// CreateTypeError function creates a JavaScript TypeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional napi_value with the string for the error code to be
// associated with the error.
// [in] msg: napi_value that references a JavaScript String to be used as the
// message for the Error.
// N-API version: 1
func CreateTypeError(env Env, code Value, msg Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_type_error(env, code, msg, &res)
	return Value(res), Status(status)
}

// CreateRangeError function creates a JavaScript RangeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional napi_value with the string for the error code to be
// associated with the error.
// [in] msg: napi_value that references a JavaScript String to be used as the
// message for the Error.
// N-API version: 1
func CreateRangeError(env Env, code Value, msg Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_range_error(env, code, msg, &res)
	return Value(res), Status(status)
}

// GetAndClearLastException function returns true if an exception is pending.
// This function can be called even if there is a pending JavaScript exception.
// [in] env: The environment that the API is invoked under.
// The function returns the exception if one is pending, NULL otherwise.
// N-API version: 1
func GetAndClearLastException(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_and_clear_last_exception(env, &res)
	return Value(res), Status(status)
}

// IsExceptionPending function ...
// [in] env: The environment that the API is invoked under.
// Boolean value that is set to true if an exception is pending.
// N-API version: 1
func IsExceptionPending(env Env) (bool, Status) {
	var res C.bool
	var status = C.napi_is_exception_pending(env, &res)
	return bool(res), Status(status)
}

// FatalException function triggers an 'uncaughtException' in JavaScript.
// Useful if an async callback throws an exception with no way to recover.
// [in] env: The environment that the API is invoked under.
// [in] err: The error that is passed to 'uncaughtException'.
// N-API version: 3
func FatalException(env Env, error Value) Status {
	return Status(C.napi_fatal_exception(env, error))
}

// FatalError function thrown a fatal error o immediately terminate the
// process.
// [in] location: Optional location at which the error occurred.
// [in] location_len: The length of the location in bytes, or NAPI_AUTO_LENGTH
// if it is null-terminated.
// [in] message: The message associated with the error.
// [in] message_len: The length of the message in bytes, or NAPI_AUTO_LENGTH if
// it is null-terminated.
// This function can be called even if there is a pending JavaScript exception.
// The function call does not return, the process will be terminated.
// N-API version: 1
func FatalError(location string, msg string) {
	clocation := C.CString(location)
	defer C.free(unsafe.Pointer(clocation))
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	C.napi_fatal_error(clocation, C.NAPI_AUTO_LENGTH, cmsg, C.NAPI_AUTO_LENGTH)
	return
}

// Object Lifetime management
// As N-API calls are made, handles to objects in the heap for the underlying VM
// may be returned as napi_values. These handles must hold the objects 'live'
// until they are no longer required by the native code, otherwise the objects
// could be collected before the native code was finished using them. As object
// handles are returned they are associated with a 'scope'. The lifespan for the
// default scope is tied to the lifespan of the native method call. The result is
// that, by default, handles remain valid and the objects associated with these
// handles will be held live for the lifespan of the native method call. In many
// cases, however, it is necessary that the handles remain valid for either a
// shorter or longer lifespan than that of the native method.
// N-API only supports a single nested hierarchy of scopes. There is only one
// active scope at any time, and all new handles will be associated with that
// scope while it is active. Scopes must be closed in the reverse order from
// which they are opened. In addition, all scopes created within a native method
// must be closed before returning from that method.
// When nesting scopes, there are cases where a handle from an inner scope needs
// to live beyond the lifespan of that scope. N-API supports an 'escapable scope'
// in order to support this case. An escapable scope allows one handle to be
// 'promoted' so that it 'escapes' the current scope and the lifespan of the
// handle changes from the current scope to that of the outer scope.

// OnpenHandleScope function opens a new scope.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func OnpenHandleScope(env Env) (HandleScope, Status) {
	var res C.napi_handle_scope
	var status = C.napi_open_handle_scope(env, &res)
	return HandleScope(res), Status(status)
}

// CloseHandleScope function closes the scope passed in. Scopes must be
// closed in the reverse order from which they were created.
// [in] env: The environment that the API is invoked under.
// [in] scope: napi_value representing the scope to be closed.
// This function can be called even if there is a pending JavaScript exception.
// N-API version: 1
func CloseHandleScope(env Env, scope HandleScope) Status {
	return Status(C.napi_close_handle_scope(env, scope))
}

// OnpenEscapableHandleScope function opens a new scope from which one object
// can be promoted to the outer scope.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func OnpenEscapableHandleScope(env Env) (EscapableHandleScope, Status) {
	var res C.napi_escapable_handle_scope
	var status = C.napi_open_escapable_handle_scope(env, &res)
	return EscapableHandleScope(res), Status(status)
}

// CloseEscapableHandleScope function closes the scope passed in. Scopes must
// be closed in the reverse order from which they were created.
// [in] env: The environment that the API is invoked under.
// [in] scope: napi_value representing the scope to be closed.
// This function can be called even if there is a pending JavaScript exception.
// N-API version: 1
func CloseEscapableHandleScope(env Env, scope EscapableHandleScope) Status {
	return Status(C.napi_close_escapable_handle_scope(env, scope))
}

// EscapeHandle function promotes the handle to the JavaScript object so that
// it is valid for the lifetime of the outer scope. It can only be called once
// per scope. If it is called more than once an error will be returned.
// [in] env: The environment that the API is invoked under.
// [in] scope: napi_value representing the current scope.
// [in] escapee: napi_value representing the JavaScript Object to be escaped.
// This function can be called even if there is a pending JavaScript exception.
// N-API version: 1
func EscapeHandle(env Env, scope EscapableHandleScope, escapee Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_escape_handle(env, scope, escapee, &res)
	return Value(res), Status(status)
}

// References to objects with a lifespan longer than that of the native method
// In some cases an addon will need to be able to create and reference objects
// with a lifespan longer than that of a single native method invocation.
// For example, to create a constructor and later use that constructor in a
// request to creates instances, it must be possible to reference the constructor
// object across many different instance creation requests. This would not be
// possible with a normal handle returned as a NapiValue as described in the
// earlier section. The lifespan of a normal handle is managed by scopes and all
// scopes must be closed before the end of a native method.

// N-API provides methods to create persistent references to an object. Each
// persistent reference has an associated count with a value of 0 or higher. The
// count determines if the reference will keep the corresponding object live.
// References with a count of 0 do not prevent the object from being collected
// and are often called 'weak' references. Any count greater than 0 will prevent
// the object from being collected.

// References must be deleted once they are no longer required by the addon.
// When a reference is deleted it will no longer prevent the corresponding object
// from being collected. Failure to delete a persistent reference will result in
// a 'memory leak' with both the native memory for the persistent reference and
// the corresponding object on the heap being retained forever.

// There can be multiple persistent references created which refer to the same
// object, each of which will either keep the object live or not based on its
// individual count.

// CreateReference function creates a new reference with the specified
// reference count to the Object passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing the Object to which we want a reference.
// [in] initial_refcount: Initial reference count for the new reference.
// N-API version: 1
func CreateReference(env Env, value Value, refCount uint) (Ref, Status) {
	var res C.napi_ref
	var status = C.napi_create_reference(env, value, C.uint(refCount), &res)
	return Ref(res), Status(status)
}

// DeleteReference function deletes the reference passed in.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref to be deleted.
// This API can be called even if there is a pending JavaScript exception.
// N-API version: 1
func DeleteReference(env Env, ref Ref) Status {
	var status = C.napi_delete_reference(env, ref)
	return Status(status)
}

// ReferenceRef function  increments the reference count for the reference
// passed in and returns the resulting reference count.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref for which the reference count will be incremented.
// N-API version: 1
func ReferenceRef(env Env, ref Ref) (uint, Status) {
	var res C.uint
	var status = C.napi_reference_ref(env, ref, &res)
	return uint(res), Status(status)
}

// ReferenceUnref function ecrements the reference count for the reference
// passed in and returns the resulting reference count.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref for which the reference count will be decremented.
// N-API version: 1
func ReferenceUnref(env Env, ref Ref) (uint, Status) {
	var res C.uint
	var status = C.napi_reference_unref(env, ref, &res)
	return uint(res), Status(status)
}

// GetReferenceValue function returns the NapiValue representing the
// JavaScript Object associated with the NapiRef. Otherwise, result will be
// NULL.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref for which we requesting the corresponding Object.
// N-API version: 1
func GetReferenceValue(env Env, ref Ref) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_reference_value(env, ref, &res)
	return Value(res), Status(status)
}

// AddEnvCleanupHook function ...
func AddEnvCleanupHook(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_ok
	return Value(res), Status(status)
}

// RemoveCleaupHook function ...
func RemoveCleaupHook(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_ok
	return Value(res), Status(status)
}

// CreateArray function returns an N-API value corresponding to a JavaScript
// Array type. JavaScript arrays are described in Section 22.1 of the ECMAScript
// Language Specification.
// [in] env: The environment that the N-API call is invoked under.
// N-API version: 1
func CreateArray(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_array(env, &res)
	return Value(res), Status(status)
}

// CreateArrayWithLength function returns an N-API value corresponding to a
// JavaScript Array type. The Array's length property is set to the passed-in
// length parameter. However, the underlying buffer is not guaranteed to be
// pre-allocated by the VM when the array is created - that behavior is left to
// the underlying VM implementation.
// avaScript arrays are described in Section 22.1 of the ECMAScript Language
// Specification.
// [in] env: The environment that the API is invoked under.
// [in] length: The initial length of the Array.
// N-API version: 1
func CreateArrayWithLength(env Env, length uint) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_array_with_length(env, C.size_t(length), &res)
	return Value(res), Status(status)
}

// CreateArrayBuffer function returns N-API value corresponding to a
// JavaScript `ArrayBuffer`. ArrayBuffer is a data stucture used to represent
// fixed-length binary data buffers. They are normally used as backing-buffer for
// `TypedArray` objects. The ArrayBuffer allocated will have an underlying byte
// buffer whose size is determined by the length parameter that's passed in. The
// underlying buffer is optionally returned back to the caller in case the caller
// wants to directly manipulate the buffer.
// This buffer can only be written to directly from native code.
// To write to this buffer from JavaScript, a typed array or DataView object
// would need to be created.
// JavaScript ArrayBuffer objects are described in Section 24.1 of the ECMAScript
// Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] length: The length in bytes of the array buffer to create.
// [out] data: Pointer to the underlying byte buffer of the ArrayBuffer.
// [out] result: A napi_value representing a JavaScript ArrayBuffer.
// N-API version: 1
func CreateArrayBuffer(env Env, length uint) (Value, unsafe.Pointer, Status) {
	var res C.napi_value
	var data unsafe.Pointer
	var status = C.napi_create_arraybuffer(env, C.size_t(length), &data, &res)
	return Value(res), data, Status(status)
}

// CreateBuffer function returns N-API value that allocates a node::Buffer
// object. While this is still a fully-supported data structure, in most cases
// musing a TypedArray will suffice.
// [in] env: The environment that the API is invoked under.
// [in] size: Size in bytes of the underlying buffer.
// [out] data: Raw pointer to the underlying buffer.
// [out] result: A napi_value representing a node::Buffer.
// N-API version: 1
func CreateBuffer(env Env, length uint) (Value, unsafe.Pointer, Status) {
	var res C.napi_value
	var data unsafe.Pointer
	var status = C.napi_create_buffer(env, C.size_t(length), &data, &res)
	return Value(res), data, Status(status)
}

// CreateBufferCopy function  allocates a node::Buffer object and initializes
// it with data copied from the passed-in buffer. While this is still a
// fully-supported data structure, in most cases using a TypedArray will suffice.
// [in] env: The environment that the API is invoked under.
// [in] length: Size in bytes of the input buffer (should be the same as the size
// of the new buffer).
// [in] data: Raw pointer to the underlying buffer to copy from.
// [out] result_data: Pointer to the new Buffer's underlying data buffer.
// [out] result: A napi_value representing a node::Buffer.
// N-API version: 1
func CreateBufferCopy(env Env, length uint, raw unsafe.Pointer) (Value, unsafe.Pointer, Status) {
	var res C.napi_value
	var data unsafe.Pointer
	var status = C.napi_create_buffer_copy(env, C.size_t(length), raw, &data, &res)
	return Value(res), data, Status(status)
}

// CreateExternal function allocates a JavaScript value with external data
// attached to it. This is used to pass external data through JavaScript code, so
// it can be retrieved later by native code. The API allows the caller to pass in
// a finalize callback, in case the underlying native resource needs to be
// cleaned up when the external JavaScript value gets collected.
// [in] env: The environment that the API is invoked under.
// [in] data: Raw pointer to the external data.
// [in] finalize_cb: Optional callback to call when the external value is being
// collected.
// [in] finalize_hint: Optional hint to pass to the finalize callback during
// collection.
// [out] result: A napi_value representing an external value.
// The created value is not an object, and therefore does not support additional
// properties. It is considered a distinct value type `napi_external`.
// N-API version: 1
func CreateExternal(env Env, raw unsafe.Pointer) (Value, Status) {
	var res C.napi_value
	// Remember to handle napi_finalize finalize_cb and void* finalize_hint
	var status = C.napi_create_external(env, raw, nil, nil, &res)
	return Value(res), Status(status)
}

// CreateExternalArrayBuffer function returns an N-API value corresponding to
// a JavaScript ArrayBuffer. The underlying byte buffer of the ArrayBuffer is
// externally allocated and managed. The caller must ensure that the byte buffer
// remains valid until the finalize callback is called.
// [in] env: The environment that the API is invoked under.
// [in] external_data: Pointer to the underlying byte buffer of the ArrayBuffer.
// [in] byte_length: The length in bytes of the underlying buffer.
// [in] finalize_cb: Optional callback to call when the ArrayBuffer is being
// collected.
// [in] finalize_hint: Optional hint to pass to the finalize callback during
// collection.
// [out] result: A napi_value representing a JavaScript ArrayBuffer.
// JavaScript ArrayBuffers are described in Section 24.1 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateExternalArrayBuffer(env Env, length uint, raw unsafe.Pointer) (Value, Status) {
	var res C.napi_value
	// Remember to handle napi_finalize finalize_cb and void* finalize_hint
	var status = C.napi_create_external_arraybuffer(env, raw, C.size_t(length), nil, nil, &res)
	return Value(res), Status(status)
}

// CreateExternalBuffer function allocates a node::Buffer object and
// initializes it with data backed by the passed in buffer. While this is still a
// fully-supported data structure, in most cases using a TypedArray will suffice.
// [in] env: The environment that the API is invoked under.
// [in] length: Size in bytes of the input buffer (should be the same as the size
// of the new buffer).
// [in] data: Raw pointer to the underlying buffer to copy from.
// [in] finalize_cb: Optional callback to call when the ArrayBuffer is being
// collected.
// [in] finalize_hint: Optional hint to pass to the finalize callback during
// collection.
// [out] result: A napi_value representing a node::Buffer.
// Remember that fsor Node.js >=4 Buffers are Uint8Array.
//  N-API version: 1
func CreateExternalBuffer(env Env, length uint, raw unsafe.Pointer) (Value, Status) {
	var res C.napi_value
	// Remember to handle napi_finalize finalize_cb and void* finalize_hint
	var status = C.napi_create_external_buffer(env, C.size_t(length), raw, nil, nil, &res)
	return Value(res), Status(status)
}

// CreateObject function allocates a default JavaScript Object. It is the
// equivalent of doing new Object() in JavaScript.
// The JavaScript Object type is described in Section 6.1.7 of the ECMAScript
// Language Specification.
// [in] env: The environment that the API is invoked under.
// [out] result: A napi_value representing a JavaScript Object.
// Returns napi_ok if the API succeeded.
// N-API version: 1
func CreateObject(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_object(env, &res)
	return Value(res), Status(status)
}

// CreateSymbol function creates a JavaScript Symbol object from a
// UTF8-encoded C string.
// The JavaScript Symbol type is described in Section 19.4 of the ECMAScript
// Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] description: Optional napi_value which refers to a JavaScript String to
// be set as the description for the symbol.
// [out] result: A napi_value representing a JavaScript Symbol.
// Returns napi_ok if the API succeeded.
// N-API version: 1
func CreateSymbol(env Env, value Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_symbol(env, value, &res)
	return Value(res), Status(status)
}

// CreateTypedArray function JavaScript TypedArray object over an existing
// ArrayBuffer.  TypedArray objects provide an array-like view over an underlying
// data buffer where each element has the same underlying binary scalar datatype.
// It's required that:
// (length * size_of_element) + byte_offset should be <= the size in bytes of the
// array passed in. If not, a RangeError exception is raised.
// [in] env: The environment that the API is invoked under.
// [in] type: Scalar datatype of the elements within the TypedArray.
// [in] length: Number of elements in the TypedArray.
// [in] arraybuffer: ArrayBuffer underlying the typed array.
// [in] byte_offset: The byte offset within the ArrayBuffer from which to start
// projecting the TypedArray.
// [out] result: A napi_value representing a JavaScript TypedArray.
// JavaScript TypedArray objects are described in Section 22.2 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateTypedArray(env Env, arrayType TypedArrayType, lenght uint, value Value, offset uint) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_typedarray(env, (C.napi_typedarray_type)(arrayType), C.size_t(lenght), (C.napi_value)(value), C.size_t(offset), &res)
	return Value(res), Status(status)
}

// CreateDataview function creates a JavaScript DataView object over an
// existing ArrayBuffer. DataView objects provide an array-like view over an
// underlying data buffer, but one which allows items of different size and type
// in the ArrayBuffer.
// [in] env: The environment that the API is invoked under.
// [in] length: Number of elements in the DataView.
// [in] arraybuffer: ArrayBuffer underlying the DataView.
// [in] byte_offset: The byte offset within the ArrayBuffer from which to start
// projecting the DataView.
// [out] result: A napi_value representing a JavaScript DataView.
// It is required that byte_length + byte_offset is less than or equal to the
// size in bytes of the array passed in. If not, a RangeError exception is
// raised.
// JavaScript DataView objects are described in Section 24.3 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateDataview(env Env, length uint, offset uint, value Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_dataview(env, C.size_t(length), (C.napi_value)(value), C.size_t(offset), &res)
	return Value(res), Status(status)
}

// CreateInt32 function creates JavaScript Number from the C int32_t type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateInt32(env Env, value int32) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_int32(env, C.int(value), &res)
	return Value(res), Status(status)
}

// CreateUInt32 function creates JavaScript Number from the C uint32_t type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateUInt32(env Env, value uint32) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_uint32(env, C.uint(value), &res)
	return Value(res), Status(status)
}

// CreateInt64 function creates JavaScript Number from the C int64_t type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateInt64(env Env, value int64) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_int64(env, C.int64_t(value), &res)
	return Value(res), Status(status)
}

// CreateDouble function creates JavaScript Number from the C double type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateDouble(env Env, value float64) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_double(env, C.double(value), &res)
	return Value(res), Status(status)
}

// CreateBigintInt64 function creates JavaScript BigInt from the C int64_t
// type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// N-API version: -
func CreateBigintInt64(env Env, value int64) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_bigint_int64(env, C.int64_t(value), &res)
	return Value(res), Status(status)
}

// CreateBigintUInt64 function creates JavaScript BigInt from the C uint64_t
// type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// N-API version: -
func CreateBigintUInt64(env Env, value uint64) (Value, Status) {
	var res C.napi_value
	var status = C.napi_create_bigint_uint64(env, C.uint64_t(value), &res)
	return Value(res), Status(status)
}

// CreateBigintWords function converts an array of unsigned 64-bit words into
// a single BigInt value.
// [in] env: The environment that the API is invoked under.
// [in] sign_bit: Determines if the resulting BigInt will be positive or
// negative.
// [in] word_count: The length of the words array.
// [in] words: An array of uint64_t little-endian 64-bit words.
// [out] result: A napi_value representing a JavaScript BigInt.
// Returns napi_ok if the API succeeded.
// N-API version: -
func CreateBigintWords(env Env, sign int, words []uint64) (Value, Status) {
	var res C.napi_value
	var raw = (unsafe.Pointer(&words[0]))
	defer C.free(raw)
	var status = C.napi_create_bigint_words(env, C.int(sign), C.size_t(len(words)), (*C.uint64_t)(raw), &res)
	return Value(res), Status(status)
}

// CreateStringLatin1 function creates a JavaScript String object from an
// ISO-8859-1-encoded C string. The native string is copied.
// [in] env: The environment that the API is invoked under.
// [in] str: Character buffer representing an ISO-8859-1-encoded string.
// The JavaScript String type is described in Section 6.1.4 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateStringLatin1(env Env, str string) (Value, Status) {
	var res C.napi_value
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	var status = C.napi_create_string_latin1(env, cstr, C.NAPI_AUTO_LENGTH, &res)
	return Value(res), Status(status)
}

// CreateStringUtf16 function creates a JavaScript String object from a
// UTF16-LE-encoded C string. The native string is copied.
// [in] env: The environment that the API is invoked under.
// [in] str: Character buffer representing a UTF16-LE-encoded string.
// The JavaScript String type is described in Section 6.1.4 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateStringUtf16(env Env, str string) (Value, Status) {
	var res C.napi_value
	cstr := (*C.ushort)(cstring(str))
	defer C.free(unsafe.Pointer(cstr))
	var status = C.napi_create_string_utf16(env, cstr, C.NAPI_AUTO_LENGTH, &res)
	return Value(res), Status(status)
}

// CreateStringUtf8 function creates a JavaScript String object from a
// UTF8-encoded C string. The native string is copied.
// [in] env: The environment that the API is invoked under.
// [in] str: Character buffer representing a UTF8-encoded string.
// The JavaScript String type is described in Section 6.1.4 of the ECMAScript
// Language Specification.
// N-API version: 1
func CreateStringUtf8(env Env, str string) (Value, Status) {
	var res C.napi_value
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	var status = C.napi_create_string_utf8(env, cstr, C.NAPI_AUTO_LENGTH, &res)
	return Value(res), Status(status)
}

// GetArrayLength function returns the length of an array.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing the JavaScript Array whose length is
// being queried.
// [out] result: uint32 representing length of the array.
// Returns napi_ok if the API succeeded.
// Array length is described in Section 22.1.4.1 of the ECMAScript Language
// Specification.
// N-API version: 1
func GetArrayLength(env Env, value Value) (uint32, Status) {
	var res C.uint32_t
	var status = C.napi_get_array_length(env, value, &res)
	return uint32(res), Status(status)
}

// GetArrayBufferInfo function the underlying data buffer of a node::Buffer
// and it's length.
// [in] env: The environment that the API is invoked under.
//  N-API version: 1
func GetArrayBufferInfo(env Env, value Value) (unsafe.Pointer, uint, Status) {
	var data unsafe.Pointer
	var length C.size_t
	var status = C.napi_get_buffer_info(env, value, &data, &length)
	return data, uint(length), Status(status)
}

// GetPrototype function returns a N-API value representing the prototype of
// the given object.
// [in] env: The environment that the API is invoked under.
// [in] object: napi_value representing JavaScript Object whose prototype to
// return. This returns the equivalent of Object.getPrototypeOf (which is not the
// same as the function's prototype property).
// [out] result: napi_value representing prototype of the given object.
// N-API version: 1
func GetPrototype(env Env, object Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_prototype(env, (C.napi_value)(object), &res)
	return Value(res), Status(status)
}

// GetTypedArrayInfo function returns various properties of a typed array.
// [in] env: The environment that the API is invoked under.
// [in] typedarray: napi_value representing the TypedArray whose properties to
// query.
// [out] type: Scalar datatype of the elements within the TypedArray.
// [out] length: The number of elements in the TypedArray.
// [out] data: The data buffer underlying the TypedArray adjusted by the
// byte_offset value so that it points to the first element in the TypedArray.
// [out] arraybuffer: The ArrayBuffer underlying the TypedArray.
// [out] byte_offset: The byte offset within the underlying native array at
// which the first element of the arrays is located. The value for the data
// parameter has already been adjusted so that data points to the first element
// in the array. Therefore, the first byte of the native array would be at
// data - byte_offset.
// Warning: Use caution while using this API since the underlying data buffer is
// managed by the VM.
// N-API version: 1
func GetTypedArrayInfo(env Env, value Value) (Value, TypedArrayType, uint, unsafe.Pointer, uint, Status) {
	var arrayType C.napi_typedarray_type
	var length C.size_t
	var data unsafe.Pointer
	var arraybuffer C.napi_value
	var offset C.size_t
	var status = C.napi_get_typedarray_info(env, value, &arrayType, &length, &data, &arraybuffer, &offset)
	return Value(arraybuffer), TypedArrayType(arrayType), uint(length), data, uint(offset), Status(status)
}

// GetDataviewInfo function eturns various properties of a DataView.
// [in] env: The environment that the API is invoked under.
// [in] dataview: napi_value representing the DataView whose properties to query.
// [out] byte_length: Number of bytes in the DataView.
// [out] data: The data buffer underlying the DataView.
// [out] arraybuffer: ArrayBuffer underlying the DataView.
// [out] byte_offset: The byte offset within the data buffer from which to start
// projecting the DataView.
// N-API version: 1
func GetDataviewInfo(env Env, value Value) (Value, uint, uint, Status) {
	var length C.size_t
	var data unsafe.Pointer
	var arraybuffer C.napi_value
	var offset C.size_t
	var status = C.napi_get_dataview_info(env, value, &length, &data, &arraybuffer, &offset)
	return Value(arraybuffer), uint(length), uint(offset), Status(status)
}

// GetValueBool function returns the C boolean primitive equivalent of the
// given JavaScript Boolean.
// [in] env: The environment that the API is invoked under.
// [in] value: NapiValue representing JavaScript Boolean.
// Returns napi_ok if the API succeeded. If a non-boolean NapiValue is passed
// in it returns napi_boolean_expected.
// N-API version: 1
func GetValueBool(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_get_value_bool(env, value, &res)
	return bool(res), Status(status)
}

// GetValueDouble function returns the C double primitive equivalent of the
// given JavaScript Number.
// [in] env: The environment that the API is invoked under.
// [in] value: NapiValue representing JavaScript Number.
// Returns napi_ok if the API succeeded. If a non-number NapiValue is passed in
// it returns napi_number_expected.
// N-API version: 1
func GetValueDouble(env Env, value Value) (float64, Status) {
	var res C.double
	var status = C.napi_get_value_double(env, value, &res)
	return float64(res), Status(status)
}

// GetValueBigintInt64 function returns the C int64_t primitive equivalent of
// the given JavaScript BigInt. If needed it will truncate the value, setting
// lossless to false.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript BigInt.
// Returns napi_ok if the API succeeded. If a non-BigInt is passed in it returns
// napi_bigint_expected.
// N-API version: -
func GetValueBigintInt64(env Env, value Value) (int64, bool, Status) {
	var res C.int64_t
	var lossless C.bool
	var status = C.napi_get_value_bigint_int64(env, value, &res, &lossless)
	return int64(res), bool(lossless), Status(status)
}

// GetValueBigintUInt64 function returns the C uint64_t primitive equivalent
// of the given JavaScript BigInt. If needed it will truncate the value, setting
// lossless to false.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript BigInt.
// Returns napi_ok if the API succeeded. If a non-BigInt is passed in it returns
// napi_bigint_expected.
// N-API version: -
func GetValueBigintUInt64(env Env, value Value) (uint64, bool, Status) {
	var res C.uint64_t
	var lossless C.bool
	var status = C.napi_get_value_bigint_uint64(env, value, &res, &lossless)
	return uint64(res), bool(lossless), Status(status)
}

// GetValueBigintWords function returns a single `BigInt` value into a sign
// bit, 64-bit little endian array, and the number of elements backed in the
// array. The sign_bit and words may be both set to NULL, in order to get only
// word_count.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript BigInt.
// [out] sign_bit: Integer representing if the JavaScript BigInt is positive or
// negative.
// [in/out] word_count: Must be initialized to the length of the words array.
// Upon return, it will be set to the actual number of words that would be
// needed to store this BigInt.
// [out] words: Pointer to a pre-allocated 64-bit word array.
// N-API version: -
func GetValueBigintWords(env Env, value Value) (unsafe.Pointer, uint, int, Status) {
	var count C.size_t
	var sign C.int
	var words unsafe.Pointer
	var status = C.napi_get_value_bigint_words(env, value, &sign, &count, (*C.uint64_t)(words))
	return words, uint(count), int(sign), Status(status)
}

// GetValueExternal function returns external data pointer that was
// previously passed to NapiCreateExternal.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript external value.
// [out] result: Pointer to the data wrapped by the JavaScript external value.
// If a non-external napi_value is passed in it returns napi_invalid_arg.
// N-API version: 1
func GetValueExternal(env Env, value Value) (unsafe.Pointer, Status) {
	var res unsafe.Pointer
	var status = C.napi_get_value_external(env, value, &res)
	return res, Status(status)
}

// GetValueInt32 function returns the C int32 primitive equivalent of the
// given JavaScript Number.
// If the number exceeds the range of the 32 bit integer, then the result is
// truncated to the equivalent of the bottom 32 bits. This can result in a large
// positive number becoming a negative number if the value is > 2^31 -1.
// Non-finite number values (NaN, +Infinity, or -Infinity) set the result to
// zero.
// N-API version: 1
func GetValueInt32(env Env, value Value) (int32, Status) {
	var res C.int32_t
	var status = C.napi_get_value_int32(env, value, &res)
	return int32(res), Status(status)
}

// GetValueInt64 function returns the C int64 primitive equivalent of the
// given JavaScript Number.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript Number.
// Returns napi_ok if the API succeeded. If a non-number NapiValue is passed in
// it returns napi_number_expected.
// Number values outside the range of Number.MIN_SAFE_INTEGER -(2^53 - 1) - Number.MAX_SAFE_INTEGER (2^53 - 1)
// will lose precision.
// Non-finite number values (NaN, +Infinity, or -Infinity) set the result to
// zero.
// N-API version: 1
func GetValueInt64(env Env, value Value) (int64, Status) {
	var res C.int64_t
	var status = C.napi_get_value_int64(env, value, &res)
	return int64(res), Status(status)
}

// GetValueStringLatin1 function returns the ISO-8859-1-encoded string
// corresponding the value passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript string.
// [in] buf: Buffer to write the ISO-8859-1-encoded string into. If NULL is
// passed in, the length of the string (in bytes) is returned.
// [in] bufsize: Size of the destination buffer. When this value is insufficient,
// the returned string will be truncated.
// [out] result: Number of bytes copied into the buffer, excluding the null
// terminator.
// Returns napi_ok if the API succeeded. If a non-String napi_value is passed in
// it returns napi_string_expected.
// N-API version: 1
func GetValueStringLatin1(env Env, value Value, len uint) (string, Status) {
	var buf (*C.char)
	var res C.size_t
	var status = C.napi_get_value_string_latin1(env, value, buf, C.size_t(len), &res)
	return string(C.GoStringN(buf, C.int(res))), Status(status)
}

// GetValueStringUtf8 function returns the UTF16-encoded string corresponding
// the value passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript string.
// [in] buf: Buffer to write the UTF8-encoded string into. If NULL is passed in,
// the length of the string (in bytes) is returned.
// [in] bufsize: Size of the destination buffer. When this value is insufficient,
// the returned string will be truncated.
// [out] result: Number of bytes copied into the buffer, excluding the null
// terminator.
// Returns napi_ok if the API succeeded. If a non-String napi_value is passed in
// it returns napi_string_expected.
// N-API version: 1
func GetValueStringUtf8(env Env, value Value, len uint) (string, Status) {
	var buf (*C.char)
	var res C.size_t
	var status = C.napi_get_value_string_utf8(env, value, buf, C.size_t(len), &res)
	return string(C.GoStringN(buf, C.int(res))), Status(status)
}

// GetValueStringUtf16 function returns the UTF16-encoded string
// corresponding the value passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript string.
// [in] buf: Buffer to write the UTF16-LE-encoded string into. If NULL is passed
// in, the length of the string (in 2-byte code units) is returned.
// [in] bufsize: Size of the destination buffer. When this value is insufficient,
// the returned string will be truncated.
// [out] result: Number of 2-byte code units copied into the buffer, excluding
// the null terminator.
// Returns napi_ok if the API succeeded. If a non-String napi_value is passed in
// it returns napi_string_expected.
// N-API version: 1
func GetValueStringUtf16(env Env, value Value, len uint) (string, Status) {
	var buf (*C.ushort)
	var res C.size_t
	var status = C.napi_get_value_string_utf16(env, value, buf, C.size_t(len), &res)
	var str = bytes.NewBuffer(C.GoBytes(unsafe.Pointer(buf), C.int(res))).String()
	return str, Status(status)
}

// GetValueUint32 function returns the C primitive equivalent of the
// given NapiValue as a uint32_t
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript Number.
// Returns napi_ok if the API succeeded. If a non-number NapiValue is passed in
// it returns napi_number_expected.
// N-API version: 1
func GetValueUint32(env Env, value Value) (uint32, Status) {
	var res C.uint32_t
	var status = C.napi_get_value_uint32(env, value, &res)
	return uint32(res), Status(status)
}

// GetBoolean function returns the JavaScript singleton object that is used
// to represent the given boolean value.
// [in] env: The environment that the API is invoked under.
// [in] value: The value of the boolean to retrieve.
// N-API version: 1
func GetBoolean(env Env, value bool) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_boolean(env, C.bool(value), &res)
	return Value(res), Status(status)
}

// GetGlobal function returns JavaScript global object.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func GetGlobal(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_global(env, &res)
	return Value(res), Status(status)
}

// GetNull function returns JavaScript null object.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func GetNull(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_null(env, &res)
	return Value(res), Status(status)
}

// GetUndefined function returns JavaScript Undefined value.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func GetUndefined(env Env) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_undefined(env, &res)
	return Value(res), Status(status)
}

// CoerceToBool function implements the abstract operation ToBoolean() as
// defined in Section 7.1.2 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func CoerceToBool(env Env, value Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_coerce_to_bool(env, value, &res)
	return Value(res), Status(status)
}

// CoerceToNumber function implements the abstract operation ToNumber() as
// defined in Section 7.1.3 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func CoerceToNumber(env Env, value Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_coerce_to_number(env, value, &res)
	return Value(res), Status(status)
}

// CoerceToObject function implements the abstract operation ToObject() as
// defined in Section 7.1.13 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func CoerceToObject(env Env, value Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_coerce_to_object(env, value, &res)
	return Value(res), Status(status)
}

// CoerceToString function mplements the abstractoperation ToString() as
// defined in Section 7.1.13 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func CoerceToString(env Env, value Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_coerce_to_string(env, value, &res)
	return Value(res), Status(status)
}

// TypeOf function is similar to invoke the typeof Operator on the object as
// defined in Section 12.5.5 of the ECMAScript Language Specification.
// However, it has support for detecting an External value. If value has a type
// that is invalid, an error is returned.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value whose type to query.
// If the type of value is not a known ECMAScript type and value is not an
// External value napi_invalid_arg will be returned.
// N-API version: 1
func TypeOf(env Env, value Value) (ValueType, Status) {
	var res C.napi_valuetype
	var status = C.napi_typeof(env, value, &res)
	return ValueType(res), Status(status)
}

// InstanceOf function is similar to invoke the instanceof Operator on the
// object as defined in Section 12.10.4 of the ECMAScript Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] object: The JavaScript value to check.
// [in] constructor: The JavaScript function object of the constructor function
// to check against.
// N-API version: 1
func InstanceOf(env Env, object Value, constructor Value) (bool, Status) {
	var res C.bool
	var status = C.napi_instanceof(env, object, constructor, &res)
	return bool(res), Status(status)
}

// IsArray function is similar to invoke the IsArray operation on the object
// as defined in Section 7.2.2 of the ECMAScript Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func IsArray(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_is_array(env, value, &res)
	return bool(res), Status(status)
}

// IsArrayBuffer function checks if the Object passed in is an array buffer.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func IsArrayBuffer(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_is_arraybuffer(env, value, &res)
	return bool(res), Status(status)
}

// IsBuffer function  checks if the Object passed in is a buffer.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func IsBuffer(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_is_buffer(env, value, &res)
	return bool(res), Status(status)
}

// IsTypedArray function checks if the Object passed in is a typed array.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func IsTypedArray(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_is_typedarray(env, value, &res)
	return bool(res), Status(status)
}

// IsDataview function checks if the Object passed in is a DataView.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func IsDataview(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_is_dataview(env, value, &res)
	return bool(res), Status(status)
}

// StrictEquals function is simnilar to invoke the Strict Equality algorithm
// as defined in Section 7.2.14 of the ECMAScript Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] lhs: The JavaScript value to check.
// [in] rhs: The JavaScript value to check against.
// N-API version: 1
func StrictEquals(env Env, lhs Value, rhs Value) (bool, Status) {
	var res C.bool
	var status = C.napi_strict_equals(env, lhs, rhs, &res)
	return bool(res), Status(status)
}

// GetPropertyNames function returns the names of the enumerable properties
// of object as an array of strings. The properties of object whose key is a
// symbol will not be included.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the properties.
// N-API version: 1
func GetPropertyNames(env Env, object Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_property_names(env, object, &res)
	return Value(res), Status(status)
}

// SetProperty function set a property on the Object passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object on which to set the property.
// [in] key: The name of the property to set.
// [in] value: The property value.
// N-API version: 1
func SetProperty(env Env, object Value, key Value, value Value) Status {
	var status = C.napi_set_property(env, object, key, value)
	return Status(status)
}

// GetProperty function gets the requested property from the Object
// passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the property.
// [in] key: The name of the property to retrieve.
// N-API version: 1
func GetProperty(env Env, object Value, key Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_property(env, object, key, &res)
	return Value(res), Status(status)
}

// HasProperty function checks if the Object passed in has the named
// property.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] key: The name of the property whose existence to check.
// N-API version: 1
func HasProperty(env Env, object Value, key Value) (bool, Status) {
	var res C.bool
	var status = C.napi_has_property(env, object, key, &res)
	return bool(res), Status(status)
}

// DeleteProperty function attempts to delete the key own property from
// object.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] key: The name of the property to delete.
// N-API version: 1
func DeleteProperty(env Env, object Value, key Value) (bool, Status) {
	var res C.bool
	var status = C.napi_delete_property(env, object, key, &res)
	return bool(res), Status(status)
}

// HasOwnProperty function checks if the Object passed in has the named own
// property. key must be a string or a Symbol, or an error will be thrown. N-API
// will not perform any conversion between data types.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] key: The name of the own property whose existence to check.
// N-API version: 1
func HasOwnProperty(env Env, object Value, key Value) (bool, Status) {
	var res C.bool
	var status = C.napi_has_own_property(env, object, key, &res)
	return bool(res), Status(status)
}

// SetNamedProperty function set a property on the Object passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object on which to set the property.
// [in] utf8Name: The name of the property to set.
// [in] value: The property value.
// N-API version: 1
func SetNamedProperty(env Env, object Value, key string, value Value) Status {
	var ckey = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var status = C.napi_set_named_property(env, object, ckey, value)
	return Status(status)
}

// GetNamedProperty function gets the requested property from the Object
// passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the property.
// [in] utf8Name: The name of the property to get.
// N-API version: 1
func GetNamedProperty(env Env, object Value, key string) (Value, Status) {
	var res C.napi_value
	var ckey = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var status = C.napi_get_named_property(env, object, ckey, &res)
	return Value(res), Status(status)
}

// HasNamedProperty function checks if the Object passed in has the named
// property.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] utf8Name: The name of the property whose existence to check.
// N-API version: 1
func HasNamedProperty(env Env, object Value, key string) (bool, Status) {
	var res C.bool
	var ckey = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var status = C.napi_has_named_property(env, object, ckey, &res)
	return bool(res), Status(status)
}

// SetElement function sets and element on the Object passed in.
// [in] object: The object from which to set the properties.
// [in] index: The index of the property to set.
// [in] value: The property value.
// N-API version: 1
func SetElement(env Env, object Value, index uint, value Value) Status {
	var status = C.napi_set_element(env, object, C.uint32_t(index), value)
	return Status(status)
}

// GetElement function gets the element at the requested index.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the property.
// [in] index: The index of the property to get.
// N-API version: 1
func GetElement(env Env, object Value, index uint) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_element(env, object, C.uint32_t(index), &res)
	return Value(res), Status(status)
}

// HasElement function returns if the Object passed in has an element at the
// requested index.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] index: The index of the property whose existence to check.
// N-API version: 1
func HasElement(env Env, object Value, index uint) (bool, Status) {
	var res C.bool
	var status = C.napi_has_element(env, object, C.uint32_t(index), &res)
	return bool(res), Status(status)
}

// DeleteElement function attempts to delete the specified index from object.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] index: The index of the property to delete.
// N-API version: 1
func DeleteElement(env Env, object Value, index uint) (bool, Status) {
	var res C.bool
	var status = C.napi_delete_element(env, object, C.uint32_t(index), &res)
	return bool(res), Status(status)
}

func unixNano(env Env, info CallbackInfo) Value {
	fmt.Println("unixNano ...")
	now := time.Now()
	value, _ := CreateInt64(env, now.UnixNano())
	return value
}

type Ctx struct {
	caller *Caller
	env    Env
}

var ctxMap = make(map[C.napi_env]Ctx)

//export GoHandler
func GoHandler(env C.napi_env, info C.napi_callback_info) C.napi_value {
	fmt.Printf("Env in GoHandler %p\n", env)
	ctx := ctxMap[env]
	caller := ctx.caller
	return (C.napi_value)(caller.Cb(Env(env), CallbackInfo(info)))
}

func DefineProperties(env Env, value Value, properties []Property) Status {
	raw := make([]PropertyDescriptor, len(properties))
	for i := range properties {
		properties[i].getRaw()
		//raw[i] = prop
	}
	name := C.CString("unixNano")
	defer C.free(unsafe.Pointer(name))
	hello, _ := CreateStringUtf8(env, "hello")
	caller := &Caller{
		Cb: unixNano,
	}
	ctxMap[env] = Ctx{
		caller: caller,
		env:    env,
	}
	fmt.Printf("Env in define %p\n", env)
	desc := PropertyDescriptor{
		utf8name:   name,
		name:       nil,
		method:     (Callback)(unsafe.Pointer(C.InvokeGoHandler)),
		getter:     nil,
		setter:     nil,
		value:      hello,
		attributes: C.napi_default,
		data:       nil,
	}
	raw[0] = desc
	var props = (unsafe.Pointer(&raw[0]))
	var status = C.napi_define_properties(env, value, C.size_t(len(properties)), (*C.napi_property_descriptor)(props))
	return Status(status)
}

// Working with JavaScript Functions
// N-API provides a set of APIs that allow JavaScript code to call back into
// native code.  N-API APIs that support calling back into native code take in a
// callback functions represented by the NapiCallback type.
// When the JavaScript VM calls back to native code, the NapiCallback function
// provided is invoked.
// Additionally, N-API provides a set of functions which allow calling JavaScript
// functions from native code. One can either call a function like a regular
// JavaScript function call, or as a constructor function.

// Any non-NULL data which is passed to this API via the data field of the
// NapiPropertyDescriptor items can be associated with object and freed whenever
// object is garbage-collected by passing both object and the data to
// NapiAddFinalizer.

// CallFunction function allows a JavaScript function object to be called
// from a native add-on. This is the primary mechanism of calling back from the
// add-on's native code into JavaScript.
// For the special case of calling into JavaScript after an async operation,
// see NapiMakeCallback.
// [in] env: The environment that the API is invoked under.
// [in] recv: The this object passed to the called function.
// [in] func: napi_value representing the JavaScript function to be invoked.
// [in] argc: The count of elements in the argv array.
// [in] argv: Array of napi_values representing JavaScript values passed in as
// arguments to the function.
// N-API version: 1
func CallFunction(env Env, receiver Value, function Value, arguments []Value) (Value, Status) {
	var res C.napi_value
	var args = unsafe.Pointer(&arguments[0])
	var status = C.napi_call_function(env, receiver, function, C.size_t(len(arguments)), (*C.napi_value)(args), &res)
	return Value(res), Status(status)
}

// CreateFunction function allows an add-on author to create a function
// object in native code. This is the primary mechanism to allow calling into
// the add-on's native code from JavaScript.
// [in] env: The environment that the API is invoked under.
// [in] utf8Name: The name of the function encoded as UTF8. This is visible
// within JavaScript as the new function object's name property.
// [in] length: The length of the utf8name in bytes, or NAPI_AUTO_LENGTH if it
// is null-terminated.
// [in] cb: The native function which should be called when this function object
// is invoked.
// [in] data: User-provided data context. This will be passed back into the
// function when invoked later.
// N-API version: 1
func CreateFunction(env Env, name string, cb CCallback) (Value, Status) {
	/*caller := &Caller{
		Cb: cb,
	}*/
	var res C.napi_value
	/*var cname = C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var status = C.napi_create_function(env, cname, C.NAPI_AUTO_LENGTH, (Callback)(C.Callback(unsafe.Pointer(caller))), nil, &res)*/
	var status = C.napi_ok
	return Value(res), Status(status)
}

// GetCbInfo function is used within a callback function to retrieve details
// about the call like the arguments and the this pointer from a given callback
// info.
// [in] env: The environment that the API is invoked under.
// [in] cbinfo: The callback info passed into the callback function.
// [in-out] argc: Specifies the size of the provided argv array and receives the
// actual count of arguments.
// [out] argv: Buffer to which the napi_value representing the arguments are copied. If there are more arguments than the provided count, only the requested number of arguments are copied. If there are fewer arguments provided than claimed, the rest of argv is filled with napi_value values that represent undefined.
// [out] this: Receives the JavaScript this argument for the call.
// [out] data: Receives the data pointer for the callback.
// N-API version: 1
func GetCbInfo(env Env, cbinfo CallbackInfo) ([]Value, Value, unsafe.Pointer, Status) {
	var staticArgc C.size_t = 5
	var argc = staticArgc
	staticArgv := make([]Value, 5)
	var dynamicArgv []Value
	var cStaticArgv = (unsafe.Pointer(&staticArgv[0]))
	var thisArg C.napi_value
	var data unsafe.Pointer
	var status = C.napi_get_cb_info(env, cbinfo, &argc, (*C.napi_value)(cStaticArgv), &thisArg, &data)
	if argc > staticArgc {
		dynamicArgv = make([]Value, int(argc))
		var cDynamicArgv = (unsafe.Pointer(&dynamicArgv[0]))
		var status = C.napi_get_cb_info(env, cbinfo, &argc, (*C.napi_value)(cDynamicArgv), nil, nil)
		return dynamicArgv, Value(thisArg), data, Status(status)
	}
	arguments := make([]Value, int(argc))
	for i := 0; i < int(argc); i++ {
		arguments[i] = staticArgv[i]
	}
	return arguments, Value(thisArg), data, Status(status)
}

// GetNewTarget function returns the new.target of the constructor call. If
// the current callback is not a constructor call, the result is NULL.
// [in] env: The environment that the API is invoked under.
// [in] cbinfo: The callback info passed into the callback function.
// N-API version: 1
func GetNewTarget(env Env, cbinfo CallbackInfo) (Value, Status) {
	var res C.napi_value
	var status = C.napi_get_new_target(env, cbinfo, &res)
	return Value(res), Status(status)
}

// NewInstance function  is used to instantiate a new JavaScript value using
// a given NapiValue that represents the constructor for the object.
// [in] env: The environment that the API is invoked under.
// [in] cons: napi_value representing the JavaScript function to be invoked as a
// constructor.
// [in] argc: The count of elements in the argv array.
// [in] argv: Array of JavaScript values as napi_value representing the
// arguments to the constructor.
// [out] result: napi_value representing the JavaScript object returned, which in
// this case is the constructed object.
// N-API version: 1
func NewInstance(env Env, ctor Value, arguments []Value) (Value, Status) {
	var res C.napi_value
	var args = unsafe.Pointer(&arguments[0])
	// defer C.free(args)
	var status = C.napi_new_instance(env, ctor, C.size_t(len(arguments)), (*C.napi_value)(args), &res)
	return Value(res), Status(status)
}

//Object Wrap
// N-API offers a way to "wrap" C++ classes and instances so that the class
// constructor and methods can be called from JavaScript.
// The NapiDefineClass function defines a JavaScript class with constructor, s
// tatic properties and methods, and instance properties and methods that
// correspond to the C++ class.
// When JavaScript code invokes the constructor, the constructor callback uses
// NapiWrap to wrap a new C++ instance in a JavaScript object, then returns the
// wrapper object.
// When JavaScript code invokes a method or property accessor on the class, the
// corresponding NapiCallback C++ function is invoked.
// For wrapped objects it may be difficult to distinguish between a function
// called on a class prototype and a function called on an instance of a class.
// A common pattern used to address this problem is to save a persistent
// reference to the class constructor for later instanceof checks.

// DefineClass function defines a JavaScript class that corresponds to
// a C++ class.
// The C++ constructor callback should be a static method on the class that calls
// the actual class constructor, then wraps the new C++ instance in a JavaScript
// object, and returns the wrapper object.
// The JavaScript constructor function returned from napi_define_class is often
// saved and used later, to construct new instances of the class from native
// code, and/or check whether provided values are instances of the class. In that
// case, to prevent the function value from being garbage-collected, create a
// persistent reference to it using NapiCreateReference and ensure the
// reference count is kept >= 1.
// [in] env: The environment that the API is invoked under.
// [in] utf8name: Name of the JavaScript constructor function; this is not
// required to be the same as the C++ class name, though it is recommended for
// clarity.
// [in] length: The length of the utf8name in bytes, or NAPI_AUTO_LENGTH if it
// is null-terminated.
// [in] constructor: Callback function that handles constructing instances of
// the class. (This should be a static method on the class, not an actual C++
// constructor function.)
// [in] data: Optional data to be passed to the constructor callback as the data
// property of the callback info.
// [in] property_count: Number of items in the properties array argument.
// [in] properties: Array of property descriptors describing static and instance
// data properties, accessors, and methods on the class.
// See documentation for NapiPropertyDescriptor function.
// [out] result: A napi_value representing the constructor function for the
// class.
// Any non-NULL data which is passed to this API via the data parameter or via
// the data field of the NapiPropertyDescriptor array items can be associated
// with the resulting JavaScript constructor (which is returned in the result
// parameter) and freed whenever the class is garbage-collected by passing both
// the JavaScript function and the data to NapiAddFinalizer.
// N-API version: 1
func DefineClass(env Env, name string, ctor Callback, properties []PropertyDescriptor) (Value, Status) {
	var res C.napi_value
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	props := unsafe.Pointer(&properties[0])
	var status = C.napi_define_class(env, cname, C.NAPI_AUTO_LENGTH, ctor, nil, C.size_t(len(properties)), (*C.napi_property_descriptor)(props), &res)
	return Value(res), Status(status)
}

// Wrap function wraps a native instance in a JavaScript object. The native
// instance can be retrieved later using NapiUnwrap().
// [in] env: The environment that the API is invoked under.
// [in] js_object: The JavaScript object that will be the wrapper for the native
// object.
// [in] native_object: The native instance that will be wrapped in the
// JavaScript object.
// [in] finalize_cb: Optional native callback that can be used to free the native
// instance when the JavaScript object is ready for garbage-collection.
// [in] finalize_hint: Optional contextual hint that is passed to the finalize
// callback.
// [out] result: Optional reference to the wrapped object.
// When JavaScript code invokes a constructor for a class that was defined using
// NapiDefineClass(), the NapiCallback for the constructor is invoked. After
// constructing an instance of the native class, the callback must then call
// NapiWrap() to wrap the newly constructed instance in the already-created
// JavaScript object that is the this argument to the constructor callback. That
// this object was created from the constructor function's prototype, so it
// already has definitions of all the instance properties and methods.
// Typically when wrapping a class instance, a finalize callback should be
// provided that simply deletes the native instance that is received as the data
// argument to the finalize callback.
// The optional returned reference is initially a weak reference, meaning it has
// a reference count of 0. Typically this reference count would be incremented
// temporarily during async operations that require the instance to remain valid.
// N-API version: 1
func Wrap(env Env, value Value, native unsafe.Pointer) (Ref, Status) {
	var res C.napi_ref
	// TODO napi_wrap(napi_env env, napi_value js_object, void* native_object, napi_finalize finalize_cb, void* finalize_hint, napi_ref* result);
	var status = C.napi_wrap(env, value, native, nil, nil, &res)
	return Ref(res), Status(status)
}

// Unwrap function retrieves a native instance that was previously wrapped
// in a JavaScript object using NapiWrap().
// [in] env: The environment that the API is invoked under.
// [in] js_object: The object associated with the native instance.
// [out] result: Pointer to the wrapped native instance.
// When JavaScript code invokes a method or property accessor on the class, the
// corresponding NapiCallback is invoked. If the callback is for an instance
// method or accessor, then the this argument to the callback is the wrapper
// object; the wrapped C++ instance that is the target of the call can be
// obtained then by calling NapiUnwrap() on the wrapper object.
// N-API version: 1
func Unwrap(env Env, value Value) (unsafe.Pointer, Status) {
	var res unsafe.Pointer
	// napi_remove_wrap(napi_env env, napi_value js_object, void** result)
	var status = C.napi_unwrap(env, value, &res)
	return res, Status(status)
}

// RemoveWrap function retrieves a native instance that was previously
// wrapped in the JavaScript object using NapiWrap() and removes the wrapping.
// If a finalize callback was associated with the wrapping, it will no longer be
// called when the JavaScript object becomes garbage-collected.
// [in] env: The environment that the API is invoked under.
// [in] js_object: The object associated with the native instance.
// [out] result: Pointer to the wrapped native instance.
// N-API version: 1
func RemoveWrap(env Env, value Value) (unsafe.Pointer, Status) {
	var res unsafe.Pointer
	// TODO napi_remove_wrap(napi_env env, napi_value js_object, void** result)s
	var status = C.napi_remove_wrap(env, value, &res)
	return res, Status(status)
}

// AddFinalizer function adds a NapiFinalize callback which will be called
// when the JavaScript object is ready for garbage collection.
// [in] env: The environment that the API is invoked under.
// [in] js_object: The JavaScript object to which the native data will be
// attached.
// [in] native_object: The native data that will be attached to the JavaScript
// object.
// [in] finalize_cb: Native callback that will be used to free the native data
// when the JavaScript object is ready for garbage-collection.
// [in] finalize_hint: Optional contextual hint that is passed to the finalize
// callback.
// [out] result: Optional reference to the JavaScript object.
// This API is similar to NapiWrap() except that:
//  - the native data cannot be retrieved later using Napinwrap(),
//  - nor can it be removed later using NapiRemoveWrap(),
//  - the API can be called multiple times with different data items in order to
//    attach each of them to the JavaScript object.
// Caution: The optional returned reference (if obtained) should be deleted via
// NapiDeleteReference ONLY in response to the finalize callback invocation. If
// it is deleted before, then the finalize callback may never be invoked.
// Therefore, when obtaining a reference a finalize callback is also required in
// order to enable correct disposal of the reference.
// N-API version: 1
func AddFinalizer(env Env, obj Value, native unsafe.Pointer, finalizer *FinalizeCaller, hint unsafe.Pointer) (Ref, Status) {
	var res C.napi_ref
	var fn = C.FinalizeCallback(unsafe.Pointer(finalizer))
	var status = C.napi_add_finalizer(env, obj, native, fn, hint, &res)
	return Ref(res), Status(status)
}

// Simple Asynchronous Operations
// Add-on modules often need to leverage asynchronous helpers from libuv as part
// of their implementation. This allows them to schedule work to be executed
// asynchronously so that their methods can return in advance of the working
// being completed. This is important in order to allow them to avoid blocking
// overall execution of the Node.js application.
// N-API provides an ABI-stable interface for these supporting functions which
// covers the most common asynchronous use cases.

// N-API defines work structure which is used asynchronous worker. Insstances are
// created or deleted with CreateAsyncWork and DeleteAssyncWork. The execute and
// complete callbacks are functions that will be invoked when the executor is
// ready to execute and when it completes its task respectively.

// The execute function should avoid making any N-API calls that could result in
// the execution of JavaSscript of interaction with JavaScript objects. Most
// often, any code that needs to make N-API calls should be made in the complete
// callback instead. Avoid unsing the Env parameter in the execute callback as it
// will likely execute JavaScript.

// These functions implement the following interfaces:
// typedef void (*napi_async_execute_callback)(napi_env env, void* data);
// typedef void (*napi_async_complete_callback)(napi_env env, napi_status status, void* data);
// When these methods are invoked, the data parameter passed will be the add-on
// provided data that was passed into CreateAsyncWork function.

// Once created the asynchronous worker can be queued for the execution using
// QueueAsyncWork function. The CancelAyncWork function can be used if the work
// needs to be cancelled before the work has started its execution. After calling
// CancelAsyncWork the complete callback will be invoked with a status of
// Cancelled. The work should not be deleted before the complete callback
// invocation, even when it was cancelled.

// CreateAsyncWork function allocates a work object that is used to execute logic
// asynchronously.
// [in] env: The environment that the API is invoked under.
// [in] resource: An optional object associated with the async work that will be
// passed to possible async hooks.
// [in] name: Identifier for the kind of resource that is being provided for
// diagnostic information exposed by the async hooks API.
// [in] execute: The native function which should be called to execute the logic
// asynchronously. The given function is called from a worker pool thread and can
// execute in parallel with the main event loop thread.
// [in] complete: The native function which will be called when the asynchronous
// logic is completed or is cancelled. The given function is called from the main
// event loop thread.
// [in] data: User-provided data context. This will be passed back into the
// execute and complete functions.
// [out] result: Returns the handle to the newly created async work.
// N-API version: 1
func CreateAsyncWork(env Env, resource Value, name Value, execute *AsyncExecuteCaller, complete *AsyncCompleteCaller, data unsafe.Pointer) (AsyncWork, Status) {
	var res C.napi_async_work
	cexecute := C.AsyncExecuteCallback(unsafe.Pointer(execute))
	ccomplete := C.AsyncCompleteCallback(unsafe.Pointer(complete))
	cdata := unsafe.Pointer(data)
	var status = C.napi_create_async_work(env, resource, name, cexecute, ccomplete, cdata, &res)
	return AsyncWork(res), Status(status)
}

// DeleteAsyncWork function frees the previously allocated work object.
// This API can be called even if there is a pending JavaScript exception.
// [in] env: The environment that the API is invoked under.
// [in] work: The handle returned by the call to CreateAsyncWork.
// N-API version: 1
func DeleteAsyncWork(env Env, work AsyncWork) Status {
	var status = C.napi_delete_async_work(env, work)
	return Status(status)
}

// QueueAsyncWork function requests that the previously allocated work be
// scheduled for the execution.
// Once it returns successfully, this function must not be called again with the
// same AsyncWork item or the result will be undefined.
// [in] env: The environment that the API is invoked under.
// [in] work: The handle returned by the call to CreateAsyncWork function.s
// N-API version: 1
func QueueAsyncWork(env Env, work AsyncWork) Status {
	var status = C.napi_queue_async_work(env, work)
	return Status(status)
}

// CancelAsyncWork function cancels queued work if it has not yet been started.
// If it has already started executing, it cannot be cancelled and GenericFailure
// will be returned. If successful, the complete callback will be invoked with a
// status value of Cancelled. The work should not be deleted before the complete
// callback invocation, even if it has been successfully cancelled.
// This API can be called even if there is a pending JavaScript exception.
// [in] env: The environment that the API is invoked under.
// [in] work: The handle returned by the call to CreateAsyncWork.
// N-API version: 1
func CancelAsyncWork(env Env, work AsyncWork) Status {
	var status = C.napi_cancel_async_work(env, work)
	return Status(status)
}

// Custom Asynchronous Operations
// The simple asynchronous work may not be appropriate for every scenario. When
// using any other asynchronous mechanism, the following APIs are necessary to
// ensure an asynchronous operation is properly tracked by the runtime.

// AsyncInit function initializes a new asynchronous context.
// [in] env: The environment that the API is invoked under.
// [in] resource: An optional object associated with the asynchronous work that
// will be passed to possible async hooks.
// [in] name: Identifier for the kind of resource that is being provided for
// diagnostic information exposed by the async hooks.
// [out] result: The initialized asynchronous context.
// N-API version: 1
func AsyncInit(env Env, resource Value, name Value) (AsyncContext, Status) {
	var res C.napi_async_context
	var status = C.napi_async_init(env, resource, name, &res)
	return AsyncContext(res), Status(status)
}

// AsyncDestroy function destroys the passed asynchronous context. The function
// can be called even if there is a pending JavaScript exception.
// [in] env: The environment that the API is invoked under.
// [in] ctx: The asynchronous context to be destroyed.
// N-API version: 1
func AsyncDestroy(env Env, ctx AsyncContext) Status {
	var status = C.napi_async_destroy(env, ctx)
	return Status(status)
}

// MakeCallback function allows a JavaScript function object to be called from
// a native add-on.
// This function is similar to CallFunction. However, it is used to call from
// native code back into JavaScript after running from async operation.
// [in] env: The environment that the API is invoked under.
// [in] ctx: Context for the async operation that is invoking the callback. This
// should normally be a value previously obtained from AsyncInit function.
// However nil is also allowed, which indicates the current async context
// (if any) is to be used for the callback.
// [in] recv: The this object passed to the called function.
// [in] fn: The JavaScript function to be invoked.
// [in] args: Slice of JavaScript values as Value representing the arguments to
// the function.
// [out] result: Value representing the JavaScript object returned.
// N-API version: 1
func MakeCallback(env Env, ctx AsyncContext, recv Value, fn Value, args []Value) (Value, Status) {
	var res C.napi_value
	var argv = (*C.napi_value)(unsafe.Pointer(&args[0]))
	var argc = C.size_t(len(args))
	var status = C.napi_make_callback(env, ctx, recv, fn, argc, argv, &res)
	return Value(res), Status(status)
}

// OpenCallbackScope function opens the required scope.
// [in] env: The environment that the API is invoked under.
// [in] resource: An object associated with the async work that will be passed to
// possible async hooks .
// [in] ctx: Context for the async operation that is invoking the callback. This
// should be a value previously obtained from Asyncinit.
// [out] result: The newly created scope.
// There are cases where it is necessary to have the equivalent of the scope
// associated with a callback in place when making certain N-API calls.
// N-API version: 3
func OpenCallbackScope(env Env, resource Value, ctx AsyncContext) (CallbackScope, Status) {
	var scope C.napi_callback_scope
	var status = C.napi_open_callback_scope(env, resource, ctx, &scope)
	return CallbackScope(scope), Status(status)
}

// CloseCallbackScope function closes the required scope.
// [in] env: The environment that the API is invoked under.
// [in] scope: The scope to be closed.
// This function can be called even if there is a pending JavaScript exception.
// N-API version: 3
func CloseCallbackScope(env Env, scope CallbackScope) Status {
	var status = C.napi_close_callback_scope(env, scope)
	return Status(status)
}

// GetNodeVersion function fills the version struct with the major, minor,
// and patch version of Node.js that is currently running, and the release field
// with the value of process.release.name.
// [in] env: The environment that the API is invoked under.
// The returned buffer is statically allocated and does not need to be freed.
// N-API version: 1
func GetNodeVersion(env Env) (NodeVersion, Status) {
	var res *C.napi_node_version
	var status = C.napi_get_node_version(env, &res)
	return NodeVersion(res), Status(status)
}

// GetVersion function returns the highest version of N-API supported.
// [in] env: The environment that the API is invoked under.
// This function returns the highest N-API version supported by the Node.js
// runtime. N-API is planned to be additive such that newer releases of Node.js
// may support additional API functions. In order to allow an addon to use a
// newer function when running with versions of Node.js that support it, while
// providing fallback behavior when running with Node.js versions that don't
// support it.
// N-API version: 1
func GetVersion(env Env) (uint32, Status) {
	var res C.uint32_t
	var status = C.napi_get_version(env, &res)
	return uint32(res), Status(status)
}

// AdjustExternalMemory function gives V8 an indication of the amount of
// externally allocated memory that is kept alive by JavaScript objects
// (i.e. a JavaScript object that points to its own memory allocated by a native
// module). Registering externally allocated memory will trigger global garbage
// collections more often than it would otherwise.
// [in] env: The environment that the API is invoked under.
// [in] change_in_bytes: The change in externally allocated memory that is kept
// alive by JavaScript objects.
// N-API version: 1
func AdjustExternalMemory(env Env, changeInBytes int64) (int64, Status) {
	var res C.int64_t
	var status = C.napi_adjust_external_memory(env, C.int64_t(changeInBytes), &res)
	return int64(res), Status(status)
}

// Promises
// N-API provides facilities for creating Promise objects as described in
// Section 25.4 of the ECMA specification. It implements promises as a pair of
// objects.
// When a promise is created a "deferred" object is created and returned
// alongside the Promise. The deferred object is bound to the created Promise and
// is the only means to resolve or reject the Promise. The deferred object will
// be freed when the Promise will be resolde or rejected. The Promise object may
// be returned to JavaScript where it can be used in the usual fashion.

// CreatePromise function creates a deferred object and a JavaScript promise.
// [in] env: The environment that the API is invoked under.
// [out] deferred: A newly created deferred object which can later be used to
// resolve or reject the associated promise.
// [out] promise: The JavaScript promise associated with the deferred object.
// N-API version: 1
func CreatePromise(env Env) (Value, Deferred, Status) {
	var promise C.napi_value
	var deferred C.napi_deferred
	var status = C.napi_create_promise(env, &deferred, &promise)
	return Value(promise), Deferred(deferred), Status(status)
}

// ResolveDeferred function resolves a JavaScript promise by way of the deferred
// object with which it is asssociated. It can only be used to resolve JavaScript
// promises for which the corresponding deferred object is available.
// [in] env: The environment that the API is invoked under.
// [in] deferred: The deferred object whose associated promise to resolve.
// [in] rejection: The value with which to reject the promise.
// The deferred object is freed upon successful completion.
// N-API version: 1
func ResolveDeferred(env Env, deferred Deferred, resolution Value) Status {
	var status = C.napi_resolve_deferred(env, deferred, resolution)
	return Status(status)
}

// RejectDeferred function rejects a JavaScript promise by way of the deferred
// object with which it is associated. It can only be used to reject JavaScript
// promise for which the corresponding deferred object is available.
// [in] env: The environment that the API is invoked under.
// [in] deferred: The deferred object whose associated promise to resolve.
// [in] rejection: The value with which to reject the promise.
// The deferred object is freed upon successful completion.
// N-API version: 1
func RejectDeferred(env Env, deferred Deferred, rejection Value) Status {
	var status = C.napi_reject_deferred(env, deferred, rejection)
	return Status(status)
}

// IsPromise function checks whether a promise object is a native promise
// object.
// [in] env: The environment that the API is invoked under.
// [in] promise: The promise to examine
// [out] is_promise: Flag indicating whether promise is a native promise object
// that is, a promise object created by the underlying JavaScript engine.
// N-API version: 1
func IsPromise(env Env, value Value) (bool, Status) {
	var res C.bool
	var status = C.napi_is_promise(env, value, &res)
	return bool(res), Status(status)
}

// RunScript function function execute a JavaScript script passend in like a
// string.
// [in] env: The environment that the API is invoked under.
// [in] script: A JavaScript string containing the script to execute.
// [out] result: The value resulting from having executed the script.
// N-API version: 1
func RunScript(env Env, script Value) (Value, Status) {
	var res C.napi_value
	var status = C.napi_run_script(env, script, &res)
	return Value(res), Status(status)
}

// GetUvEventLoop function retrieves the current event loop associated with a
// specific environment.
// [in] env: The environment that the API is invoked under.
// [out] loop: The current libuv loop instance.
// N-API version: 2
func GetUvEventLoop(env Env) (UVLoop, Status) {
	var loop UVLoop
	var status = C.napi_get_uv_event_loop(env, &loop)
	return loop, Status(status)
}

// Asynchronous Thread-safe Function Calls
// JavaScript functions can normally only be called from a native addon's main
// thread. If an addon creates additional threads, then N-API functions that
// require a Env, Value, or Ref must not be called from those threads.
// When an addon has additional threads and JavaScript functions need to be
// invoked based on the processing completed by those threads, those threads must
// communicate with the addon's main thread so that the main thread can invoke
// the JavaScript function on their behalf. The thread-safe function APIs provide
// an easy way to do this.

// CreateThreadsafeFunction function ...
// [in] env: The environment that the API is invoked under.
// [in] fn: An optional JavaScript function to call from another thread. It must
// be provided if nil is passed to call_js_cb.
// [in] resource: An optional object associated with the async work that will be
// passed to possible async hooks.
// [in] maxQueueSize: Maximum size of the queue. 0 for no limit.
// [in] initialThreadCount: The initial number of threads, including the main
// thread, which will be making use of this function.
// [in] data: Optional data to be passed to finalize.
// [in] finalizer: Optional function to call when the thread-safe function is
// being destroyed.
// [in] context: Optional data to attach to the resulting thread-safe function.
// [in] js: Optional callback which calls the JavaScript function in response to
// a call on a different thread. This callback will be called on the main thread.
// If not given, the JavaScript function will be called with no parameters and
// with undefined as its this value.
// N-API version: 4
func CreateThreadsafeFunction(env Env, fn Value, resource Value, name Value, maxQueueSize uint, initialThreadCount uint, data unsafe.Pointer, finalizer *FinalizeCaller, ctx unsafe.Pointer, tsfn *ThreadsafeFunctionsCaller) (ThreadsafeFunction, Status) {
	var res C.napi_threadsafe_function
	ctsfn := C.ThreadsafeFunctionCallback(unsafe.Pointer(tsfn))
	cfinalize := C.FinalizeCallback(unsafe.Pointer(finalizer))
	var status = C.napi_create_threadsafe_function(env, fn, resource, name, C.size_t(maxQueueSize), C.size_t(initialThreadCount), data, cfinalize, ctx, ctsfn, &res)
	return ThreadsafeFunction(res), Status(status)
}

// GetThreadsafeFunctionContext function ...
// [in] func: The thread-safe function for which to retrieve the context.
// [out] result: The location where to store the context.
// This API may be called from any thread which makes use of thread-safe
// function.
// N-API version: 4
func GetThreadsafeFunctionContext(fn ThreadsafeFunction) (unsafe.Pointer, Status) {
	var res unsafe.Pointer
	var status = C.napi_get_threadsafe_function_context(fn, &res)
	return res, Status(status)
}

// CallThreadsafeFunction function ...
// [in] fn: The asynchronous thread-safe JavaScript function to invoke.
// [in] data: Data to send into JavaScript via the callback provided during the
// creation of the thread-safe JavaScript function.
// [in] isBlocking: Flag whose value can be either NapiTsfnBlocking to indicate
// that the call should block if the queue is full or NapiTsfnNonBlocking to
// indicate that the call should return immediately with a status of
// QueueFull whenever the queue is full.
// This function may be called from any thread which makes use of the thread-safe
// function.
// N-API version: 4
func CallThreadsafeFunction(fn ThreadsafeFunction, data unsafe.Pointer, mode ThreadsafeFunctionCallMode) Status {
	var status = C.napi_call_threadsafe_function(fn, data, mode)
	return Status(status)
}

// AcquireThreadsafeFunction function should be called before use thread-safe
// function. This prevents the thread-safe function to be destroyed when all
// other threads have stopped making use of it. This function maybe called from
// any thread which will start making use of thread-safe function.
// N-API version: 4
func AcquireThreadsafeFunction(fn ThreadsafeFunction) Status {
	var status = C.napi_acquire_threadsafe_function(fn)
	return Status(status)
}

// ReleaseThreadsafeFunction function should be called when a thread stops making
// use of thread-safe function. Passing the thread-safe function to any
// thread-safe APIs could result in an undefined behaviour. This function may be
// called from any thread which will stop making use of thread-safe function.
// [in] fn: The asynchronous thread-safe JavaScript function whose reference
// count to decrement.
// [in] mode: Flag whose value can be:
//  - napi_tsfn_release: the current thread will make no further calls to
//  thread-safe function.
//  - napi_tsfn_abort: the current thread, no other htread should make any
//  further calls to thread-safe function.
// If set to napi_tsfn_abort, further calls to napi_call_threadsafe_function()
// will return napi_closing and no further values will be placed in the queue.
// N-API version: 4
func ReleaseThreadsafeFunction(fn ThreadsafeFunction, mode TheradsafeFunctionReleaseMode) Status {
	var status = C.napi_release_threadsafe_function(fn, mode)
	return Status(status)
}

// RefThreadsafeFunction function is used to indicate that the event loop running
// on the main thread should not exit until the threadsafe function has been
// destroyed. This function may only be called from the main thread.
// [in] env: The environment that the API is invoked under.
// [in] fn: The thread-safe function to reference.
// N-API version: 4
func RefThreadsafeFunction(env Env, fn ThreadsafeFunction) Status {
	var status = C.napi_ref_threadsafe_function(env, fn)
	return Status(status)
}

// UnrefThreadsafeFunction function is used to indicate that the event loop
// running on the main thread may exit before the threadsafe function is
// destroyed. This function may only be called from the main thread.
// [in] env: The environment that the API is invoked under.
// [in] fn: The thread-safe function to unreference.
// N-API version: 4
func UnrefThreadsafeFunction(env Env, fn ThreadsafeFunction) Status {
	var status = C.napi_unref_threadsafe_function(env, fn)
	return Status(status)
}

// CCallback  ...
type CCallback func(Env, CallbackInfo) Value

// Caller contains a callback to call
type Caller struct {
	Cb CCallback
}

//export CallCallback
func CallCallback(wrap unsafe.Pointer, env C.napi_env, info C.napi_callback_info) C.napi_value {
	caller := (*Caller)(wrap)
	return (C.napi_value)(caller.Cb(Env(env), CallbackInfo(info)))
}

// CAsyncExecuteCallback  ...
type CAsyncExecuteCallback func(Env, unsafe.Pointer)

// AsyncExecuteCaller contains a callback to call
type AsyncExecuteCaller struct {
	Cb CAsyncExecuteCallback
}

//export CallAsyncExecuteCallback
func CallAsyncExecuteCallback(wrap unsafe.Pointer, env C.napi_env, data unsafe.Pointer) {
	caller := (*AsyncExecuteCaller)(wrap)
	caller.Cb(env, data)
}

// CAsyncExecuteCallback  ...
type CAsyncCompleteCallback func(Env, Status, unsafe.Pointer)

// AsyncExecuteCaller contains a callback to call
type AsyncCompleteCaller struct {
	Cb CAsyncCompleteCallback
}

//export CallAsyncCompleteCallback
func CallAsyncCompleteCallback(wrap unsafe.Pointer, env C.napi_env, status C.napi_status, data unsafe.Pointer) {
	caller := (*AsyncCompleteCaller)(wrap)
	caller.Cb(env, status, data)
}

// CFinalizeCallback  ...
type CFinalizeCallback func(Env, unsafe.Pointer, unsafe.Pointer)

// FinalizeCaller contains a callback to call
type FinalizeCaller struct {
	Cb CFinalizeCallback
}

//export CallFinalizeCallback
func CallFinalizeCallback(wrap unsafe.Pointer, env C.napi_env, data unsafe.Pointer, hint unsafe.Pointer) {
	caller := (*FinalizeCaller)(wrap)
	caller.Cb(env, data, hint)
}

// CThreadsafeFunctionsCallback  ...
type CThreadsafeFunctionsCallback func(Env, Value, unsafe.Pointer, unsafe.Pointer)

// ThreadsafeFunctionsCaller contains a callback to call
type ThreadsafeFunctionsCaller struct {
	Cb CThreadsafeFunctionsCallback
}

//export CallThreadsafeFunctionCallback
func CallThreadsafeFunctionCallback(wrap unsafe.Pointer, env C.napi_env, fn C.napi_value, ctx unsafe.Pointer, data unsafe.Pointer) {
	caller := (*ThreadsafeFunctionsCaller)(wrap)
	caller.Cb(env, fn, ctx, data)
}

// Property ...
type Property struct {
	Name   string
	Method *Caller
}

// GetRaw ...
func (prop *Property) getRaw() PropertyDescriptor {
	name := C.CString(prop.Name)
	defer C.free(unsafe.Pointer(name))
	desc := PropertyDescriptor{
		utf8name:   name,
		name:       nil,
		method:     nil, //(Callback)(C.Callback(unsafe.Pointer(prop.Method))), //nil,
		getter:     nil,
		setter:     nil,
		value:      nil,
		attributes: C.napi_default,
		data:       nil,
	}
	return desc
}
