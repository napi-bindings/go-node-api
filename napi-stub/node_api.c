#include "node_api.h"
#include <stdio.h>

// ##############################################################################
// Stub for JavaScript API
// ##############################################################################

NAPI_EXTERN napi_status
napi_get_last_error_info(
    napi_env env, 
    const napi_extended_error_info** result) {
        return napi_ok;
}

// Getters for defined singletons
NAPI_EXTERN napi_status 
napi_get_undefined(
    napi_env env, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_null(
    napi_env env, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_global(
    napi_env env, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_boolean(
    napi_env env, 
    bool value, 
    napi_value* result) {
        return napi_ok;
}

// Methods to create Primitive types/Objects
NAPI_EXTERN napi_status 
napi_create_object(
    napi_env env, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_array(
    napi_env env, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_array_with_length(
    napi_env env, 
    size_t length, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_double(
    napi_env env, 
    double value, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_int32(
    napi_env env, 
    int32_t value, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_uint32(
    napi_env env, 
    uint32_t value, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_int64(
    napi_env env, 
    int64_t value, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_string_latin1(
    napi_env env, 
    const char* str, 
    size_t length, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_string_utf8(
    napi_env env, 
    const char* str, 
    size_t length, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_string_utf16(
    napi_env env, 
    const char16_t* str, 
    size_t length, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_symbol(
    napi_env env, 
    napi_value description, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_function(
    napi_env env, 
    const char* utf8name, 
    size_t length, 
    napi_callback cb, 
    void* data, napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_error(
    napi_env env, 
    napi_value code, 
    napi_value msg, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_type_error(
    napi_env env, 
    napi_value code, 
    napi_value msg, 
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_range_error(
    napi_env env, 
    napi_value code, 
    napi_value msg, 
    napi_value* result) {
        return napi_ok;
}

// Methods to get the native napi_value from Primitive type
NAPI_EXTERN napi_status 
napi_typeof(
    napi_env env, 
    napi_value value, 
    napi_valuetype* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_value_double(
    napi_env env, 
    napi_value value, 
    double* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_value_int32(
    napi_env env, 
    napi_value value, 
    int32_t* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_value_uint32(
    napi_env env, 
    napi_value value, 
    uint32_t* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_value_int64(
    napi_env env, 
    napi_value value, 
    int64_t* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_value_bool(
    napi_env env,
    napi_value value, 
    bool* result) {
        return napi_ok;
}

// Copies LATIN-1 encoded bytes from a string into a buffer.
NAPI_EXTERN napi_status 
napi_get_value_string_latin1(
    napi_env env, 
    napi_value value,
    char* buf, 
    size_t bufsize, 
    size_t* result) {
        return napi_ok;
}

// Copies UTF-8 encoded bytes from a string into a buffer.
NAPI_EXTERN napi_status 
napi_get_value_string_utf8(
    napi_env env,
    napi_value value,
    char* buf,
    size_t bufsize,
    size_t* result) {
        return napi_ok;
}

// Copies UTF-16 encoded bytes from a string into a buffer.
NAPI_EXTERN napi_status 
napi_get_value_string_utf16(
    napi_env env,
    napi_value value,
    char16_t* buf,
    size_t bufsize,
    size_t* result) {
        return napi_ok;
}

// Methods to coerce values
// These APIs may execute user scripts
NAPI_EXTERN napi_status 
napi_coerce_to_bool(
    napi_env env,
    napi_value value,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_coerce_to_number(
    napi_env env,
    napi_value value,
    napi_value* result) {
        return napi_ok;
    }
NAPI_EXTERN napi_status 
napi_coerce_to_object(
    napi_env env,
    napi_value value,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_coerce_to_string(
    napi_env env,
    napi_value value,
    napi_value* result) {
        return napi_ok;
}

// Methods to work with Objects
NAPI_EXTERN napi_status 
napi_get_prototype(
    napi_env env,
    napi_value object,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_property_names(
    napi_env env,
    napi_value object,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_set_property(
    napi_env env,
    napi_value object,
    napi_value key,
    napi_value value) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_has_property(
    napi_env env,
    napi_value object,
    napi_value key,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_property(
    napi_env env,
    napi_value object,
    napi_value key,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_delete_property(
    napi_env env,
    napi_value object,
    napi_value key,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_has_own_property(
    napi_env env,
    napi_value object,
    napi_value key,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_set_named_property(
    napi_env env,
    napi_value object,
    const char* utf8name,
    napi_value value) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_has_named_property(
    napi_env env,
    napi_value object,
    const char* utf8name,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_named_property(
    napi_env env,
    napi_value object,
    const char* utf8name,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_set_element(
    napi_env env,
    napi_value object,
    uint32_t index,
    napi_value value) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_has_element(
    napi_env env,
    napi_value object,
    uint32_t index,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_element(
    napi_env env,
    napi_value object,
    uint32_t index,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_delete_element(
    napi_env env,
    napi_value object,
    uint32_t index,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status
napi_define_properties(
    napi_env env,
    napi_value object,
    size_t property_count,
    const napi_property_descriptor* properties) {
        return napi_ok;
}

// Methods to work with Arrays
NAPI_EXTERN napi_status 
napi_is_array(
    napi_env env,
    napi_value value,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_array_length(
    napi_env env,
    napi_value value,
    uint32_t* result) {
        return napi_ok;
}

// Methods to compare values
NAPI_EXTERN napi_status 
napi_strict_equals(
    napi_env env,
    napi_value lhs,
    napi_value rhs,
    bool* result) {
        return napi_ok;
}

// Methods to work with Functions
NAPI_EXTERN napi_status 
napi_call_function(
    napi_env env,
    napi_value recv,
    napi_value func,
    size_t argc,
    const napi_value* argv,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_new_instance(
    napi_env env,
    napi_value constructor,
    size_t argc,
    const napi_value* argv,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_instanceof(
    napi_env env,
    napi_value object,
    napi_value constructor,
    bool* result) {
        return napi_ok;
}

// Methods to work with napi_callbacks

// Gets all callback info in a single call. (Ugly, but faster.)
NAPI_EXTERN napi_status 
napi_get_cb_info(
    napi_env env,       
    napi_callback_info cbinfo,
    size_t* argc,  
    napi_value* argv,
    napi_value* this_arg,
    void** data) {
        return napi_ok;
} 

NAPI_EXTERN napi_status 
napi_get_new_target(
    napi_env env,
    napi_callback_info cbinfo,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status
napi_define_class(
    napi_env env,
    const char* utf8name,
    size_t length,
    napi_callback constructor,
    void* data,
    size_t property_count,
    const napi_property_descriptor* properties,
    napi_value* result) {
        return napi_ok;
}

// Methods to work with external data objects
NAPI_EXTERN napi_status 
napi_wrap(
    napi_env env,
    napi_value js_object,
    void* native_object,
    napi_finalize finalize_cb,
    void* finalize_hint,
    napi_ref* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_unwrap(
    napi_env env,
    napi_value js_object,
    void** result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_remove_wrap(
    napi_env env,
    napi_value js_object,
    void** result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_external(
    napi_env env,
    void* data,
    napi_finalize finalize_cb,
    void* finalize_hint,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_value_external(
    napi_env env,
    napi_value value,
    void** result) {
        return napi_ok;
}

// Set initial_refcount to 0 for a weak reference, >0 for a strong reference.
NAPI_EXTERN napi_status 
napi_create_reference(
    napi_env env,
    napi_value value,
    uint32_t initial_refcount,
    napi_ref* result) {
        return napi_ok;
}

// Deletes a reference. The referenced value is released, and may
// be GC'd unless there are other references to it.
NAPI_EXTERN napi_status 
napi_delete_reference(
    napi_env env, 
    napi_ref ref) {
        return napi_ok;
}

// Increments the reference count, optionally returning the resulting count.
// After this call the  reference will be a strong reference because its
// refcount is >0, and the referenced object is effectively "pinned".
// Calling this when the refcount is 0 and the object is unavailable
// results in an error.
NAPI_EXTERN napi_status 
napi_reference_ref(
    napi_env env,
    napi_ref ref,
    uint32_t* result) {
        return napi_ok;
}

// Decrements the reference count, optionally returning the resulting count.
// If the result is 0 the reference is now weak and the object may be GC'd
// at any time if there are no other references. Calling this when the
// refcount is already 0 results in an error.
NAPI_EXTERN napi_status 
napi_reference_unref(
    napi_env env,
    napi_ref ref,
    uint32_t* result) {
        return napi_ok;
}

// Attempts to get a referenced value. If the reference is weak,
// the value might no longer be available, in that case the call
// is still successful but the result is NULL.
NAPI_EXTERN napi_status 
napi_get_reference_value(
    napi_env env,
    napi_ref ref,
    napi_value* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_open_handle_scope(
    napi_env env,
    napi_handle_scope* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_close_handle_scope(
    napi_env env,
    napi_handle_scope scope) {
        return napi_ok;
}
NAPI_EXTERN napi_status
napi_open_escapable_handle_scope(
    napi_env env,
    napi_escapable_handle_scope* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status
napi_close_escapable_handle_scope(
    napi_env env,
    napi_escapable_handle_scope scope) {
        return napi_ok;
    }

NAPI_EXTERN napi_status 
napi_escape_handle(
    napi_env env,
    napi_escapable_handle_scope scope,
    napi_value escapee,
    napi_value* result) {
        return napi_ok;
}

// Methods to support error handling
NAPI_EXTERN napi_status 
napi_throw(
    napi_env env, 
    napi_value error) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_throw_error(
    napi_env env,
    const char* code,
    const char* msg) {
        return napi_ok;
    }
NAPI_EXTERN napi_status 
napi_throw_type_error(
    napi_env env,
    const char* code,
    const char* msg) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_throw_range_error(
    napi_env env,
    const char* code,
    const char* msg) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_is_error(
    napi_env env,
    napi_value value,
    bool* result) {
        return napi_ok;
}

// Methods to support catching exceptions
NAPI_EXTERN napi_status 
napi_is_exception_pending(
    napi_env env, 
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_and_clear_last_exception(
    napi_env env,
    napi_value* result) {
        return napi_ok;
}

// Methods to work with array buffers and typed arrays
NAPI_EXTERN napi_status 
napi_is_arraybuffer(
    napi_env env,
    napi_value value,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_arraybuffer(
    napi_env env,
    size_t byte_length,
    void** data,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status
napi_create_external_arraybuffer(
    napi_env env,
    void* external_data,
    size_t byte_length,
    napi_finalize finalize_cb,
    void* finalize_hint,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_arraybuffer_info(
    napi_env env,
    napi_value arraybuffer,
    void** data,
    size_t* byte_length) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_is_typedarray(
    napi_env env,
    napi_value value,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_typedarray(
    napi_env env,
    napi_typedarray_type type,
    size_t length,
    napi_value arraybuffer,
    size_t byte_offset,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_typedarray_info(
    napi_env env,
    napi_value typedarray,
    napi_typedarray_type* type,
    size_t* length,
    void** data,
    napi_value* arraybuffer,
    size_t* byte_offset) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_create_dataview(
    napi_env env,
    size_t length,
    napi_value arraybuffer,
    size_t byte_offset,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_is_dataview(
    napi_env env,
    napi_value value,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_dataview_info(
    napi_env env,
    napi_value dataview,
    size_t* bytelength,
    void** data,
    napi_value* arraybuffer,
    size_t* byte_offset) {
        return napi_ok;
}

// version management
NAPI_EXTERN napi_status 
napi_get_version(
    napi_env env, 
    uint32_t* result) {
        return napi_ok;
}

// Promises
NAPI_EXTERN napi_status 
napi_create_promise(
    napi_env env,
    napi_deferred* deferred,
    napi_value* promise) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_resolve_deferred(
    napi_env env,
    napi_deferred deferred,
    napi_value resolution) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_reject_deferred(
    napi_env env,
    napi_deferred deferred,
    napi_value rejection) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_is_promise(
    napi_env env,
    napi_value promise,
    bool* is_promise) {
        return napi_ok;
}

// Running a script
NAPI_EXTERN napi_status 
napi_run_script(
    napi_env env,
    napi_value script,
    napi_value* result) {
        return napi_ok;
}

// Memory management
NAPI_EXTERN napi_status 
napi_adjust_external_memory(
    napi_env env,
    int64_t change_in_bytes,
    int64_t* adjusted_value) {
        return napi_ok;
}

#ifdef NAPI_EXPERIMENTAL

// Dates
NAPI_EXTERN napi_status 
napi_create_date(
    napi_env env,
    double time,
    napi_value* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_is_date(
    napi_env env,
    napi_value value,
    bool* is_date) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_get_date_value(
    napi_env env,
    napi_value value,
    double* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_create_bigint_int64(
    napi_env env,
    int64_t value,
    napi_value* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_create_bigint_uint64(
    napi_env env,
    uint64_t value,
    napi_value* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_create_bigint_words(
    napi_env env,
    int sign_bit,
    size_t word_count,
    const uint64_t* words,
    napi_value* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_get_value_bigint_int64(
    napi_env env,
    napi_value value,
    int64_t* result,
    bool* lossless) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_get_value_bigint_uint64(
    napi_env env,
    napi_value value,
    uint64_t* result,
    bool* lossless) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_get_value_bigint_words(
    napi_env env,
    napi_value value,
    int* sign_bit,
    size_t* word_count,
    uint64_t* words) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_add_finalizer(
    napi_env env,
    napi_value js_object,
    void* native_object,
    napi_finalize finalize_cb,
    void* finalize_hint,
    napi_ref* result) {
        return napi_ok;
}
#endif  // NAPI_EXPERIMENTAL

// ##############################################################################
// Stub for runtime API
// ##############################################################################

NAPI_EXTERN NAPI_NO_RETURN 
void napi_fatal_error(
    const char* location,
    size_t location_len,
    const char* message,
    size_t message_len);

// Methods for custom handling of async operations
NAPI_EXTERN napi_status 
napi_async_init(
    napi_env env,
    napi_value async_resource,
    napi_value async_resource_name,
    napi_async_context* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_async_destroy(
    napi_env env,
    napi_async_context async_context) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_make_callback(
    napi_env env,
    napi_async_context async_context,
    napi_value recv,
    napi_value func,
    size_t argc,
    const napi_value* argv,
    napi_value* result) {
        return napi_ok;
}

// Methods to provide node::Buffer functionality with napi types
NAPI_EXTERN napi_status 
napi_create_buffer(
    napi_env env,
    size_t length,
    void** data,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_external_buffer(
    napi_env env,
    size_t length,
    void* data,
    napi_finalize finalize_cb,
    void* finalize_hint,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_create_buffer_copy(
    napi_env env,
    size_t length,
    const void* data,
    void** result_data,
    napi_value* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_is_buffer(
    napi_env env,
    napi_value value,
    bool* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_get_buffer_info(
    napi_env env,
    napi_value value,
    void** data,
    size_t* length) {
        return napi_ok;
}

// Methods to manage simple async operations
NAPI_EXTERN napi_status 
napi_create_async_work(
    napi_env env,
    napi_value async_resource,
    napi_value async_resource_name,
    napi_async_execute_callback execute,
    napi_async_complete_callback complete,
    void* data,
    napi_async_work* result) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_delete_async_work(
    napi_env env,
    napi_async_work work) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_queue_async_work(
    napi_env env,
    napi_async_work work) {
        return napi_ok;
}
NAPI_EXTERN napi_status 
napi_cancel_async_work(
    napi_env env,
    napi_async_work work) {
        return napi_ok;
}

#if NAPI_VERSION >= 2

// Return the current libuv event loop for a given environment
NAPI_EXTERN napi_status 
napi_get_uv_event_loop(
    napi_env env,
    struct uv_loop_s** loop) {
        return napi_ok;
}

#endif  // NAPI_VERSION >= 2

#if NAPI_VERSION >= 3

NAPI_EXTERN napi_status 
napi_fatal_exception(napi_env env, napi_value err);

NAPI_EXTERN napi_status 
napi_add_env_cleanup_hook(
    napi_env env,
    void (*fun)(void* arg),
    void* arg) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_remove_env_cleanup_hook(
    napi_env env,
    void (*fun)(void* arg),
    void* arg) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_open_callback_scope(
    napi_env env,
    napi_value resource_object,
    napi_async_context context,
    napi_callback_scope* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status 
napi_close_callback_scope(
    napi_env env,
    napi_callback_scope scope) {
        return napi_ok;
}

#endif  // NAPI_VERSION >= 3

#if NAPI_VERSION >= 4

// Calling into JS from other threads
NAPI_EXTERN napi_status
napi_create_threadsafe_function(
    napi_env env,
    napi_value func,
    napi_value async_resource,
    napi_value async_resource_name,
    size_t max_queue_size,
    size_t initial_thread_count,
    void* thread_finalize_data,
    napi_finalize thread_finalize_cb,
    void* context,
    napi_threadsafe_function_call_js call_js_cb,
    napi_threadsafe_function* result) {
        return napi_ok;
}

NAPI_EXTERN napi_status
napi_get_threadsafe_function_context(
    napi_threadsafe_function func,
    void** result) {
        return napi_ok;
}

NAPI_EXTERN napi_status
napi_call_threadsafe_function(
    napi_threadsafe_function func,
    void* data,
    napi_threadsafe_function_call_mode is_blocking) {
        return napi_ok;
}

NAPI_EXTERN napi_status
napi_acquire_threadsafe_function(
    napi_threadsafe_function func) {
        return napi_ok;
}

NAPI_EXTERN napi_status
napi_release_threadsafe_function(
    napi_threadsafe_function func,
    napi_threadsafe_function_release_mode mode) {
        return napi_ok;
}

NAPI_EXTERN napi_status
napi_unref_threadsafe_function(
    napi_env env, 
    napi_threadsafe_function func) {
        return napi_ok;
}

NAPI_EXTERN napi_status
napi_ref_threadsafe_function(
    napi_env env, 
    napi_threadsafe_function func) {
        return napi_ok;
}

#endif  // NAPI_VERSION >= 4
