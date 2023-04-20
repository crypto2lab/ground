package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/OneOfOne/xxhash"
	"github.com/crypto2lab/ground/internal/chainspec"
	"github.com/crypto2lab/ground/internal/primitives"
)

const defaultBadgerLocation = "./tmp"

func startClient(chainspecFilePath string) error {
	spec := chainspec.ChainSpec{}
	chainSpecBytes, err := os.ReadFile(chainspecFilePath)
	if err != nil {
		return fmt.Errorf("while reading chainspec: %w", err)
	}

	err = json.Unmarshal(chainSpecBytes, &spec)
	if err != nil {
		return fmt.Errorf("while unmarshaling chainspec: %w", err)
	}

	database := NewDatabase()
	err = database.Open()
	if err != nil {
		return fmt.Errorf("while opening database: %w", err)
	}
	defer database.Close()

	err = spec.StoreGenesis(database)
	if err != nil {
		return fmt.Errorf("while storing genesis: %w", err)
	}

	retrieveAliceTNTBalance(database)
	return nil
}

func retrieveAliceTNTBalance(db *Database) {
	pubKey := primitives.PublicAddress("5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY")
	encPubKey, err := pubKey.Encode()
	if err != nil {
		log.Fatal(err)
		return
	}

	currency := primitives.Currency("TNT")
	encCurrency, err := currency.Encode()
	if err != nil {
		log.Fatal(err)
		return
	}

	//mod::module::publickey::currency
	storageKey := fmt.Sprintf("mod::%s::%s::%s",
		"ACCOUNT::", encPubKey, encCurrency)

	hasher := xxhash.NewS64(0)
	hasher.WriteString(storageKey)

	result := hasher.Sum64()
	aliceAccountKey := make([]byte, 8)
	binary.LittleEndian.PutUint64(aliceAccountKey, result)

	value, err := db.Get(aliceAccountKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	var expected uint64 = 1000000000000
	got := binary.LittleEndian.Uint64(value)

	fmt.Printf("expected=%v, got=%v\n", expected, got)
}

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
			Runtime: wasmHexBlob,
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

	startClient(testNetChainSpec)
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
