package parser

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
	Parse(src []byte) (ParsedSrc, error)
}

type FileParser struct{}

func (fp *FileParser) Parse(src []byte) (*File, error) {
	f := NewFile()
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	pf, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
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
				bd = bt.String()[1 : len(bt.String())-1]
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
				f.Imports = fp.parseImports(dec.Specs)
			case token.CONST:
				f.Constants = fp.parseConstants(dec.Specs)
			case token.VAR:
				f.Vars = fp.parseVars(dec.Specs)
			case token.TYPE:
				fp.parseType(dec.Specs, &f)
			default:
				logrus.Info("Skipping unknown Token Type")
			}
		}
	}
	//fmt.Println(f.String())
	return &f, nil
}
func (fp *FileParser) parseType(ds []ast.Spec, f *File) {
	for _, sp := range ds {
		tsp, ok := sp.(*ast.TypeSpec)
		if !ok {
			logrus.Debug("Type spec is not TypeSpec type, odd, skipping")
			continue
		}
		switch tsp.Type.(type) {
		case *ast.InterfaceType:
			ift := tsp.Type.(*ast.InterfaceType)
			mth := fp.parseFieldListAsMethods(ift.Methods)
			intr := NewInterface(tsp.Name.Name, mth)
			intr.Methods = mth
			f.Interfaces = append(f.Interfaces, intr)
		case *ast.StructType:
			st := tsp.Type.(*ast.StructType)
			str := NewStruct(tsp.Name.Name, fp.parseFieldListAsNamedTypes(st.Fields))
			f.Structs = append(f.Structs, str)
		default:
			logrus.Info("Skipping unknown type")
		}
	}
}
func (fp *FileParser) parseImports(ds []ast.Spec) []NamedTypeValue {
	imports := []NamedTypeValue{}
	for _, sp := range ds {
		isp, ok := sp.(*ast.ImportSpec)
		if !ok {
			logrus.Debug("Import spec is not ImportSpec type, odd, skipping")
			continue
		}
		ip := NewNameType("", "")
		if isp.Name != nil {
			ip.Name = isp.Name.Name
		}
		if isp.Path != nil {
			ip.Type = isp.Path.Value
		}
		imports = append(imports, ip)
	}
	return imports
}
func (fp *FileParser) parseVars(ds []ast.Spec) []NamedTypeValue {
	vars := []NamedTypeValue{}
	for _, sp := range ds {
		vsp, ok := sp.(*ast.ValueSpec)
		if !ok {
			logrus.Debug("Var spec is not ValueSpec type, odd, skipping")
			continue
		}
		tp, ok := vsp.Type.(*ast.Ident)
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
			vars = append(vars, NewNameTypeValue(tp.Name, vsp.Names[0].Name, bd))
		} else {
			vars = append(vars, NewNameType(tp.Name, vsp.Names[0].Name))
		}

	}
	return vars
}
func (fp *FileParser) parseConstants(ds []ast.Spec) []NamedTypeValue {
	constants := []NamedTypeValue{}
	for _, sp := range ds {
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
		tp, ok := vsp.Type.(*ast.Ident)
		if !ok {
			logrus.Debug("Spec type not  Ident type, odd, skipping")
			continue
		}
		constants = append(constants, NewNameTypeValue(tp.Name, vsp.Names[0].Name, bd))
	}
	return constants
}
func (fp *FileParser) parseFieldListAsNamedTypes(list *ast.FieldList) []NamedTypeValue {
	ntv := []NamedTypeValue{}
	if list != nil {
		for i, p := range list.List {
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
				names = append(names, typ[:1]+fmt.Sprintf("%d", i))
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
func (fp *FileParser) parseFieldListAsMethods(list *ast.FieldList) []Method {
	mth := []Method{}
	if list != nil {
		for _, p := range list.List {
			switch t := p.Type.(type) {
			case *ast.FuncType:
				m := Method{
					Name: p.Names[0].Name,
				}
				m.Parameters = fp.parseFieldListAsNamedTypes(t.Params)
				m.Results = fp.parseFieldListAsNamedTypes(t.Results)
				mth = append(mth, m)
			default:
				logrus.Info("Skipping unknown type")
			}
		}
	}
	return mth
}
func NewFileParser() *FileParser {
	return &FileParser{}
}
