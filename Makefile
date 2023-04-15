tinygo-runtime-build:
	tinygo build -o ./target/runtime.wasm -target wasm ./cmd/runtime/...

run: tinygo-runtime-build
	go run ./cmd/client/... ./target/runtime.wasm