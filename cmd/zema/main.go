package main

import (
	"fmt"
	"os"
	"zema/internal/zema"
)

type code int

const (
	_ int = iota
	initCode
	fatalCode
)

func main() {
	_zema, err := zema.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(initCode)
	}
	if err := _zema.Run(); err != nil {
		fmt.Println(err)
		os.Exit(fatalCode)
	}
}
