package parser

import (
	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/templates"
	"go/format"
	"golang.org/x/tools/imports"
)

type ParsedSrc interface {
	String() string
}
type File struct {
	Package    string
	Imports    []NamedTypeValue
	Constants  []NamedTypeValue
	Vars       []NamedTypeValue
	Interfaces []Interface
	Structs    []Struct
	Methods    []Method
}

type Struct struct {
	Name string
	Vars []NamedTypeValue
}
type Interface struct {
	Name    string
	Methods []Method
}
type Method struct {
	Name       string
	Struct     NamedTypeValue
	Body       string
	Parameters []NamedTypeValue
	Results    []NamedTypeValue
}
type NamedTypeValue struct {
	Name     string
	Type     string
	Value    string
	HasValue bool
}

func (f *File) String() string {
	s, err := template.NewEngine().Execute("file", f)
	if err != nil {
		logrus.Panic(err)
	}
	dt, err := imports.Process(f.Package, []byte(s), nil)
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}

func (m *Method) String() string {
	str := ""
	if m.Struct.Name != "" {
		s, err := template.NewEngine().ExecuteString("{{template \"struct_function\" .}}", m)
		if err != nil {
			logrus.Panic(err)
		}
		str = s
	} else {
		s, err := template.NewEngine().ExecuteString("{{template \"func\" .}}", m)
		if err != nil {
			logrus.Panic(err)
		}
		str = s
	}
	dt, err := format.Source([]byte(str))
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}

func (s *Struct) String() string {
	str, err := template.NewEngine().ExecuteString("{{template \"struct\" .}}", s)
	if err != nil {
		logrus.Panic(err)
	}
	dt, err := format.Source([]byte(str))
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}
func (i *Interface) String() string {
	str, err := template.NewEngine().ExecuteString("{{template \"interface\" .}}", i)
	if err != nil {
		logrus.Panic(err)
	}
	dt, err := format.Source([]byte(str))
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}

func NewNameType(name string, tp string) NamedTypeValue {
	return NamedTypeValue{
		Name:     name,
		Type:     tp,
		HasValue: false,
	}
}
func NewNameTypeValue(name string, tp string, vl string) NamedTypeValue {
	return NamedTypeValue{
		Name:     name,
		Type:     tp,
		HasValue: true,
		Value:    vl,
	}
}

func NewMethod(name string, str NamedTypeValue, body string, parameters, results []NamedTypeValue) Method {
	return Method{
		Name:       name,
		Struct:     str,
		Body:       body,
		Parameters: parameters,
		Results:    results,
	}
}

func NewInterface(name string, methods []Method) Interface {
	return Interface{
		Name:    name,
		Methods: methods,
	}
}
func NewStruct(name string, vars []NamedTypeValue) Struct {
	return Struct{
		Name: name,
		Vars: vars,
	}
}

func NewFile() File {
	return File{
		Interfaces: []Interface{},
		Imports:    []NamedTypeValue{},
		Structs:    []Struct{},
		Vars:       []NamedTypeValue{},
		Constants:  []NamedTypeValue{},
		Methods:    []Method{},
	}
}
