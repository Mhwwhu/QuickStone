package wasmrunner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmRuntime struct {
	Context    context.Context
	Runtime    wazero.Runtime
	Module     api.Module
	Memory     api.Memory
	WasiModule api.Closer
	// modName string
}

type WasmResult struct {
	Values []uint64
	Err    error
}

func NewWasmRuntime(ctx context.Context, wasmBytes []byte) (*WasmRuntime, error) {
	// 1. 创建 runtime
	rt := wazero.NewRuntime(ctx)

	// 2. 编译模块
	compiled, err := rt.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}
	defer compiled.Close(ctx)

	wasiModule, err := wasi_snapshot_preview1.Instantiate(ctx, rt)
	if err != nil {
		panic(err) // 这一步必须成功
	}

	// 3. 实例化模块（注意：很多 wasm 需要配置 WASI 等，这里简化）
	mod, err := rt.InstantiateModule(ctx, compiled, wazero.NewModuleConfig())
	if err != nil {
		return nil, err
	}

	// 4. 拿到导出的 memory 对象，一般名字就是 "memory"
	mem := mod.ExportedMemory("memory")
	if mem == nil {
		return nil, fmt.Errorf("module has no exported memory \"memory\"")
	}

	// 5. 塞到自定义的 runtime 结构里
	return &WasmRuntime{
		Context:    ctx,
		Runtime:    rt,
		Module:     mod,
		Memory:     mem,
		WasiModule: wasiModule,
	}, nil
}

func (r *WasmRuntime) Close(ctx context.Context) error {
	var err1, err2 error
	if r.Module != nil {
		err1 = r.Module.Close(ctx)
	}
	if r.Runtime != nil {
		err2 = r.Runtime.Close(ctx)
	}
	if err1 != nil {
		return err1
	}
	return err2
}

// ABI: funcName(inPtr i32, inLen i32, outPtr i32, outCap i32) -> i32(outLen)
func (r *WasmRuntime) Call(
	ctx context.Context,
	funcName string,
	in any,
) (json.RawMessage, error) {
	// 1. 输入序列化
	inBytes, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal input: %w", err)
	}

	mod := r.Module
	mem := r.Memory

	// 3. 拿到 alloc + 目标函数
	allocFn := mod.ExportedFunction("alloc")
	if allocFn == nil {
		return nil, fmt.Errorf("alloc function not found")
	}

	handleFn := mod.ExportedFunction(funcName)
	if handleFn == nil {
		return nil, fmt.Errorf("function %q not found", funcName)
	}

	// 4. 给输入申请内存
	inLen := uint32(len(inBytes))
	inPtr, err := callAlloc(ctx, allocFn, inLen)
	if err != nil {
		return nil, fmt.Errorf("alloc input: %w", err)
	}
	if ok := mem.Write(inPtr, inBytes); !ok {
		return nil, fmt.Errorf("write input to wasm memory failed")
	}

	// 5. 给输出申请内存
	// 简单策略：输出缓冲区容量 = max(输入长度 * 4, 1024)
	outCap := inLen * 4
	if outCap < 1024 {
		outCap = 1024
	}

	outPtr, err := callAlloc(ctx, allocFn, outCap)
	if err != nil {
		return nil, fmt.Errorf("alloc output: %w", err)
	}

	// 6. 调用 handle(inPtr, inLen, outPtr, outCap) -> outLen
	res, err := handleFn.Call(
		ctx,
		uint64(inPtr),
		uint64(inLen),
		uint64(outPtr),
		uint64(outCap),
	)
	if err != nil || len(res) == 0 {
		return nil, fmt.Errorf("call %s: %w", funcName, err)
	}

	outLen := uint32(res[0]) // i32 返回值放在低 32 位

	if outLen > outCap {
		return nil, fmt.Errorf("wasm returned length %d > output capacity %d", outLen, outCap)
	}

	// 7. 从 wasm memory 读出输出 JSON
	outBytes, ok := mem.Read(outPtr, outLen)
	if !ok {
		return nil, fmt.Errorf("read output from wasm memory failed")
	}

	// 不在这里 free，依赖实例销毁一起回收
	return json.RawMessage(outBytes), nil
}

// 小工具：从 alloc 调用里拿到 ptr
func callAlloc(ctx context.Context, allocFn api.Function, size uint32) (uint32, error) {
	res, err := allocFn.Call(ctx, uint64(size))
	if err != nil || len(res) == 0 {
		return 0, err
	}
	return uint32(res[0]), nil
}
