package parser

import (
	"fmt"
	"golang.org/x/tools/imports"
	"testing"
)

func TestName(t *testing.T) {
	s, err := imports.Process("test", []byte(`
	package main
	func main(){
	s:= generator.ServiceGenerator{}
	fmt.Println(s)
	}
	`), nil)
	fmt.Println(string(s), err)

}
