/**
 * Copyright 2026 AiT authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

// Command nestscan fails when production Go has control-flow nesting depth >= 3
// (if / for / range / switch / select), matching go-structure §9.
//
// Usage (from module root):
//
//	go run ./scripts/nestscan.go ./internal ./cmd ./pkg
//
// Skips *_test.go, mock/, and pkg/third-party/.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	roots := os.Args[1:]
	if len(roots) == 0 {
		fmt.Fprintln(os.Stderr, "usage: nestscan <dir>...")
		os.Exit(2)
	}

	fset := token.NewFileSet()
	var hits int
	for _, root := range roots {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			if strings.Contains(path, string(filepath.Separator)+"mock"+string(filepath.Separator)) ||
				strings.Contains(path, "third-party") {
				return nil
			}

			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, err)
				hits++
				return nil
			}
			for _, d := range f.Decls {
				fn, ok := d.(*ast.FuncDecl)
				if !ok || fn.Body == nil {
					continue
				}
				hits += scanNest(fset, path, fn.Name.Name, fn.Body, 0)
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk %s: %v\n", root, err)
			os.Exit(1)
		}
	}

	if hits > 0 {
		fmt.Fprintf(os.Stderr, "nestscan: %d control-flow nesting depth >= 3 site(s)\n", hits)
		os.Exit(1)
	}
}

func scanNest(fset *token.FileSet, path, fname string, n ast.Node, depth int) int {
	hits := 0
	ast.Inspect(n, func(node ast.Node) bool {
		if node == nil || node == n {
			return true
		}
		var body *ast.BlockStmt
		switch x := node.(type) {
		case *ast.IfStmt:
			body = x.Body
		case *ast.ForStmt:
			body = x.Body
		case *ast.RangeStmt:
			body = x.Body
		case *ast.SwitchStmt:
			body = x.Body
		case *ast.TypeSwitchStmt:
			body = x.Body
		case *ast.SelectStmt:
			body = x.Body
		case *ast.FuncLit:
			hits += scanNest(fset, path, fname+"$lit", x.Body, 0)
			return false
		default:
			return true
		}

		nd := depth + 1
		if nd >= 3 && body != nil {
			pos := fset.Position(node.Pos())
			fmt.Fprintf(os.Stderr, "%s:%d %s depth=%d\n", path, pos.Line, fname, nd)
			hits++
		}
		if body != nil {
			hits += scanNest(fset, path, fname, body, nd)
		}
		if ifs, ok := node.(*ast.IfStmt); ok && ifs.Else != nil {
			hits += scanNest(fset, path, fname, ifs.Else, nd)
		}
		return false
	})
	return hits
}
