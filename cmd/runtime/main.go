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
