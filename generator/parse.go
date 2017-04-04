package generator

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

type Parser interface {
	Parse(src []byte) ParsedSrc
}

type FileParser struct{}

func (fp *FileParser) Parse(src []byte) *File {
	f := NewFile()
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	pf, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	f.Package = pf.Name.Name
	for _, v := range pf.Decls {
		if dec, ok := v.(*ast.FuncDecl); ok {
			st := []NamedTypeValue{}
			pr := []NamedTypeValue{}
			rs := []NamedTypeValue{}
			if dec.Recv != nil {
				st = fp.parseFieldListAsNamedTypes(dec.Recv)
			}
			if dec.Type != nil {
				pr = fp.parseFieldListAsNamedTypes(dec.Type.Params)
				rs = fp.parseFieldListAsNamedTypes(dec.Type.Results)
			}
			bd := ""
			if dec.Body != nil {
				fst := token.NewFileSet()
				bt := bytes.NewBufferString("")
				err := format.Node(bt, fst, dec.Body)
				bd = bt.String()
				if err != nil {
					logrus.Panic(err)
				}
			}
			str := NamedTypeValue{}
			if len(st) > 0 {
				str = st[0]
			}
			fc := NewMethod(dec.Name.String(), str, bd, pr, rs)
			f.Methods = append(f.Methods, fc)
		}
		if dec, ok := v.(*ast.GenDecl); ok {
			switch dec.Tok {
			case token.IMPORT:
				for _,sp := range dec.Specs{
					isp, ok := sp.(*ast.ImportSpec)
					if !ok {
						logrus.Debug("Import spec is not ImportSpec type, odd, skipping")
						continue
					}
					ip := NewNameType("","")
					if isp.Name != nil{
						ip.Name = isp.Name.Name
					}
					if isp.Path != nil {
						ip.Type = isp.Path.Value
					}
					fmt.Println(ip)
					f.Imports = append(f.Imports,ip)
				}
			case token.CONST:
				for _,sp := range dec.Specs{
					vsp, ok := sp.(*ast.ValueSpec)
					if !ok {
						logrus.Debug("Constant spec is not ValueSpec type, odd, skipping")
						continue
					}
					fst := token.NewFileSet()
					bt := bytes.NewBufferString("")
					err := format.Node(bt, fst, vsp.Values[0])
					bd := bt.String()
					if err != nil {
						logrus.Panic(err)
					}
					tp,ok := vsp.Type.(*ast.Ident)
					if !ok {
						logrus.Debug("Spec type not  Ident type, odd, skipping")
						continue
					}
					f.Constants = append(f.Constants,NewNameTypeValue(tp.Name,vsp.Names[0].Name,bd))
				}
			case token.VAR:
				for _,sp := range dec.Specs{
					vsp, ok := sp.(*ast.ValueSpec)
					if !ok {
						logrus.Debug("Var spec is not ValueSpec type, odd, skipping")
						continue
					}
					tp,ok := vsp.Type.(*ast.Ident)
					if !ok {
						logrus.Debug("Spec type not  Ident type, odd, skipping")
						continue
					}
					if len(vsp.Values) > 0 {
						fst := token.NewFileSet()
						bt := bytes.NewBufferString("")
						err := format.Node(bt, fst, vsp.Values[0])
						bd := bt.String()
						if err != nil {
							logrus.Panic(err)
						}
						f.Vars = append(f.Constants,NewNameTypeValue(tp.Name,vsp.Names[0].Name,bd))
					} else {
						f.Vars = append(f.Constants,NewNameType(tp.Name,vsp.Names[0].Name))
					}

				}

			default:
				logrus.Info("Skipping unknown Token Type")
			}
		}
	}
	fmt.Println(f.String())
	return &f
}
func (fp *FileParser) parseFieldListAsNamedTypes(list *ast.FieldList) []NamedTypeValue {
	ntv := []NamedTypeValue{}
	if list != nil {
		for _, p := range list.List {
			var typ string
			switch t := p.Type.(type) {
			case *ast.Ident:
				logrus.Debug("Type Ident, i.e. a built-in type")
				typ = t.Name

			case *ast.SelectorExpr:
				logrus.Debug("Type Selector, i.e. a third-party type")
				selectorIdent, ok := t.X.(*ast.Ident)
				if !ok {
					logrus.Debug("Selector X isn't an Ident; very odd, skipping")
					continue
				}
				typ = fmt.Sprintf(fmt.Sprintf("%s.%s", selectorIdent.Name, t.Sel.Name))
			case *ast.StarExpr:
				starIndent, ok := t.X.(*ast.Ident)
				if !ok {
					logrus.Debug("Selector X isn't an Ident; very odd, skipping")
					continue
				}
				typ = "*" + starIndent.Name
			default:
				logrus.Info("Skipping unknown Field Type")
				continue
			}
			logrus.Debug(fmt.Sprintf("Type %s", typ))

			// Potentially N names
			var names []string
			for _, ident := range p.Names {
				names = append(names, ident.Name)
			}
			if len(names) == 0 {
				// Anonymous named type, give it a default name
				names = append(names, "somename") // TODO(pb): generator
			}
			for _, name := range names {
				namedType := NewNameType(name, typ)
				logrus.Debug(fmt.Sprintf("NamedType %+v", namedType))
				ntv = append(ntv, namedType)
			}
		}
	}
	return ntv
}
func NewFileParser() *FileParser {
	return &FileParser{}
}
