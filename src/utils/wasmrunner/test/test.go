package main

import (
	onuploadpluginmodels "QuickStone/src/models/pluginModels/onUploadPluginModels"
	"QuickStone/src/utils/wasmrunner"
	"context"
	"encoding/json"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	file, err := os.ReadFile("/home/mhwwhu/projects/QuickStone/src/plugin/on_upload_plugins/hello/main.wasm")
	if err != nil {
		log.Fatalf("Failed to read WASM file: %v", err)
	}
	runtime, err := wasmrunner.NewWasmRuntime(ctx, file)
	if err != nil {
		log.Fatalf("Failed to create wasm runtime: %v", err)
	}

	defer runtime.Close(ctx)
	msg, err := runtime.Call(ctx, "process", "hello")
	if err != nil {
		log.Fatalf("Failed to call wasm function: %v", err)
	}
	var ret onuploadpluginmodels.OnUploadAction
	err = json.Unmarshal(msg, &ret)
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("from wasm: %s", ret)
}
