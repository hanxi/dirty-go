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
	"os"
	"path/filepath"
	"strings"
)

func genNewFunc(name string) string {
	newFuncTmpl := `func New%s() *%s {
	p := &%s{}
	p.self = p
	p.root = p
	return p
}

`
	return fmt.Sprintf(newFuncTmpl, name, name, name)
}

func genWrapArr(name, fieldName, fieldType, arrEltType string) string {
	structTmpl := `type %s struct {
	Base
	%s %s
}

`
	structStr := fmt.Sprintf(structTmpl, name, fieldName, fieldType)

	newFuncTmpl := `func New%s() *%s {
	p := &%s{}
	p.%s = make(%s, 0)
	p.self = p
	p.root = p
	return p
}

`
	newFuncStr := fmt.Sprintf(newFuncTmpl, name, name, name, fieldName, fieldType)

	newFuncFromSliceTmpl := `func New%sFromSlice(%s %s) *%s {
	p := &%s{}
	p.%s = make(%s, 0)
	p.%s = append(p.%s, %s...)
	p.self = p
	p.root = p
	return p
}

`
	newFuncFromSliceStr := fmt.Sprintf(newFuncFromSliceTmpl, name, fieldName, fieldType, name, name, fieldName, fieldType, fieldName, fieldName, fieldName)

	appendFuncTmpl := `func (p *%s) Append(value %s) {
	if p == nil {
		return
	}
	p.%s = append(p.%s, value)
	value.root = p.root
	p.NotifyDirty()
}

`
	appendFuncStr := fmt.Sprintf(appendFuncTmpl, name, arrEltType, fieldName, fieldName)

	foreachFuncTmpl := `func (p *%s) Foreach(f func(%s)) {
	if p == nil {
		return
	}
	for _, v := range p.%s {
		f(v)
	}
}

`
	foreachFuncStr := fmt.Sprintf(foreachFuncTmpl, name, arrEltType, fieldName)

	return structStr + newFuncStr + newFuncFromSliceStr + appendFuncStr + foreachFuncStr
}

func genWrapMap(name, fieldName, fieldType, mapKeyType, mapValueType string) string {
	structTmpl := `type %s struct {
	Base
	%s %s
}

`
	structStr := fmt.Sprintf(structTmpl, name, fieldName, fieldType)

	newFuncTmpl := `func New%s() *%s {
	p := &%s{}
	p.%s = make(%s, 0)
	p.self = p
	p.root = p
	return p
}

`
	newFuncStr := fmt.Sprintf(newFuncTmpl, name, name, name, fieldName, fieldType)

	newFuncFromMapTmpl := `func New%sFromMap(%s %s) *%s {
	p := &%s{}
	p.%s = make(%s)
	for k, v := range %s {
		p.%s[k] = v
	}
	p.self = p
	p.root = p
	return p
}

`
	newFuncFromMapStr := fmt.Sprintf(newFuncFromMapTmpl, name, fieldName, fieldType, name, name, fieldName, fieldType, fieldName, fieldName)

	getFuncTmpl := `func (p *%s) Get(key %s) %s {
	if p == nil {
		return nil
	}
	return p.%s[key]
}

`
	getFuncStr := fmt.Sprintf(getFuncTmpl, name, mapKeyType, mapValueType, fieldName)

	setFuncTmpl := `func (p *%s) Set(key %s, value %s) {
	if p == nil {
		return
	}
	p.%s[key] = value
	value.root = p.root
	p.NotifyDirty()
}

`
	setFuncStr := fmt.Sprintf(setFuncTmpl, name, mapKeyType, mapValueType, fieldName)

	deleteFuncTmpl := `func (p *%s) Delete(key %s) {
	if p == nil {
		return
	}
	delete(p.%s, key)
	p.NotifyDirty()
}

`
	deleteFuncStr := fmt.Sprintf(deleteFuncTmpl, name, mapKeyType, fieldName)

	foreachFuncTmpl := `func (p *%s) Foreach(f func(%s, %s)) {
	if p == nil {
		return
	}
	for k, v := range p.%s {
		f(k, v)
	}
}

`
	foreachFuncStr := fmt.Sprintf(foreachFuncTmpl, name, mapKeyType, mapValueType, fieldName)

	return structStr + newFuncStr + newFuncFromMapStr + getFuncStr + setFuncStr + deleteFuncStr + foreachFuncStr
}

func genStruct(fset *token.FileSet, name string, list []*ast.Field) string {
	structTmpl := `type %s struct {
	Base%s
}

`

	eles := ""
	wraps := ""
	for _, field := range list {
		fieldName := field.Names[0].Name
		fieldType := getTypeString(fset, field.Type)

		ele := fieldName + " " + fieldType

		// wrap array
		arrType, isArr := field.Type.(*ast.ArrayType)
		if isArr {
			wrapStructName := "Arr" + name + strings.Title(fieldName)
			ele = "_wrap_" + fieldName + " *" + wrapStructName

			arrEltType := getTypeString(fset, arrType.Elt)
			wraps += genWrapArr(wrapStructName, fieldName, fieldType, arrEltType)
		}

		// wrap map
		mapType, isMap := field.Type.(*ast.MapType)
		if isMap {
			wrapStructName := "Map" + name + strings.Title(fieldName)
			ele = "_wrap_" + fieldName + " *" + wrapStructName

			mapKeyType := getTypeString(fset, mapType.Key)
			mapValueType := getTypeString(fset, mapType.Value)
			wraps += genWrapMap(wrapStructName, fieldName, fieldType, mapKeyType, mapValueType)
		}

		eles += "\n\t" + ele
	}
	return fmt.Sprintf(structTmpl, name, eles) + wraps
}

// Generate Get and Set methods for fields
func genGetSetFunc(fset *token.FileSet, name string, list []*ast.Field) string {
	setTmpl := `func (p *%s) Set%s(value %s) {
	if p == nil {
		return
	}
	p.%s = value
	p.NotifyDirty()
}

`
	getTmpl := `func (p *%s) Get%s() %s {
	if p == nil {
		return %s
	}
	return p.%s
}

`

	setStarStructTmpl := `func (p *%s) Set%s(value %s) {
	if p == nil {
		return
	}
	p.%s = value
	value.root = p.root
	p.NotifyDirty()
}

`

	setWrapTmpl := `func (p *%s) Set%s(value %s) {
	if p == nil {
		return
	}
	p._wrap_%s = value
	value.root = p.root
	for _, v := range value.%s {
		v.root = p.root
	}
	p.NotifyDirty()
}

`
	getWrapTmpl := `func (p *%s) Get%s() %s {
	if p == nil {
		return nil
	}
	return p._wrap_%s
}

`
	funcListStr := ""
	for _, field := range list {
		fieldName := field.Names[0].Name
		fieldType := getTypeString(fset, field.Type)

		if fieldIsStarStruct(field) {
			funcListStr += fmt.Sprintf(setStarStructTmpl, name, strings.Title(fieldName), fieldType, fieldName)
		} else {
			if fieldIsArrayStarStruct(field) {
				wrapType := "*Arr" + name + strings.Title(fieldName)
				funcListStr += fmt.Sprintf(setWrapTmpl, name, strings.Title(fieldName), wrapType, fieldName, fieldName)
			} else if fieldIsMapStarStruct(field) {
				wrapType := "*Map" + name + strings.Title(fieldName)
				funcListStr += fmt.Sprintf(setWrapTmpl, name, strings.Title(fieldName), wrapType, fieldName, fieldName)
			} else {
				funcListStr += fmt.Sprintf(setTmpl, name, strings.Title(fieldName), fieldType, fieldName)
			}
		}

		if fieldIsArrayStarStruct(field) {
			wrapType := "*Arr" + name + strings.Title(fieldName)
			funcListStr += fmt.Sprintf(getWrapTmpl, name, strings.Title(fieldName), wrapType, fieldName)
		} else if fieldIsMapStarStruct(field) {
			wrapType := "*Map" + name + strings.Title(fieldName)
			funcListStr += fmt.Sprintf(getWrapTmpl, name, strings.Title(fieldName), wrapType, fieldName)
		} else {
			funcListStr += fmt.Sprintf(getTmpl, name, strings.Title(fieldName), fieldType, getZeroValue(fieldType), fieldName)
		}
	}
	return funcListStr
}

func genDirtyOut(tmpl, out, outPackageName string) {
	src, err := ioutil.ReadFile(tmpl)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	headTmpl := `// Code generated by dirty-go; DO NOT EDIT.
package %s

`
	outStr := fmt.Sprintf(headTmpl, outPackageName)

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

			// Generate struct
			structStr := genStruct(fset, typeSpec.Name.Name, structType.Fields.List)
			outStr += structStr

			// Generate NewFunc
			newFuncStr := genNewFunc(typeSpec.Name.Name)
			outStr += newFuncStr

			// Generate Get and Set methods for fields
			getSetFuncStr := genGetSetFunc(fset, typeSpec.Name.Name, structType.Fields.List)
			outStr += getSetFuncStr
		}
	}
	writeCodeToFile(out, outStr)
}

func writeCodeToFile(out, code string) {
	// reload out code
	fsetOut := token.NewFileSet()
	fOut, err := parser.ParseFile(fsetOut, "", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// format out code and write to file
	outputFile, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	format.Node(outputFile, fsetOut, fOut)
	fmt.Println("generate:", out)
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

func writeBaseCode(out, outPackageName string) {
	baseTmpl := `// Code generated by dirty-go; DO NOT EDIT.
package %s

type Observer interface {
	OnDirty(interface{})
}

type DataObject interface {
	NotifyDirty()
	Attach(Observer)
}

type Base struct {
	DataObject
	observer Observer
	root     DataObject
	self     DataObject
}

func (x *Base) NotifyDirty() {
	if x.observer != nil {
		x.observer.OnDirty(x)
	}
	if x.root != nil && x.root != x.self {
		x.root.NotifyDirty()
	}
}

func (x *Base) Attach(o Observer) {
	x.observer = o
}
`

	outStr := fmt.Sprintf(baseTmpl, outPackageName)
	writeCodeToFile(out, outStr)
}

func isDir(f string) bool {
	fileInfo, err := os.Stat(f)
	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		return true
	}
	return false
}

func getFileNames(dir string) []string {
	files, _ := ioutil.ReadDir(dir)

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}

func main() {
	var in string
	var out string
	flag.StringVar(&in, "in", "", "The input directory.")
	flag.StringVar(&out, "out", "", "The output directory.")
	flag.Parse()

	if in == "" || out == "" {
		fmt.Println("Need input and output directory")
		flag.Usage()
		os.Exit(1)
	}

	if !isDir(in) {
		fmt.Println("input must be directory")
		flag.Usage()
		os.Exit(1)

	}

	if !isDir(out) {
		err := os.MkdirAll(out, 0755)
		if err != nil {
			fmt.Println(err)
			flag.Usage()
			os.Exit(1)
		}
	}

	outPackageName := filepath.Base(out)
	fileNames := getFileNames(in)
	for _, name := range fileNames {
		infile := filepath.Join(in, name)
		outfile := filepath.Join(out, name)
		genDirtyOut(infile, outfile, outPackageName)
	}

	outfile := filepath.Join(out, "base.go")
	writeBaseCode(outfile, outPackageName)
}
