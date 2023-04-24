package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto2lab/ground/internal/chainspec"
)

const defaultBadgerLocation = "./tmp"

func importRuntimeToChainSpec(runtimeFilePath string, output string) error {
	runtimePath := os.Args[1]

	runtimeBytes, err := os.ReadFile(runtimePath)
	if err != nil {
		return fmt.Errorf("reading runtime: %w", err)
	}

	wasmHexBlob := hex.EncodeToString(runtimeBytes)
	wasmHexBlob = "0x" + wasmHexBlob

	genericChainSpec := chainspec.ChainSpec{
		Genesis: &chainspec.Genesis{
			Code: wasmHexBlob,
			Accounts: []*chainspec.ChainSpecAccount{
				{
					PublicAddress: "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
					Currency:      "TNT",
					Balance:       1000000000000,
				},
				{
					PublicAddress: "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty",
					Currency:      "TNT",
					Balance:       1000000000000,
				},
			},
		},
	}

	chainSpecBytes, err := json.MarshalIndent(genericChainSpec, "", "    ")
	if err != nil {
		return fmt.Errorf("while marshaling chainspec: %w", err)
	}

	err = os.WriteFile(output, chainSpecBytes, os.ModeType|os.ModePerm)
	if err != nil {
		return fmt.Errorf("while writing chain_spec.json: %w", err)
	}

	return nil
}

func main() {
	const testNetChainSpec = "./chain/testnet/chain_spec.json"
	runtimePath := os.Args[1]
	err := importRuntimeToChainSpec(runtimePath, testNetChainSpec)
	if err != nil {
		log.Fatalf("importing runtime to chainspec: %s", err)
	}

	client, err := StartClient(testNetChainSpec)
	if err != nil {
		log.Fatalf("starting client: %s", err)
	}

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownSig

	client.Stop()
	log.Printf("client sucessfully shutdown!\n")
}

// func main() {
// 	runtimePath := os.Args[1]

// 	runtimeBytes, err := os.ReadFile(runtimePath)
// 	if err != nil {
// 		log.Fatalf("while reading runtime: %s", err)
// 	}

// 	ctx := context.Background()
// 	r := wazero.NewRuntime(ctx)
// 	defer r.Close(ctx)

// 	_, err = r.NewHostModuleBuilder("env").
// 		NewFunctionBuilder().WithFunc(logString).Export("log").
// 		Instantiate(ctx)
// 	if err != nil {
// 		log.Fatalf("binding env:log function: %s", err)
// 	}
// 	wasi_snapshot_preview1.MustInstantiate(ctx, r)

// 	mod, err := r.Instantiate(ctx, runtimeBytes)
// 	if err != nil {
// 		log.Fatalf("while instantiating wasm blob: %s", err)
// 	}

// 	greet := mod.ExportedFunction("greet")
// 	greeting := mod.ExportedFunction("greeting")
// 	malloc := mod.ExportedFunction("malloc")
// 	free := mod.ExportedFunction("free")

// 	name := "Eclesio"
// 	nameSize := uint64(len(name))

// 	allocResults, err := malloc.Call(ctx, nameSize)
// 	if err != nil {
// 		log.Fatalf("while allocating: %s", err)
// 	}
// 	namePtr := allocResults[0]
// 	defer free.Call(ctx, namePtr)

// 	ok := mod.Memory().Write(uint32(namePtr), []byte(name))
// 	if !ok {
// 		log.Fatalf("Memory.Write(%d, %d) out of range of memory size %d",
// 			namePtr, nameSize, mod.Memory().Size())
// 	}

// 	_, err = greet.Call(ctx, namePtr, nameSize)
// 	if err != nil {
// 		log.Fatalf("while calling greet: %s", err)
// 	}

// 	results, err := greeting.Call(ctx, namePtr, nameSize)
// 	if err != nil {
// 		log.Fatalf("while calling greeting: %s", err)
// 	}

// 	encodedPtrData := results[0]
// 	ptr := uint32(encodedPtrData >> 32)
// 	size := uint32(encodedPtrData)

// 	bytes, ok := mod.Memory().Read(ptr, size)
// 	if !ok {
// 		log.Fatalf("Memory.Read(%d, %d) out of range of memory size %d",
// 			ptr, size, mod.Memory().Size())
// 	}

// 	fmt.Println("go >>", string(bytes))
// }

// func logString(_ context.Context, m api.Module, offset, byteCount uint32) {
// 	buf, ok := m.Memory().Read(offset, byteCount)
// 	if !ok {
// 		log.Fatalf("failed to read the memory(%d, %d), out of range", offset, byteCount)
// 	}
// 	fmt.Println(string(buf))
// }
