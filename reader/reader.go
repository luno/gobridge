package reader

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func ParseFile(path string, mName string) (map[string]map[string]string, map[string][]FunctionSignature, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet()
	data := make(map[string]map[string]string)
	fs := make(map[string][]FunctionSignature)

	err = readFile(fset, mName, string(b), data, fs)
	if err != nil {
		return nil, nil, err
	}

	return data, fs, nil
}

func readFile(fset *token.FileSet, mName string, fContents string, data map[string]map[string]string, fs map[string][]FunctionSignature) error {
	node, err := parser.ParseFile(fset, "", fContents, parser.ParseComments)
	if err != nil {
		return err
	}

	r := &Reader{
		fset: fset,
	}

	sp := make(map[string]string)
	// Traverse file tree for interface methods
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if t.Name.IsExported() {
				switch assertion := t.Type.(type) {
				case *ast.StructType:
					sp = r.ListStructProperties(assertion)
					data[t.Name.Name] = sp
				case *ast.InterfaceType:
					apiFs := r.ListInterfaceMethods(assertion)
					fs[t.Name.Name] = apiFs
				}
			}
		case *ast.ImportSpec:
			if !strings.Contains(t.Path.Value, mName) {
				return false
			}

			path := strings.ReplaceAll(t.Path.Value, `"`, "")
			pathList := strings.Split(path, "/")
			path = strings.Join(pathList[1:], "/")

			fi, err := ioutil.ReadDir("./" + path)
			if err != nil {
				panic(err)
			}

			for _, f := range fi {
				b, err := ioutil.ReadFile("./" + path + "/" + f.Name())
				if err != nil {
					panic(err)
				}

				err = readFile(fset, mName, string(b), data, fs)
				if err != nil {
					panic(err)
				}
			}
		}
		return true
	})

	return nil
}

type Reader struct {
	fset *token.FileSet
}

// ListInterfaceMethods returns the function signatures of all the interface methods
func (r *Reader) ListStructProperties(it *ast.StructType) map[string]string {
	sf := make(map[string]string)
	for _, field := range it.Fields.List {
		// Name
		name := field.Names[0].Name
		// Type
		switch t := field.Type.(type) {
		case *ast.Ident:
			sf[name] = t.Name
		case *ast.SelectorExpr:
			p2 := t.Sel
			sf[name] = p2.Name
		case *ast.ArrayType:
			sf[name] = "[]" + r.importTypeFromASTExpr(t.Elt)
		case *ast.MapType:
			p1 := r.importTypeFromASTExpr(t.Key)
			p2 := r.importTypeFromASTExpr(t.Value)
			sf[name] = fmt.Sprintf("map[%s]%s", p1, p2)
		default:
			fmt.Println("Unable to determine type:")
			ast.Print(r.fset, t)
		}
	}
	return sf
}

type FunctionSignature struct {
	Name    string
	Params  map[string]string
	Results map[string]string
}

// ListInterfaceMethods returns the function signatures of all the interface methods
func (r *Reader) ListInterfaceMethods(it *ast.InterfaceType) []FunctionSignature {
	var fsSlice []FunctionSignature

	methods := it.Methods.List
	for _, method := range methods {
		fn, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		sig := r.CheckFunctionSignature(fn)
		sig.Name = method.Names[0].Name
		fsSlice = append(fsSlice, sig)
	}

	return fsSlice
}

func (r *Reader) CheckFunctionSignature(fn *ast.FuncType) FunctionSignature {
	var fs FunctionSignature
	fs.Results = make(map[string]string)
	fs.Params = make(map[string]string)

	for _, param := range fn.Params.List {
		n, v := r.parseVariables(param)
		if strings.ToLower(v) == "context" ||
			strings.ToLower(v) == "error" {
			continue
		}

		fs.Params[n] = v
	}

	for _, result := range fn.Results.List {
		n, v := r.parseVariables(result)
		if strings.ToLower(v) == "context" ||
			strings.ToLower(v) == "error" {
			continue
		}

		fs.Results[n] = v
	}

	return fs
}

func (r *Reader) parseVariables(field *ast.Field) (name string, importTyp string) {
	// Unnamed variables
	if len(field.Names) < 1 {
		return r.importTypeFromASTExpr(field.Type), r.importTypeFromASTExpr(field.Type)
	}

	return field.Names[0].Name, r.importTypeFromASTExpr(field.Type)
}

func (r *Reader) importTypeFromASTExpr(expr ast.Expr) string {
	switch s := expr.(type) {
	case *ast.Ident:
		return importTypeFromASTIdent(s)
	case *ast.SelectorExpr:
		return importTypeFromASTIdent(s.Sel)
	case *ast.ArrayType:
		return "[]" + r.importTypeFromASTExpr(s.Elt)
	case *ast.MapType:
		p1 := r.importTypeFromASTExpr(s.Key)
		p2 := r.importTypeFromASTExpr(s.Value)
		return fmt.Sprintf("map[%s]%s", p1, p2)
	default:
		return ""
	}
}

// importTypeFromASTIdent collects the type name and can be used for non-import types in the file tree
func importTypeFromASTIdent(ident *ast.Ident) string {
	return ident.Name
}

