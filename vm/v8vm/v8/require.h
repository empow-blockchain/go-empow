#ifndef EMPOW_V8_REQUIRE_H
#define EMPOW_V8_REQUIRE_H

#include "sandbox.h"

//extern char *requireModule(SandboxPtr, const char *);
void InitRequire(Isolate *isolate, Local<ObjectTemplate> globalTpl);
void NewNativeRequire(const FunctionCallbackInfo<Value> &info);

#endif // EMPOW_V8_REQUIRE_H