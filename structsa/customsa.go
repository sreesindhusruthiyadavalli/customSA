package structsa

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"reflect"
)

const (
	Doc = "Tool to check usage of interfaces during self stim"
	unusedIntDoc = "Check for unused interfaces"
)


// Analyzer runs static analysis.
var Analyzer = &analysis.Analyzer{
	Name:     "checkInterfaceinStim",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
		(*ast.InterfaceType)(nil),
		(*ast.TypeSpec)(nil),
	}
	intfs := make([]string, 0)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch t := n.(type){
		case *ast.InterfaceType:
			//interfaces := t.Interface
			//fmt.Printf("interfaces %+v\n",interfaces)
		case *ast.StructType:
			//filter fields of struct type
			for _, field := range t.Fields.List {
				fmt.Println(field.Type, reflect.TypeOf(field.Type), reflect.ValueOf(field.Type).Interface())
				fmt.Println(field.Names[0].Name)
				switch field.Type.(type){
				case *ast.Ident:
					identName := field.Type.(*ast.Ident).Name
					fmt.Printf("ident: %+v, identName %+v\n",field.Type, identName)
					fmt.Printf("interfaces %+v\n", intfs)
					for _, intf := range intfs{
						if identName == intf {
							pass.Reportf(t.Pos(), "struct has %s Interface type.please fix this to any struct", identName)
						}
					}
				case *ast.InterfaceType:
					fmt.Printf("Inttype:%+v\n",field.Type)
				}
			}
		case *ast.TypeSpec:
			// which are public
			fmt.Printf("typespec %+v\n",t.Name)
			if t.Name.IsExported() {
				switch t.Type.(type) {
				// and are interfaces
				case *ast.InterfaceType:
					fmt.Printf("InterfaceType:%+v\n", t.Name.Name)
					intfs = append(intfs, t.Name.Name)
				}
			}
		}
	})
	return nil, nil
}

var UnusedInterfaceAnalyzer = &analysis.Analyzer{
	Name: "checkUnusedAnalyzer",
	Doc: unusedIntDoc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run: customrun,
}

func customrun(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.InterfaceType)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.StructType)(nil),
	}
	intfs := make([]string, 0)
	intfMaps := make(map[string]int)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch t := n.(type){
		case *ast.InterfaceType:
			//interfaces := t.Interface
			//fmt.Printf("interfaces %+v\n",interfaces)
		case *ast.TypeSpec:
			// which are public
			fmt.Printf("typespec %+v\n",t.Name)
			if t.Name.IsExported() {
				switch t.Type.(type) {
				// and are interfaces
				case *ast.InterfaceType:
					fmt.Printf("InterfaceType:%+v\n", t.Name.Name)
					intfs = append(intfs, t.Name.Name)
					intfMaps[t.Name.Name] = 1
				}
			}
		case *ast.StructType:
			for _, field := range t.Fields.List {
				fmt.Println(field.Type, reflect.TypeOf(field.Type), reflect.ValueOf(field.Type).Interface())
				fmt.Println(field.Names[0].Name)
				switch field.Type.(type){
				case *ast.Ident:
					identName := field.Type.(*ast.Ident).Name
					fmt.Printf("ident: %+v, identName %+v\n",field.Type, identName)
					//Remove interfaces which are used
					fmt.Printf("intfMaps %+v\n", intfMaps)
					for k := range intfMaps{
						if identName == k{
							fmt.Printf("used interface %+v\n",k)
							delete(intfMaps, k)
							fmt.Printf("intfMaps  after delete %+v\n", intfMaps)
						}
					}
				case *ast.InterfaceType:
					fmt.Printf("Inttype:%+v\n",field.Type)
				}
			}
		}
	})

	//nodeFilter2 := []ast.Node{
	//	(*ast.StructType)(nil),
	//}
	inspect.Preorder(nil, func(n ast.Node) {
		switch t := n.(type){
		case *ast.StructType:
			for x := range intfMaps {
				pass.Reportf(t.Pos(), "unused Interface %s", x)
			}
		}
	})
	return nil, nil
}