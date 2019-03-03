package varObf

import (
	"go/ast"
	"fmt"
)

type MainSource struct {
	Assignments  []*ast.AssignStmt // :=
	Values       []*ast.ValueSpec  // consts, =
	Literals     []*ast.BasicLit   // "asdf", 1234, 0xFFFFD00D
	Imports      []*ast.ImportSpec // all the Imports
	FunctionDecl []*ast.FuncDecl
}

func ParseMainSourceFromAST(node ast.Node) *MainSource {
	ret_val := MainSource{}

	// depth first iterate over each node in the AST Tree
	ast.Inspect(node, func (n ast.Node) bool{

		// if the node is an assignment
		assignments, ok := n.(*ast.AssignStmt)
		if ok {
			// add it to the list of Assignments
			ret_val.Assignments = append(ret_val.Assignments, assignments)
			return true
		}

		ident, ok := n.(*ast.Ident)
		if ok {
			fmt.Printf("GOT IDENT: %s\nIS EXPT: %t\nSTRING: %s\n\n", ident.Name, ident.IsExported(), ident.String())
			return true
		}

		imports, ok := n.(*ast.ImportSpec)
		if ok {
			ret_val.Imports = append(ret_val.Imports, imports)
			return true
		}
		var import_names []string
		for i := range ret_val.Imports {
			import_names = append(import_names, ret_val.Imports[i].Path.Value)
		}
		// ditto
		values, ok := n.(*ast.ValueSpec)
		if ok  {
			// ditto
			ret_val.Values = append(ret_val.Values, values)
			return true
		}

		// if the node is a literal
		vars, ok := n.(*ast.BasicLit)
		if ok && !StrContains(import_names, vars.Value) {
			// add it to our list of Literals
			ret_val.Literals = append(ret_val.Literals, vars)
			return true // our evaluation is done, don't recheck the same node
		}
		functionDecl, ok := n.(*ast.FuncDecl)
		if ok {
			// add it to our list of Literals
			ret_val.FunctionDecl = append(ret_val.FunctionDecl, functionDecl)
			return true // our evaluation is done, don't recheck the same node
		}
		return true
	})
	return  &ret_val
}


func NodeContains(slice []*ast.Node, item *ast.Node, max int) bool {
	for i := 0; i < max; i++ {
		if *item == *slice[i]{
			return true
		}
	}
	return false
}

func StrContains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}