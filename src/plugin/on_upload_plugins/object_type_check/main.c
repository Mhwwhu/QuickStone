#include <stdint.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>

// 模拟堆的静态分配器
extern unsigned char __heap_base;
uintptr_t heap_start = (uintptr_t)&__heap_base;
static uint32_t current_pages;
static uint64_t heap_ptr;

// 内存分配函数：分配指定字节数的内存
__attribute__((visibility("default"))) uint64_t alloc(uint64_t size) {
    uint64_t allocated_ptr = heap_ptr;
    heap_ptr += size;  // 向后移动 heap_ptr，模拟分配内存
    return allocated_ptr;  // 返回内存块的起始地址
}

// echo 函数，接收输入字符串的指针和长度，返回复述的字符串
__attribute__((visibility("default"))) uint64_t echo(uint64_t inPtr, uint64_t inLen, uint64_t outPtr, uint64_t outCap) {
    // 获取输入字符串
    char *input = (char *)inPtr;

    // 计算输入字符串的长度
    size_t len = (size_t)inLen;

    // 如果输出缓冲区大小不够，直接返回 0
    if (len > outCap) {
        return 0;
    }

    // 将输入字符串复制到输出缓冲区
    char *output = (char *)outPtr;
    // for (size_t i = 0; i < len; i++) {
    //     output[i] = input[i];
    // }
    char msg[50];
    sprintf(msg, "\"heap_start:%lu, heap_ptr:%llu, current_pages:%lu\"", heap_start, heap_ptr, __builtin_wasm_memory_size(0));
    for(size_t i = 0; i <  strlen(msg); i++) {
        output[i] = msg[i];
    }

    // 返回写入的字节数
    return strlen(msg);
}

int main() {
    current_pages =  __builtin_wasm_memory_size(0); 
    heap_ptr = heap_start;
    return 0;
}