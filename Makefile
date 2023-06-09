tinygo-runtime-build:
	tinygo build -o ./target/runtime.wasm -scheduler=none --no-debug -tags=wasm -target=wasi ./cmd/runtime/...

run: tinygo-runtime-build
	rm -rf ./tmp
	go run ./cmd/client/... ./target/runtime.wasm