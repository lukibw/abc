package main

import (
	"log"
	"os"

	"github.com/lukibw/abc/compiler"
	"github.com/lukibw/abc/scanner"
	"github.com/lukibw/abc/vm"
)

func main() {
	debug, err := os.Create("main.log")
	if err != nil {
		log.Fatalln(err)
	}
	content, err := os.ReadFile("main.abc")
	if err != nil {
		log.Fatalln(err)
	}
	vm, err := vm.New(compiler.New(scanner.New(content)), log.New(debug, "", 0))
	if err != nil {
		log.Fatalln(err)
	}
	if err = vm.Run(); err != nil {
		log.Fatalln(err)
	}
}
