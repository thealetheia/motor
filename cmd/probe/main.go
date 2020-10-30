package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"
	"unicode"
)

func main() {
	fset := token.NewFileSet()
	packs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		panic(err)
	}

	tests := []string{}
	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					name := fn.Name.Name
					if strings.HasPrefix(fn.Name.Name, "Test") {
						tests = append(tests, name)
					}
				}
			}
		}
	}

	needle, trace := "", false
	if len(os.Args) > 1 {
		for _, r := range os.Args[1] {
			if r == '!' {
				trace = true
			}
		}
		needle = strings.Trim(os.Args[1], "! ")
	}

	sort.Strings(tests)
	for _, test := range tests {
		if len(os.Args) == 1 {
			fmt.Println(test)
			continue
		}

		j := 0
		for i := range test {
			if unicode.ToLower(rune(needle[j])) == unicode.ToLower(rune(test[i])) {
				j++
				if j == len(needle) {
					args := strings.Join(os.Args[2:], " ")
					mode := "!debug"
					if trace {
						mode = "!trace"
					}
					got := exec.Command("go", "test", "-v", "-test", test, "-args", mode, args)
					got.Stdout = os.Stdout
					got.Stderr = os.Stderr
					if err := got.Run(); err != nil {
						if exiterr, ok := err.(*exec.ExitError); ok {
							if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
								os.Exit(status.ExitStatus())
							}
						}
					}
				}
			}
		}
	}
}
