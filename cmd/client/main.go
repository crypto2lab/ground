package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	runtimePath := os.Args[1]

	runtimeBytes, err := os.ReadFile(runtimePath)
	if err != nil {
		log.Fatalf("while reading runtime: %s", err)
	}

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	_, err = r.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(logString).Export("log").
		Instantiate(ctx)
	if err != nil {
		log.Fatalf("binding env:log function: %s", err)
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.Instantiate(ctx, runtimeBytes)
	if err != nil {
		log.Fatalf("while instantiating wasm blob: %s", err)
	}

	greet := mod.ExportedFunction("greet")
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	name := "Eclesio"
	nameSize := uint64(len(name))

	allocResults, err := malloc.Call(ctx, nameSize)
	if err != nil {
		log.Fatalf("while allocating: %s", err)
	}
	namePtr := allocResults[0]
	defer free.Call(ctx, namePtr)

	ok := mod.Memory().Write(uint32(namePtr), []byte(name))
	if !ok {
		log.Fatalf("Memory.Write(%d, %d) out of range of memory size %d",
			namePtr, nameSize, mod.Memory().Size())
	}

	_, err = greet.Call(ctx, namePtr, nameSize)
	if err != nil {
		log.Fatalf("while calling greet: %s", err)
	}
}

func logString(_ context.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		log.Fatalf("failed to read the memory(%d, %d), out of range", offset, byteCount)
	}
	fmt.Println(string(buf))
}
