//go:build wasm

package main

import (
	"fmt"

	"github.com/crypto2lab/ground/wasm_runtime"
)

func main() {}

func log(message string) {
	ptr, size := wasm_runtime.String2Ptr(message)
	_log(ptr, size)
}

//go:wasm-module env
//export log
func _log(ptr uint32, size uint32)

func greet(name string) {
	log(fmt.Sprintf("wasm>> Hello! %s", name))
}

//export greet
func _greet(ptr, size uint32) {
	name := wasm_runtime.Ptr2String(ptr, size)
	greet(name)
}

//export greeting
func _greeting(ptr, size uint32) (ptrSize uint64) {
	name := wasm_runtime.Ptr2String(ptr, size)
	g := fmt.Sprintf("Hello, %s!", name)

	ptr, size = wasm_runtime.String2Ptr(g)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}
