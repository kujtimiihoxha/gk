package parser

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	p := NewFileParser()
	v, _ := p.Parse([]byte(`
	package service
type ABC struct {

}
type MyService interface {
	// Write your interface methods
	Foo(a []string) (map[string][]*parser.FileParser, error)
}

	`))
	fmt.Println(v.String())

}
