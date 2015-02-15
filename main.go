package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var withTest = flag.Bool("test", false, "Include test for generated files.")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 || len(args) > 2 {
		badArgs()
	}
	sumTypeName := args[0]

	pkgs, err := parser.ParseDir(token.NewFileSet(), ".", nil, 0)
	exitErr("parsing package:", err)

	for _, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			ty, ok := lookupType(sumTypeName, file)
			if !ok {
				continue
			}
			if _, ok := ty.Type.(*ast.InterfaceType); !ok {
				exitErr("Type "+sumTypeName+" is not an interface.", errors.New(""))
			}

			a, err := generateSumWalker(pkg, fileName, ty)
			exitErr("generating sum walker:", err)

			var testAST *ast.File
			if *withTest {
				testAST, err = generateTest(pkg, fileName, ty)
				exitErr("generating test:", err)
			}

			outName := outFileName(fileName, false)
			out, err := os.Create(outName)
			exitErr("creating file:", err)

			err = printAST(out, a)
			exitErr("printing to file:", err)

			if *withTest {
				outName = outFileName(fileName, true)
				out, err = os.Create(outName)
				exitErr("creating file:", err)

				err = printAST(out, testAST)
				exitErr("printing to file:", err)
			}
		}
	}
}

func badArgs() {
	flag.Usage()
	os.Exit(1)
}

func exitErr(msg string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg, err)
		os.Exit(1)
	}
}

func lookupType(name string, file *ast.File) (*ast.TypeSpec, bool) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			tySpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if tySpec.Name.String() == name {
				return tySpec, true
			}
		}
	}
	return nil, false
}

func outFileName(name string, test bool) string {
	ps := strings.Split(name, ".")
	stest := ""
	if test {
		stest = "_test"
	}
	return strings.Join(ps[:len(ps)-1], "") + "_sumtype" + stest + ".go"
}
