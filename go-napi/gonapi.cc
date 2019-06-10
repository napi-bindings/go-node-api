#include "gonapi.h"

#include "_cgo_export.h"

struct CallbackData {
  CallbackData(void* data) : data{data} {}
  napi_callback operator()() {
    static auto dataCopy = data;
    return [](napi_env env, napi_callback_info info) -> napi_value {
        return ExecuteCallback(dataCopy, env, info);
    };
  }
  void* data;
};

napi_callback CallbackMethod(void* caller) {
    CallbackData cb{caller};
    return cb();
}