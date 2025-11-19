#include <stdint.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

__attribute__((visibility("default"))) uint64_t alloc(uint64_t size) {
    void* p = malloc(size);
    return (uint64_t)p;
}

__attribute__((visibility("default"))) uint64_t process(uint64_t inPtr, uint64_t inLen, uint64_t outPtr, uint64_t outCap) {
    char *input = (char *)inPtr;

    size_t len = (size_t)inLen;

    char* msg = "{\"action_type\":\"none\",\"payload\":{\"messge\":\"hello wasm\"}}";
    char* out = (char*)outPtr;
    memcpy(out, msg, strlen(msg));
    return strlen(msg);
}

int main() {
    return 0;
}