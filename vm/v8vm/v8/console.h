#ifndef EMPOW_V8_CONSOLE_H
#define EMPOW_V8_CONSOLE_H

#include "sandbox.h"

// This Class Provide Console.Log Function so JS code can use Go log.
void InitConsole(Isolate *isolate, Local<ObjectTemplate> globalTpl);
void NewConsoleLog(const FunctionCallbackInfo<Value> &args);

#endif // EMPOW_V8_CONSOLE_H