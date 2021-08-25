package structsa

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"reflect"
)

const Doc = "Tool to check usage of interfaces during self stim"


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
		//chnls := make(chan string, 10)
			//make(chan string)
		switch t := n.(type){
		case *ast.InterfaceType:
			//interfaces := t.Interface
			//fmt.Printf("interfaces %+v\n",interfaces)
		case *ast.StructType:
			////filter self stims
			for _, field := range t.Fields.List {
				fmt.Println(field.Type, reflect.TypeOf(field.Type), reflect.ValueOf(field.Type).Interface())
				//fmt.Println(field.Type.Pos(), field.Type.End())
				fmt.Println(field.Names[0].Name)
				//newType := field.Type.(*ast.Ident)
				//fmt.Println(newType)
				//if newType.Obj != nil{
				//	fmt.Println(newType.Obj.Kind)
				//}
				switch field.Type.(type){
				case *ast.Ident:
					identName := field.Type.(*ast.Ident).Name
					fmt.Printf("ident: %+v, identName %+v\n",field.Type, identName)
					//chanLen := len(chnls)
					////for i:=0;  chnls {
					////	intfs = append(intfs, <-interf)
					////}
					//for i := 0; i < chanLen; i++ {
					//	intfs = append(intfs, <-chnls)
					//}
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
					//ch := make(chan string)
					//chnls <- t.Name.Name
					//chnls = append(chnls, ch)
					//fmt.Printf("hereintfsL%+v\n", intfs)
				}
			}
		}
		//interfaces := n.(*ast.InterfaceType)
		//fmt.Println(interfaces.Interface)

		//function := n.(*ast.StructType)



		//switch t := n.(type) {
		//// find variable declarations


		// operation here on functions
	})
	return nil, nil
}

