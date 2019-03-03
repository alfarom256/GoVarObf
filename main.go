package main

import (
	"fmt"
	"go/token"
	"go/parser"
	"log"
	"GoVarObf/varObf"
	"go/ast"
	"go/printer"
	"os"
)

func main(){
	welcome_str := "Hello world!!!"
	fmt.Printf("%s\n", welcome_str)
	fPath := "C:/Users/Michael/go/src/GoStrObf/main.go"
	mySource := varObf.MainSource{}
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, fPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal("fuckin a, then!")
		log.Fatal(err)
		panic(err)
	}
	mySource = *varObf.ParseMainSourceFromAST(node)

	// carve out a map (table) that will store a list of Function nodes
	// the value will be all the ident variables in the function
	AllVars := make(map[*ast.FuncDecl][]*ast.Ident, len(mySource.FunctionDecl))
	for i := range mySource.FunctionDecl{
		tmp := varObf.VarsFromFunc(mySource.FunctionDecl[i])
		AllVars[mySource.FunctionDecl[i]]= tmp
	}

	for i := range AllVars {
		varObf.ChangeVarsFuncAST(node, AllVars)
		_ = i
	}

	printer.Fprint(os.Stdout, fset, node)
}