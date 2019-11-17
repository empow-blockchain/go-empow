#ifndef EMPOW_V8_INSTRUCTION_H
#define EMPOW_V8_INSTRUCTION_H

#include "sandbox.h"
#include "stddef.h"
#include <iostream>
#include <cstring>
#include <string>

void InitInstruction(Isolate *isolate, Local<ObjectTemplate> globalTpl);
void NewEMPOWContractInstruction(const FunctionCallbackInfo<Value> &info);
void EMPOWContractInstruction_Count(const FunctionCallbackInfo<Value> &args);
void EMPOWContractInstruction_Incr(const FunctionCallbackInfo<Value> &args);

class EMPOWContractInstruction {
private:
    Sandbox* sbxPtr;
    Isolate* isolate;
    int count;
public:
    EMPOWContractInstruction(SandboxPtr ptr){
        sbxPtr = static_cast<Sandbox*>(ptr);
        isolate = sbxPtr->isolate;
        count = 0;
    }

    size_t Incr(size_t num) {
        if (sbxPtr->gasUsed > SIZE_MAX - num) {
            Local<Value> err = Exception::Error(
                String::NewFromUtf8(isolate, "EMPOWContractInstruction_Incr gas overflow size_t")
            );
            isolate->ThrowException(err);
            return 0;
        }

        sbxPtr->gasUsed += num;
        count ++;
        return sbxPtr->gasUsed;
    }
    size_t Count() {
        return count;
    }
    void MemUsageCheck(){
        size_t usedMem = MemoryUsage(isolate, sbxPtr->allocator);
        if (usedMem > sbxPtr->memLimit){
            Local<Value> err = Exception::Error(
                String::NewFromUtf8(isolate, ("EMPOWContractInstruction_Incr Memory Using too much! used: " + std::to_string(usedMem) + " Limit: " + std::to_string(sbxPtr->memLimit)).c_str())
            );
            isolate->ThrowException(err);
        }
        return;
    }
};

#endif // EMPOW_V8_INSTRUCTION_H