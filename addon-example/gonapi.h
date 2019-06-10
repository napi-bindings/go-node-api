#ifndef GO_NAPI_H
#define GO_NAPI_H

#include <node_api.h>

#ifdef __cplusplus
extern "C" {
#endif

extern napi_callback CallbackMethod(void* caller);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // GO_NAPI_H