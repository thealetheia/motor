package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"unicode"
)

func main() {
	retcode := 0
	defer func() { os.Exit(retcode) }()
	erred := func(args ...interface{}) {
		fmt.Printf("!! ")
		fmt.Println(args...)
		retcode = 1
		runtime.Goexit()
	}

	fset := token.NewFileSet()
	packs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		erred(err)
	}

	tests := []string{}
	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					name := fn.Name.Name
					if strings.HasPrefix(fn.Name.Name, "Benchmark") {
						tests = append(tests, name)
					}
				}
			}
		}
	}

	needle, trace := "", false
	if len(os.Args) > 1 {
		if os.Args[1] == "." {
			got := exec.Command("go", "test", "-bench=.", "-benchmem")
			got.Stdout = os.Stdout
			got.Stderr = os.Stderr
			if err := got.Run(); err != nil {
				erred(err)
			}
			return
		}
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
					mode := "!probe.debug"
					if trace {
						mode = "!probe.trace"
					}
					got := exec.Command("go", "test", "-bench", "-v", "-run", test, "-args", mode, args)
					got.Stdout = os.Stdout
					got.Stderr = os.Stderr
					if err := got.Run(); err != nil {
						erred(err)
					}
				}
			}
		}
	}
}
