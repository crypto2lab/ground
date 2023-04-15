package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	runtimePath := os.Args[1]

	runtimeBytes, err := os.ReadFile(runtimePath)
	if err != nil {
		log.Fatalf("while reading runtime: %s", err)
	}

	ctx := context.Background()
	runtimeInstance := wazero.NewRuntime(ctx)
	defer runtimeInstance.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtimeInstance)

	mod, err := runtimeInstance.Instantiate(ctx, runtimeBytes)
	if err != nil {
		log.Fatalf("while instantiating wasm blob: %s", err)
	}

	mul := mod.ExportedFunction("multiply")
	result, err := mul.Call(ctx, 2, 2)
	if err != nil {
		log.Fatalf("while calling multiply function: %s", err)
	}

	fmt.Printf("result: %v\n", result)
}
