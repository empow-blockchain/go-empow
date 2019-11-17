#include "storage.h"
#include <stdio.h>
#include <string.h>
#include <iostream>

static putFunc CPut = nullptr;
static hasFunc CHas = nullptr;
static getFunc CGet = nullptr;
static delFunc CDel = nullptr;
static mapPutFunc CMapPut = nullptr;
static mapHasFunc CMapHas = nullptr;
static mapGetFunc CMapGet = nullptr;
static mapDelFunc CMapDel = nullptr;
static mapKeysFunc CMapKeys = nullptr;
static mapLenFunc CMapLen = nullptr;

static globalHasFunc CGHas = nullptr;
static globalGetFunc CGGet = nullptr;
static globalMapHasFunc CGMapHas = nullptr;
static globalMapGetFunc CGMapGet = nullptr;
static globalMapKeysFunc CGMapKeys = nullptr;
static globalMapLenFunc CGMapLen = nullptr;

void InitGoStorage(putFunc put, hasFunc has, getFunc get, delFunc del,
    mapPutFunc mput, mapHasFunc mhas, mapGetFunc mget, mapDelFunc mdel, mapKeysFunc mkeys, mapLenFunc mlen,
    globalHasFunc ghas, globalGetFunc gget, globalMapHasFunc gmhas, globalMapGetFunc gmget, globalMapKeysFunc gmkeys, globalMapLenFunc gmlen) {

    CPut = put;
    CHas = has;
    CGet = get;
    CDel = del;
    CMapPut = mput;
    CMapHas = mhas;
    CMapGet = mget;
    CMapDel = mdel;
    CMapKeys = mkeys;
    CMapLen = mlen;
    CGHas = ghas;
    CGGet = gget;
    CGMapHas = gmhas;
    CGMapGet = gmget;
    CGMapKeys = gmkeys;
    CGMapLen = gmlen;
}

char* EMPOWContractStorage::Put(const CStr key, const CStr value, const CStr ramPayer) {
    size_t gasUsed = 0;
    char *ret = CPut(sbxPtr, key, value, ramPayer, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::Has(const CStr key, const CStr ramPayer, bool *result) {
    size_t gasUsed = 0;
    char *ret = CHas(sbxPtr, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::Get(const CStr key, const CStr ramPayer, CStr *result) {
    size_t gasUsed = 0;
    char *ret = CGet(sbxPtr, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::Del(const CStr key, const CStr ramPayer) {
    size_t gasUsed = 0;
    char *ret = CDel(sbxPtr, key, ramPayer, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::MapPut(const CStr key, const CStr field, const CStr value, const CStr ramPayer) {
    size_t gasUsed = 0;
    char *ret = CMapPut(sbxPtr, key, field, value, ramPayer, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;

}

char* EMPOWContractStorage::MapHas(const CStr key, const CStr field, const CStr ramPayer, bool *result) {
    size_t gasUsed = 0;
    char *ret = CMapHas(sbxPtr, key, field, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::MapGet(const CStr key, const CStr field, const CStr ramPayer, CStr *result) {
    size_t gasUsed = 0;
    char *ret = CMapGet(sbxPtr, key, field, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::MapDel(const CStr key, const CStr field, const CStr ramPayer) {
    size_t gasUsed = 0;
    char *ret = CMapDel(sbxPtr, key, field, ramPayer, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::MapKeys(const CStr key, const CStr ramPayer, CStr *result) {
    size_t gasUsed = 0;
    char *ret = CMapKeys(sbxPtr, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::MapLen(const CStr key, const CStr ramPayer, size_t *result) {
    size_t gasUsed = 0;
    char *ret = CMapLen(sbxPtr, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::GlobalHas(const CStr contract, const CStr key, const CStr ramPayer, bool *result) {
    size_t gasUsed = 0;
    char *ret = CGHas(sbxPtr, contract, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::GlobalGet(const CStr contract, const CStr key, const CStr ramPayer, CStr *result) {
    size_t gasUsed = 0;
    char *ret = CGGet(sbxPtr, contract, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::GlobalMapHas(const CStr contract, const CStr key, const CStr field, const CStr ramPayer, bool *result) {
    size_t gasUsed = 0;
    char *ret = CGMapHas(sbxPtr, contract, key, field, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::GlobalMapGet(const CStr contract, const CStr key, const CStr field, const CStr ramPayer, CStr *result) {
    size_t gasUsed = 0;
    char *ret = CGMapGet(sbxPtr, contract, key, field, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::GlobalMapKeys(const CStr contract,  const CStr key, const CStr ramPayer, CStr *result) {
    size_t gasUsed = 0;
    char *ret = CGMapKeys(sbxPtr, contract, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

char* EMPOWContractStorage::GlobalMapLen(const CStr contract, const CStr key, const CStr ramPayer, size_t *result) {
    size_t gasUsed = 0;
    char *ret = CGMapLen(sbxPtr, contract, key, ramPayer, result, &gasUsed);
    Sandbox *sbx = static_cast<Sandbox*>(sbxPtr);
    sbx->gasUsed += gasUsed;
    return ret;
}

void NewEMPOWContractStorage(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Context> context = isolate->GetCurrentContext();
    Local<Object> global = context->Global();

    Local<Value> val = global->GetInternalField(0);
    if (!val->IsExternal()) {
        std::cout << "NewEMPOWContractStorage val error" << std::endl;
        return;
    }
    SandboxPtr sbx = static_cast<SandboxPtr>(Local<External>::Cast(val)->Value());

    EMPOWContractStorage *ics = new EMPOWContractStorage(sbx);

    Local<Object> self = args.Holder();
    self->SetInternalField(0, External::New(isolate, ics));

    args.GetReturnValue().Set(self);
}

void EMPOWContractStorage_Put(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Put invalid argument length.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Put key must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> val = args[1];
    if (!val->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Put value must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[2];
    if (!val->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Put ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(valStr, val, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_Put val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char* ret = ics->Put(keyStr, valStr, ramPayerStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().SetNull();
}

void EMPOWContractStorage_Has(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 2) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Has invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Has key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[1];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Has ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    bool result;

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_Has val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->Has(keyStr, ramPayerStr, &result);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set(result);
}

void EMPOWContractStorage_Get(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 2) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Get invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Get key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[1];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Get ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    CStr resultStr = {nullptr, 0};

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_Get val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->Get(keyStr, ramPayerStr, &resultStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }

    if (resultStr.data == nullptr) {
        args.GetReturnValue().SetNull();
    } else {
        args.GetReturnValue().Set(String::NewFromUtf8(isolate, resultStr.data, String::kNormalString, resultStr.size));
        free(resultStr.data);
    }
}

void EMPOWContractStorage_Del(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 2) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Del invalid argument length.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Del key must be string.")
        );
        isolate->ThrowException(err);
        return;
    }
    Local<Value> ramPayer = args[1];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_Del ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_Del val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->Del(keyStr, ramPayerStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().SetNull();
}

void EMPOWContractStorage_MapPut(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 4) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapPut invalid argument length.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapPut key must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> field = args[1];
    if (!field->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapPut key must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> val = args[2];
    if (!val->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapPut value must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[3];
    if (!val->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapPut ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(fieldStr, field, isolate);
    NewCStrChecked(valStr, val, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_MapPut val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->MapPut(keyStr, fieldStr, valStr, ramPayerStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().SetNull();
}

void EMPOWContractStorage_MapHas(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapHas invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapHas key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> field = args[1];
    if (!field->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapHas key must be string")
        );
        isolate->ThrowException(err);
        return;
    }
    Local<Value> ramPayer = args[2];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapHas ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(fieldStr, field, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    bool result;

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_MapHas val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->MapHas(keyStr, fieldStr, ramPayerStr, &result);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set(result);
}

void EMPOWContractStorage_MapGet(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapGet invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapGet key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> field = args[1];
    if (!field->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapGet key must be string")
        );
        isolate->ThrowException(err);
        return;
    }
    Local<Value> ramPayer = args[2];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapGet ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(fieldStr, field, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    CStr resultStr = {nullptr, 0};

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_MapGet val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->MapGet(keyStr, fieldStr, ramPayerStr, &resultStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    if (resultStr.data == nullptr) {
        args.GetReturnValue().SetNull();
    } else {
        args.GetReturnValue().Set(String::NewFromUtf8(isolate, resultStr.data, String::kNormalString, resultStr.size));
        free(resultStr.data);
    }
}

void EMPOWContractStorage_MapDel(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapDel invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapDel key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> field = args[1];
    if (!field->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapDel key must be string")
        );
        isolate->ThrowException(err);
        return;
    }
    Local<Value> ramPayer = args[2];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapDel ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(fieldStr, field, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_MapDel val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->MapDel(keyStr, fieldStr, ramPayerStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().SetNull();
}

void EMPOWContractStorage_MapKeys(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 2) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapKeys invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapKeys key must be string")
        );
        isolate->ThrowException(err);
        return;
    }
    Local<Value> ramPayer = args[1];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapKeys ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    CStr resultStr = {nullptr, 0};

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_MapKeys val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->MapKeys(keyStr, ramPayerStr, &resultStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set(String::NewFromUtf8(isolate, resultStr.data, String::kNormalString, resultStr.size));
    if (resultStr.data != nullptr) free(resultStr.data);
}

void EMPOWContractStorage_MapLen(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 2) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapLen invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[0];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapLen key must be string")
        );
        isolate->ThrowException(err);
        return;
    }
    Local<Value> ramPayer = args[1];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_MapLen ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    size_t result;

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_MapLen val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->MapLen(keyStr, ramPayerStr, &result);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set((int)result);
}

void EMPOWContractStorage_GlobalHas(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalHas invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> contract = args[0];
    if (!contract->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalHas contract must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[1];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalHas key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[2];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalHas ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(contractStr, contract, isolate);
    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    bool result;

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_GlobalHas val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->GlobalHas(contractStr, keyStr, ramPayerStr, &result);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set(result);
}

void EMPOWContractStorage_GlobalGet(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalGet invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> contract = args[0];
    if (!contract->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalGet contract must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[1];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalGet key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[2];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalGet ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(contractStr, contract, isolate);
    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    CStr resultStr = {nullptr, 0};

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_GlobalGet val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->GlobalGet(contractStr, keyStr, ramPayerStr, &resultStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    if (resultStr.data == nullptr) {
        args.GetReturnValue().SetNull();
    } else {
        args.GetReturnValue().Set(String::NewFromUtf8(isolate, resultStr.data, String::kNormalString, resultStr.size));
        free(resultStr.data);
    }
}

void EMPOWContractStorage_GlobalMapHas(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 4) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapHas invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> contract = args[0];
    if (!contract->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapHas contract must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[1];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapHas key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> field = args[2];
    if (!field->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapHas field must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[3];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapHas ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(contractStr, contract, isolate);
    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(fieldStr, field, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    bool result;

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_GlobalMapHas val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->GlobalMapHas(contractStr, keyStr, fieldStr, ramPayerStr, &result);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set(result);
}

void EMPOWContractStorage_GlobalMapGet(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 4) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapGet invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> contract = args[0];
    if (!contract->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapGet contract must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[1];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapGet key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> field = args[2];
    if (!field->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapGet field must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[3];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapGet ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(contractStr, contract, isolate);
    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(fieldStr, field, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    CStr resultStr = {nullptr, 0};

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_GlobalMapGet val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->GlobalMapGet(contractStr, keyStr, fieldStr, ramPayerStr, &resultStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    if (resultStr.data == nullptr) {
        args.GetReturnValue().SetNull();
    } else {
        args.GetReturnValue().Set(String::NewFromUtf8(isolate, resultStr.data, String::kNormalString, resultStr.size));
        free(resultStr.data);
    }
}

void EMPOWContractStorage_GlobalMapKeys(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapKeys invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> contract = args[0];
    if (!contract->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapKeys contract must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[1];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapKeys key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[2];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapKeys ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(contractStr, contract, isolate);
    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    CStr resultStr = {nullptr, 0};

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_GlobalMapKeys val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->GlobalMapKeys(contractStr, keyStr, ramPayerStr, &resultStr);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set(String::NewFromUtf8(isolate, resultStr.data, String::kNormalString, resultStr.size));
    if (resultStr.data != nullptr) free(resultStr.data);
}

void EMPOWContractStorage_GlobalMapLen(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Object> self = args.Holder();

    if (args.Length() != 3) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapLen invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> contract = args[0];
    if (!contract->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapLen contract must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> key = args[1];
    if (!key->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapLen key must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> ramPayer = args[2];
    if (!ramPayer->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "EMPOWContractStorage_GlobalMapLen ramPayer must be string.")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(contractStr, contract, isolate);
    NewCStrChecked(keyStr, key, isolate);
    NewCStrChecked(ramPayerStr, ramPayer, isolate);
    size_t result;

    Local<External> extVal = Local<External>::Cast(self->GetInternalField(0));
    if (!extVal->IsExternal()) {
        std::cout << "EMPOWContractStorage_GlobalMapLen val error" << std::endl;
        return;
    }

    EMPOWContractStorage *ics = static_cast<EMPOWContractStorage *>(extVal->Value());
    char *ret = ics->GlobalMapLen(contractStr, keyStr, ramPayerStr, &result);
    if (ret != nullptr) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, ret)
        );
        isolate->ThrowException(err);
        free(ret);
        return;
    }
    args.GetReturnValue().Set((int)result);
}

void InitStorage(Isolate *isolate, Local<ObjectTemplate> globalTpl) {
    Local<FunctionTemplate> storageClass =
        FunctionTemplate::New(isolate, NewEMPOWContractStorage);
    Local<String> storageClassName = String::NewFromUtf8(isolate, "EMPOWStorage");
    storageClass->SetClassName(storageClassName);

    Local<ObjectTemplate> storageTpl = storageClass->InstanceTemplate();
    storageTpl->SetInternalFieldCount(1);
    storageTpl->Set(
        String::NewFromUtf8(isolate, "put"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_Put)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "has"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_Has)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "get"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_Get)
    );
    storageTpl->Set(
            String::NewFromUtf8(isolate, "del"),
            FunctionTemplate::New(isolate, EMPOWContractStorage_Del)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "mapPut"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_MapPut)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "mapHas"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_MapHas)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "mapGet"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_MapGet)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "mapDel"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_MapDel)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "mapKeys"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_MapKeys)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "mapLen"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_MapLen)
    );
    // todo
    storageTpl->Set(
        String::NewFromUtf8(isolate, "globalGet"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_GlobalGet)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "globalHas"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_GlobalHas)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "globalMapHas"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_GlobalMapHas)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "globalMapGet"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_GlobalMapGet)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "globalMapKeys"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_GlobalMapKeys)
    );
    storageTpl->Set(
        String::NewFromUtf8(isolate, "globalMapLen"),
        FunctionTemplate::New(isolate, EMPOWContractStorage_GlobalMapLen)
    );


    globalTpl->Set(storageClassName, storageClass);
}
