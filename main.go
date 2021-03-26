package main

import (
	"C"

	"github.com/napi-bindings/go-node-api/napi"
)
import "fmt"

func main() {
	fmt.Println(napi.Hello())
	napi.Hello1()
}
