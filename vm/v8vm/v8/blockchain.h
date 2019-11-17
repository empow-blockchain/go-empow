#ifndef EMPOW_V8_BLOCKCHAIN_H
#define EMPOW_V8_BLOCKCHAIN_H

#include "sandbox.h"
#include "stddef.h"

using namespace v8;

void InitBlockchain(Isolate *isolate, Local<ObjectTemplate> globalTpl);
void NewEMPOWBlockchain(const FunctionCallbackInfo<Value> &args);

// This Class wraps Go BlockChain function so JS contract can call them.
class EMPOWBlockchain {
private:
    SandboxPtr sbxPtr;
public:
    EMPOWBlockchain(SandboxPtr ptr): sbxPtr(ptr) {}

    char* BlockInfo(CStr *result);
    char* TxInfo(CStr *result);
    char* ContextInfo(CStr *result);
    char* Call(const CStr contract, const CStr api, const CStr args, CStr *result);
    char* CallWithAuth(const CStr contract, const CStr api, const CStr args, CStr *result);
    char* RequireAuth(const CStr accountID, const CStr permission, bool *result);
    char* Receipt(const CStr content);
    char* Event(const CStr content);
};

#endif // EMPOW_V8_BLOCKCHAIN_H
