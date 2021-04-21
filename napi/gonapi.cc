#include "gonapi.h"
#include <stdio.h>
#include <utility>
#include <vector>
#include <cassert>

#include "_cgo_export.h"


napi_value MethodP(napi_env env, napi_callback_info info) {
    napi_status status;
    napi_value world;
    status = napi_create_string_utf8(env, "world", 5, &world);
    assert(status == napi_ok);
        return world;
}

struct Context{
    void* data;
    Context(void* data) : data(data) {}
};

static std::vector<Context*> m;

static std::vector<void*> registry{};
struct AsyncExecuteCallbackWrap {
  AsyncExecuteCallbackWrap(void* data) : data{data} {}
  napi_async_execute_callback operator()() {
    static auto dataCopy = data;
    return [](napi_env env, void* data) -> void {
        return CallAsyncExecuteCallback(dataCopy, env, data);
    };
  }
  void* data;
};

struct AsyncCompleteCallbackWrap {
  AsyncCompleteCallbackWrap(void* data) : data{data} {}
  napi_async_complete_callback operator()() {
    static auto dataCopy = data;
    return [](napi_env env, napi_status status, void* data) -> void {
        return CallAsyncCompleteCallback(dataCopy, env, status, data);
    };
  }
  void* data;
};

struct FinalizeCallbackWrap {
  FinalizeCallbackWrap(void* data) : data{data} {}
  napi_finalize operator()() {
    static auto dataCopy = data;
    return [](napi_env env, void* data, void* hint) -> void {
        return CallFinalizeCallback(dataCopy, env, data, hint);
    };
  }
  void* data;
};

struct ThreadsafeFunctionCallbackWrap {
  ThreadsafeFunctionCallbackWrap(void* data) : data{data} {}
  napi_threadsafe_function_call_js operator()() {
    static auto dataCopy = data;
    return [](napi_env env, napi_value callback, void* ctx, void* data) -> void {
        return CallThreadsafeFunctionCallback(dataCopy, env, callback, ctx, data);
    };
  }
  void* data;
};

struct CallbackWrap {
  CallbackWrap(void* data) : data{data} {}
  static inline
  napi_value Wrapper(napi_env env, napi_callback_info info) {
    printf("Registry %p\n", registry[0]);
    return CallCallback(registry[0], env, info);
  }
  napi_callback operator()() {
    return [](napi_env env, napi_callback_info info) -> napi_value {
        return CallCallback(registry[0], env, info);
    };
  }
  void* data;
};



#ifdef __cplusplus
extern "C" {
#endif

napi_value InvokeGoHandler(napi_env env, napi_callback_info info) {
  return GoHandler(env, info);
}

/*napi_callback Callback(void* caller) {
  printf("Callback called\n");
  registry.push_back(caller);
  printf("Registry %p\n", registry[0]);
  return [](napi_env env, napi_callback_info info) -> napi_value {
        return CallCallback(registry[0], env, info);
    };
}*/

napi_async_execute_callback AsyncExecuteCallback(void* caller) {
  AsyncExecuteCallbackWrap cb{caller};
  return cb();

}

napi_async_complete_callback AsyncCompleteCallback(void* caller) {
  AsyncCompleteCallbackWrap cb{caller};
  return cb();
}

napi_finalize FinalizeCallback(void* caller) {
  FinalizeCallbackWrap cb{caller};
  return cb();
}

napi_threadsafe_function_call_js ThreadsafeFunctionCallback(void* caller) {
  ThreadsafeFunctionCallbackWrap cb{caller};
  return cb();
}


#ifdef __cplusplus
}  // extern "C"
#endif