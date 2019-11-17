#ifndef EMPOW_V8_CRYPTO_H
#define EMPOW_V8_CRYPTO_H

#include "sandbox.h"

// This Class Provide Console.Log Function so JS code can use Go log.
void InitCrypto(Isolate *isolate, Local<ObjectTemplate> globalTpl);
void NewCrypto(const FunctionCallbackInfo<Value> &info);

class EMPOWCrypto {
private:
    SandboxPtr sbxPtr;
public:
    EMPOWCrypto(SandboxPtr ptr): sbxPtr(ptr) {}

    CStr sha3(const CStr msg);
    int verify(const CStr algo, const CStr msg, const CStr sig, const CStr pubkey);
};

#endif // EMPOW_V8_CRYPTO_H