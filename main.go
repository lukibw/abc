package main

import (
	"bytes"
	"log"
	"os"

	"github.com/lukibw/abc/compiler"
	"github.com/lukibw/abc/scanner"
	"github.com/lukibw/abc/vm"
)

func main() {
	content, err := os.ReadFile("main.abc")
	if err != nil {
		log.Fatalln(err)
	}
	vm, err := vm.New(compiler.New(scanner.New(bytes.Runes(content))))
	if err != nil {
		log.Fatalln(err)
	}
	if err = vm.Run(); err != nil {
		log.Fatalln(err)
	}
}
