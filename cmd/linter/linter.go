package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var Analyzer = &analysis.Analyzer{
	Name: "linter",
	Doc:  "Check for panic and call os.Exit/log.Fatal outside of main function",
	Run:  run,
}

func main() {
	singlechecker.Main(Analyzer)
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		packageName := ""
		inMainFunc := false
		ast.Inspect(file, func(n ast.Node) bool {
			if fileDecl, ok := n.(*ast.File); ok {
				packageName = fileDecl.Name.Name
			}

			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				inMainFunc = funcDecl.Name.Name == "main" && packageName == "main"
			}

			if callExpr, ok := n.(*ast.CallExpr); ok {
				switch fun := callExpr.Fun.(type) {
				case *ast.Ident:
					if fun.Name == "panic" {
						pass.Reportf(callExpr.Pos(), "panic detected")
					}

				case *ast.SelectorExpr:
					if !inMainFunc {
						if pkgIdent, ok := fun.X.(*ast.Ident); ok && pkgIdent.Name == "log" {
							if strings.Contains(fun.Sel.Name, "Fatal") {
								pass.Reportf(callExpr.Pos(), "log.%s detected outside of main function", fun.Sel.Name)
							}
						}
						if pkgIdent, ok := fun.X.(*ast.Ident); ok && pkgIdent.Name == "os" {
							if fun.Sel.Name == "Exit" {
								pass.Reportf(callExpr.Pos(), "os.Exit detected outside of main function")
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
