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

// 字符串首字母大写
func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// 字符串首字母小写
func firstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

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

func genWrapArr(name, fieldName, fieldType string, fset *token.FileSet, arrEltType ast.Expr, inPackageName string) string {
	lowerFieldName := firstLower(fieldName)
	arrEltTypeStr := getTypeString(fset, arrEltType)
	arrEltStarTypeStr := getStarTypeString(fset, arrEltType)

	structTmpl := `type %s struct {
	Base
	%s %s
}

`
	structStr := fmt.Sprintf(structTmpl, name, lowerFieldName, fieldType)

	newFuncTmpl := `func New%s() *%s {
	p := &%s{}
	p.%s = make(%s, 0)
	p.self = p
	p.root = p
	return p
}

`
	newFuncStr := fmt.Sprintf(newFuncTmpl, name, name, name, lowerFieldName, fieldType)

	newFuncFromSliceTmpl := `func New%sFromSlice(%s %s) *%s {
	p := &%s{}
	p.%s = make(%s, 0)
	p.%s = append(p.%s, %s...)
	p.self = p
	p.root = p
	return p
}

`
	newFuncFromSliceStr := fmt.Sprintf(newFuncFromSliceTmpl, name, lowerFieldName, fieldType, name, name, lowerFieldName, fieldType, lowerFieldName, lowerFieldName, lowerFieldName)

	appendFuncTmpl := `func (p *%s) Append(value %s) {
	if p == nil {
		return
	}
	p.%s = append(p.%s, value)
	value.root = p.root
	p.NotifyDirty()
}

`
	appendFuncStr := fmt.Sprintf(appendFuncTmpl, name, arrEltTypeStr, lowerFieldName, lowerFieldName)

	indexFuncTmpl := `func (p *%s) Index(i int) %s {
	if p == nil {
		return nil
	}
	if i < 0 || i >= len(p.%s) {
		return nil
	}
	return p.%s[i]
}

`
	indexFuncStr := fmt.Sprintf(indexFuncTmpl, name, arrEltTypeStr, lowerFieldName, lowerFieldName)

	foreachFuncTmpl := `func (p *%s) Foreach(f func(%s)) {
	if p == nil {
		return
	}
	for _, v := range p.%s {
		f(v)
	}
}

`
	foreachFuncStr := fmt.Sprintf(foreachFuncTmpl, name, arrEltTypeStr, lowerFieldName)

	arrImportTypeStr := fmt.Sprintf("[]*%s.%s", inPackageName, arrEltStarTypeStr)
	arrFromOriginTmpl := `func (p *%s) fromOrigin(o %s)  {
	for _,v := range o {
		val := New%s()
		val.fromOrigin(v)
		p.%s = append(p.%s, val)
	}
}

`
	arrFromOrigin := fmt.Sprintf(arrFromOriginTmpl, name, arrImportTypeStr, arrEltStarTypeStr, lowerFieldName, lowerFieldName)

	arrToOriginTmpl := `func (p *%s) toOrigin() %s {
	if p == nil {
		return nil
	}
	o := make(%s, 0)
	for _,v := range p.%s {
		o = append(o, v.toOrigin())
	}
	return o
}

`
	arrToOrigin := fmt.Sprintf(arrToOriginTmpl, name, arrImportTypeStr, arrImportTypeStr, lowerFieldName)

	return structStr + newFuncStr + newFuncFromSliceStr + appendFuncStr + indexFuncStr + foreachFuncStr + arrFromOrigin + arrToOrigin
}

func genWrapMap(name, fieldName, fieldType string, fset *token.FileSet, mapKeyType ast.Expr, mapValueType ast.Expr, inPackageName string) string {
	lowerFieldName := firstLower(fieldName)
	mapKeyTypeStr := getTypeString(fset, mapKeyType)
	mapValueTypeStr := getTypeString(fset, mapValueType)
	mapValueStarTypeStr := getStarTypeString(fset, mapValueType)

	structTmpl := `type %s struct {
	Base
	%s %s
}

`
	structStr := fmt.Sprintf(structTmpl, name, lowerFieldName, fieldType)

	newFuncTmpl := `func New%s() *%s {
	p := &%s{}
	p.%s = make(%s)
	p.self = p
	p.root = p
	return p
}

`
	newFuncStr := fmt.Sprintf(newFuncTmpl, name, name, name, lowerFieldName, fieldType)

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
	newFuncFromMapStr := fmt.Sprintf(newFuncFromMapTmpl, name, lowerFieldName, fieldType, name, name, lowerFieldName, fieldType, lowerFieldName, lowerFieldName)

	getFuncTmpl := `func (p *%s) Get(key %s) %s {
	if p == nil {
		return nil
	}
	return p.%s[key]
}

`
	getFuncStr := fmt.Sprintf(getFuncTmpl, name, mapKeyTypeStr, mapValueTypeStr, lowerFieldName)

	setFuncTmpl := `func (p *%s) Set(key %s, value %s) {
	if p == nil {
		return
	}
	p.%s[key] = value
	value.root = p.root
	p.NotifyDirty()
}

`
	setFuncStr := fmt.Sprintf(setFuncTmpl, name, mapKeyTypeStr, mapValueTypeStr, lowerFieldName)

	deleteFuncTmpl := `func (p *%s) Delete(key %s) {
	if p == nil {
		return
	}
	delete(p.%s, key)
	p.NotifyDirty()
}

`
	deleteFuncStr := fmt.Sprintf(deleteFuncTmpl, name, mapKeyTypeStr, lowerFieldName)

	foreachFuncTmpl := `func (p *%s) Foreach(f func(%s, %s)) {
	if p == nil {
		return
	}
	for k, v := range p.%s {
		f(k, v)
	}
}

`
	foreachFuncStr := fmt.Sprintf(foreachFuncTmpl, name, mapKeyTypeStr, mapValueTypeStr, lowerFieldName)

	mapImportTypeStr := fmt.Sprintf("map[%s]*%s.%s", mapKeyTypeStr, inPackageName, mapValueStarTypeStr)
	mapFromOriginTmpl := `func (p *%s) fromOrigin(o %s)  {
	for k,v := range o {
		p.%s[k] = New%s()
		p.%s[k].fromOrigin(v)
	}
}

`
	mapFromOrigin := fmt.Sprintf(mapFromOriginTmpl, name, mapImportTypeStr, lowerFieldName, mapValueStarTypeStr, lowerFieldName)

	mapToOriginTmpl := `func (p *%s) toOrigin() %s {
	if p == nil {
		return nil
	}
	o := make(%s)
	for k,v := range p.%s {
		o[k] = v.toOrigin()
	}
	return o
}

`
	mapToOrigin := fmt.Sprintf(mapToOriginTmpl, name, mapImportTypeStr, mapImportTypeStr, lowerFieldName)

	return structStr + newFuncStr + newFuncFromMapStr + getFuncStr + setFuncStr + deleteFuncStr + foreachFuncStr + mapFromOrigin + mapToOrigin
}

func genStruct(fset *token.FileSet, name string, list []*ast.Field, inPackageName string) string {
	structTmpl := `type %s struct {
	Base%s
}

`

	eles := ""
	wraps := ""
	for _, field := range list {
		fieldName := field.Names[0].Name
		fieldType := getTypeString(fset, field.Type)

		upperFieldName := firstUpper(fieldName)
		lowerFieldName := firstLower(fieldName)

		ele := lowerFieldName + " " + fieldType

		// wrap array
		arrType, isArr := field.Type.(*ast.ArrayType)
		if isArr {
			wrapStructName := "Arr" + name + upperFieldName
			ele = "_wrap_" + lowerFieldName + " *" + wrapStructName

			wraps += genWrapArr(wrapStructName, lowerFieldName, fieldType, fset, arrType.Elt, inPackageName)
		}

		// wrap map
		mapType, isMap := field.Type.(*ast.MapType)
		if isMap {
			wrapStructName := "Map" + name + upperFieldName
			ele = "_wrap_" + lowerFieldName + " *" + wrapStructName

			wraps += genWrapMap(wrapStructName, lowerFieldName, fieldType, fset, mapType.Key, mapType.Value, inPackageName)
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
		upperFieldName := firstUpper(fieldName)
		lowerFieldName := firstLower(fieldName)

		if fieldIsStarStruct(field) {
			funcListStr += fmt.Sprintf(setStarStructTmpl, name, upperFieldName, fieldType, lowerFieldName)
		} else {
			if fieldIsArrayStarStruct(field) {
				wrapType := "*Arr" + name + upperFieldName
				funcListStr += fmt.Sprintf(setWrapTmpl, name, upperFieldName, wrapType, lowerFieldName, lowerFieldName)
			} else if fieldIsMapStarStruct(field) {
				wrapType := "*Map" + name + upperFieldName
				funcListStr += fmt.Sprintf(setWrapTmpl, name, upperFieldName, wrapType, lowerFieldName, lowerFieldName)
			} else {
				funcListStr += fmt.Sprintf(setTmpl, name, upperFieldName, fieldType, lowerFieldName)
			}
		}

		if fieldIsArrayStarStruct(field) {
			wrapType := "*Arr" + name + upperFieldName
			funcListStr += fmt.Sprintf(getWrapTmpl, name, upperFieldName, wrapType, lowerFieldName)
		} else if fieldIsMapStarStruct(field) {
			wrapType := "*Map" + name + upperFieldName
			funcListStr += fmt.Sprintf(getWrapTmpl, name, upperFieldName, wrapType, lowerFieldName)
		} else {
			funcListStr += fmt.Sprintf(getTmpl, name, upperFieldName, fieldType, getZeroValue(fieldType), lowerFieldName)
		}
	}
	return funcListStr
}

func genJSONFunc(fset *token.FileSet, name string, list []*ast.Field, importPackage string) string {
	fromOriginTmpl := `func (p *%s) fromOrigin(o *%s.%s) {%s
}

`
	eles := ""
	for _, field := range list {
		fieldName := field.Names[0].Name
		fieldStarType := getStarTypeString(fset, field.Type)
		upperFieldName := firstUpper(fieldName)
		lowerFieldName := firstLower(fieldName)

		ele := fmt.Sprintf("p.%s = o.%s", lowerFieldName, upperFieldName)

		if fieldIsStarStruct(field) {
			ele = fmt.Sprintf("p.%s = New%s()\np.%s.fromOrigin(o.%s)", lowerFieldName, fieldStarType, lowerFieldName, upperFieldName)
		}

		_, isArr := field.Type.(*ast.ArrayType)
		_, isMap := field.Type.(*ast.MapType)
		var wrapStructName string
		if isArr {
			wrapStructName = "Arr" + name + upperFieldName
		}
		if isMap {
			wrapStructName = "Map" + name + upperFieldName
		}
		if isArr || isMap {
			ele = fmt.Sprintf("p._wrap_%s = New%s()\np._wrap_%s.fromOrigin(o.%s)", lowerFieldName, wrapStructName, lowerFieldName, upperFieldName)
		}

		eles += "\n\t" + ele
	}
	fromOrigin := fmt.Sprintf(fromOriginTmpl, name, importPackage, name, eles)

	toOriginTmpl := `func (p *%s) toOrigin() *%s.%s {
	if p == nil {
		return nil
	}
	o := &%s.%s{} %s
	return o
}

`
	eles = ""
	for _, field := range list {
		fieldName := field.Names[0].Name
		upperFieldName := firstUpper(fieldName)
		lowerFieldName := firstLower(fieldName)

		ele := fmt.Sprintf("o.%s = p.%s", upperFieldName, lowerFieldName)

		if fieldIsStarStruct(field) {
			ele = fmt.Sprintf("o.%s = p.%s.toOrigin()", upperFieldName, lowerFieldName)
		}

		_, isArr := field.Type.(*ast.ArrayType)
		_, isMap := field.Type.(*ast.MapType)
		if isArr || isMap {
			ele = fmt.Sprintf("o.%s = p._wrap_%s.toOrigin()", upperFieldName, lowerFieldName)
		}

		eles += "\n\t" + ele
	}
	toOrigin := fmt.Sprintf(toOriginTmpl, name, importPackage, name, importPackage, name, eles)

	unmarshalJSONTmpl := `func (p *%s) UnmarshalJSON(data []byte) error {
	origin := &%s.%s{}
	if err := json.Unmarshal(data, origin); err != nil {
        return err
    }
	p.fromOrigin(origin)
    return nil
}

`
	unmarshalJSON := fmt.Sprintf(unmarshalJSONTmpl, name, importPackage, name)

	marshalJSONTmpl := `func (p *%s) MarshalJSON() ([]byte, error) {
	origin := p.toOrigin()
    return json.Marshal(origin)
}

`
	marshalJSON := fmt.Sprintf(marshalJSONTmpl, name)

	return fromOrigin + toOrigin + unmarshalJSON + marshalJSON
}

func genDirtyOut(tmpl, out, outPackageName, importPackage string) {
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

import (
	"encoding/json"
	"%s"
)

`
	outStr := fmt.Sprintf(headTmpl, outPackageName, importPackage)

	inPackageName := filepath.Base(importPackage)
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

			name := typeSpec.Name.Name
			list := structType.Fields.List

			// Generate struct
			outStr += genStruct(fset, name, list, inPackageName)

			// Generate NewFunc
			outStr += genNewFunc(name)

			// Generate Get and Set methods for fields
			outStr += genGetSetFunc(fset, name, list)

			outStr += genJSONFunc(fset, name, list, inPackageName)
		}
	}
	writeCodeToFile(out, outStr)
}

func writeCodeToFile(out, code string) {
	// reload out code
	fsetOut := token.NewFileSet()
	fOut, err := parser.ParseFile(fsetOut, "", code, parser.ParseComments)
	if err != nil {
		fmt.Println(code)
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

func getStarTypeString(fset *token.FileSet, expr ast.Expr) string {
	starExpr, ok := expr.(*ast.StarExpr)
	if ok {
		return getTypeString(fset, starExpr.X)
	}

	return getTypeString(fset, expr)
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
	var importPackage string
	flag.StringVar(&in, "in", "", "The input directory.")
	flag.StringVar(&out, "out", "", "The output directory.")
	flag.StringVar(&importPackage, "import", "", "The import package.")
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
		genDirtyOut(infile, outfile, outPackageName, importPackage)
	}

	outfile := filepath.Join(out, "base.go")
	writeBaseCode(outfile, outPackageName)
}
