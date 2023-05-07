package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func main() {
	var tmpl string
	var out string
	flag.StringVar(&tmpl, "tmpl", "", "The input template file.")
	flag.StringVar(&out, "out", "", "The output file.")
	flag.Parse()

	src, err := ioutil.ReadFile(tmpl)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("package %s\n\n", f.Name.Name)

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			fmt.Printf("func New%s() *%s {\n\tp := &%s{}\n\tp.self = p\n\tp.root = p\n\treturn p\n}\n\n", typeSpec.Name.Name, typeSpec.Name.Name, typeSpec.Name.Name)

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				fieldName := field.Names[0].Name
				fieldType := getTypeString(fset, field.Type)

				// Generate Get and Set methods for all fields
				if fieldIsStarStruct(field) {
					fmt.Printf("func (p *%s) Set%s(value %s) {\n\tif p == nil {\n\t\treturn\n\t}\n\tp.%s = value\n\tvalue.root = p.root\n\tp.NotifyDirty()\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), fieldType, fieldName)
				} else {
					if fieldIsArrayStarStruct(field) || fieldIsMapStarStruct(field) {
						setRoot := "\n\tfor _,v := range value {\n\t\tv.root = p.root\n\t}"
						fmt.Printf("func (p *%s) Set%s(value %s) {\n\tif p == nil {\n\t\treturn\n\t}\n\tp.%s = value%s\n\tp.NotifyDirty()\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), fieldType, fieldName, setRoot)
					} else {
						fmt.Printf("func (p *%s) Set%s(value %s) {\n\tif p == nil {\n\t\treturn\n\t}\n\tp.%s = value\n\tp.NotifyDirty()\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), fieldType, fieldName)
					}
				}
				fmt.Printf("func (p *%s) Get%s() %s {\n\tif p == nil {\n\t\treturn %s\n\t}\n\treturn p.%s\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), fieldType, getZeroValue(fieldType), fieldName)

				// Generate Append method for slice fields
				arrType, ok := field.Type.(*ast.ArrayType)
				if ok {
					if isStarStruct(arrType.Elt) {
						fmt.Printf("func (p *%s) Append%s(value %s) {\n\tif p == nil {\n\t\treturn\n\t}\n\tp.%s = append(p.%s, value)\n\tvalue.root = p.root\n\tp.NotifyDirty()\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), fieldType[2:], fieldName, fieldName)
					} else {
						fmt.Printf("func (p *%s) Append%s(value %s) {\n\tif p == nil {\n\t\treturn\n\t}\n\tp.%s = append(p.%s, value)\n\tp.NotifyDirty()\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), fieldType[2:], fieldName, fieldName)
					}
				}

				// Generate Get and Set methods for map fields
				mapType, ok := field.Type.(*ast.MapType)
				if ok {
					keyType := getTypeString(fset, mapType.Key)
					valueType := getTypeString(fset, mapType.Value)

					if isStarStruct(mapType.Value) {
						fmt.Printf("func (p *%s) Put%s(key %s, value %s) {\n\tif p == nil {\n\t\treturn\n\t}\n\tp.%s[key] = value\n\tvalue.root = p.root\n\tp.NotifyDirty()\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), keyType, valueType, fieldName)
					} else {
						fmt.Printf("func (p *%s) Put%s(key %s, value %s) {\n\tif p == nil {\n\t\treturn\n\t}\n\tp.%s[key] = value\n\tp.NotifyDirty()\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), keyType, valueType, fieldName)
					}

					fmt.Printf("func (p *%s) Lookup%s(key %s) %s {\n\tif p == nil {\n\t\treturn %s\n\t}\n\treturn p.%s[key]\n}\n\n", typeSpec.Name.Name, strings.Title(fieldName), keyType, valueType, getZeroValue(valueType), fieldName)
				}
			}
		}
	}
}

func getTypeString(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, expr); err != nil {
		panic(err)
	}
	return buf.String()
}

func getZeroValue(fieldType string) string {
	switch fieldType {
	case "int", "int8", "int16", "int32", "int64":
		return "0"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "bool":
		return "false"
	case "string":
		return "\"\""
	default:
		return "nil"
	}
}

func fieldIsStarStruct(field *ast.Field) bool {
	return isStarStruct(field.Type)
}

func isStarStruct(expr ast.Expr) bool {
	starExpr, ok := expr.(*ast.StarExpr)
	if ok {
		_, ok = starExpr.X.(*ast.Ident)
		if ok {
			return true
		}
	}

	return false
}

func fieldIsArrayStarStruct(field *ast.Field) bool {
	arrType, ok := field.Type.(*ast.ArrayType)
	if ok {
		if isStarStruct(arrType.Elt) {
			return true
		}
	}
	return false
}

func fieldIsMapStarStruct(field *ast.Field) bool {
	mapType, ok := field.Type.(*ast.MapType)
	if ok {
		if isStarStruct(mapType.Value) {
			return true
		}
	}
	return false
}
