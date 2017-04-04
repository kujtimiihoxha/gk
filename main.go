package main

import (
	"fmt"
	"github.com/kujtimiihoxha/gk/generator"
)

func main() {
	f := generator.File{
		Package: "test",
		Imports: []generator.NamedTypeValue{},
		Vars: []generator.NamedTypeValue{
			generator.NewNameType("test", "string"),
			generator.NewNameTypeValue("i", "int", "4"),
		},
		Constants: []generator.NamedTypeValue{
			generator.NewNameTypeValue("ic", "int", "4"),
			generator.NewNameTypeValue("s", "string", "\"myVal\""),
		},
		Interfaces: []generator.Interface{
			generator.NewInterface(
				"IName",
				[]generator.Method{
					generator.NewMethod(
						"MyMeth",
						generator.NamedTypeValue{},
						"",
						[]generator.NamedTypeValue{},
						[]generator.NamedTypeValue{},
					),
				},
			),
		},
		Structs: []generator.Struct{
			generator.NewStruct(
				"MyStruct",
				[]generator.NamedTypeValue{
					generator.NewNameType("test", "string"),
				},
			),
		},
		Methods: []generator.Method{
			generator.NewMethod(
				"Normal",
				generator.NamedTypeValue{},
				"",
				[]generator.NamedTypeValue{
					generator.NewNameType("input", "string"),
				},
				[]generator.NamedTypeValue{
					generator.NewNameType("res", "string"),
				},
			),
			generator.NewMethod(
				"WithBody",
				generator.NamedTypeValue{},
				`a:=\"hello\"
				fmt.Println(a)
				`,
				[]generator.NamedTypeValue{
					generator.NewNameType("input", "string"),
				},
				[]generator.NamedTypeValue{
					generator.NewNameType("res", "string"),
				},
			),
			generator.NewMethod(
				"WithStruct",
				generator.NamedTypeValue{Name:"mstr",Type:"*MyStruct"},
				"",
				[]generator.NamedTypeValue{
					generator.NewNameType("input", "string"),
				},
				[]generator.NamedTypeValue{
					generator.NewNameType("res", "string"),
				},
			),
			generator.NewMethod(
				"WithStructBody",
				generator.NamedTypeValue{Name:"mstr",Type:"*MyStruct"},
				`a:="hello"
			fmt.Println(a)
			`,
				[]generator.NamedTypeValue{
					generator.NewNameType("input", "string"),
				},
				[]generator.NamedTypeValue{
					generator.NewNameType("res", "string"),
				},
			),
		},
	}
	//v, e := t.Execute("service_interface", i)
	//fmt.Println(e)
	//b, _ := format.Source([]byte(v))
	fmt.Println(f.String())
}
