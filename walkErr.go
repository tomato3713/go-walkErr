package walkErr

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"strings"

	"log/slog"
	"slices"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "walkErr",
	Doc:  "walkErr shows returned error list in function",
	Run:  run,
}

var errorType = types.Universe.Lookup("error").Type()
var errorInterface = errorType.Underlying().(*types.Interface)

// FuncError has fucntion declaration and error list which this function returns.
type FuncError struct {
	FuncDecl  *ast.FuncDecl
	Errors    []*ast.Ident
	CallFuncs []*ast.Ident
}

func run(pass *analysis.Pass) (interface{}, error) {
	programLevel := new(slog.LevelVar)
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(handler))
	logger := slog.New(handler)

	for _, f := range pass.Files {
		errList, err := Inspect(pass, f, logger)
		if err != nil {
			return nil, err
		}
		for _, v := range errList {
			msg := make([]string, 0, len(v.Errors))
			for _, item := range v.CallFuncs {
				logger.Info("search error", "func", fmt.Sprintf("%#v", item))
				obj := pass.TypesInfo.ObjectOf(item)
				if obj == nil {
					continue
				}
				name := fmt.Sprintf("%s.%s", obj.Pkg().Name(), obj.Name())
				if funcInfo, ok := errList[name]; ok {
					logger.Info("join error", "func", fmt.Sprintf("%#v", funcInfo))
					for _, item := range funcInfo.Errors {
						obj := pass.TypesInfo.ObjectOf(item)
						if obj == nil {
							continue
						}
						msg = append(msg, fmt.Sprintf("%s.%s", obj.Pkg().Name(), obj.Name()))
					}
				}
			}

			for _, item := range v.Errors {
				obj := pass.TypesInfo.ObjectOf(item)
				if obj == nil {
					continue
				}
				msg = append(msg, fmt.Sprintf("%s.%s", obj.Pkg().Name(), obj.Name()))
			}

			slices.Sort(msg)
			pass.Reportf(v.FuncDecl.Pos(), "return errors: %s", strings.Join(msg, ", "))
		}
	}

	return nil, nil
}

func InspectFunction(pass *analysis.Pass, funcDecl *ast.FuncDecl, logger *slog.Logger) ([]*ast.Ident, []*ast.Ident, error) {
	info := pass.TypesInfo
	errList := make([]*ast.Ident, 0, 0)
	funcList := make([]*ast.Ident, 0, 0)

	walk := func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ReturnStmt:
			for _, result := range n.Results {
				switch result := result.(type) {
				case *ast.Ident:
					obj := info.ObjectOf(result)
					if obj == nil {
						continue
					}
					if types.Implements(obj.Type(), errorInterface) {
						if !slices.ContainsFunc(errList, func(v *ast.Ident) bool {
							return v.Obj == result.Obj
						}) {
							errList = append(errList, result)
							logger.Info("find error", "name", fmt.Sprintf("%#v", result))
						}
					}
				case *ast.CallExpr:
					logger.Info("callExpr.Fun", "value", fmt.Sprintf("%#v", result.Fun))
					if fun, ok := result.Fun.(*ast.Ident); ok {
						funcList = append(funcList, fun)
					}
				}
			}
		default:
			return true
		}
		return false
	}

	for _, stmt := range funcDecl.Body.List {
		ast.Inspect(stmt, walk)
	}
	return errList, funcList, nil
}

func Inspect(pass *analysis.Pass, f *ast.File, logger *slog.Logger) (map[string]FuncError, error) {
	results := map[string]FuncError{}

	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if slices.ContainsFunc(funcDecl.Type.Results.List, func(v *ast.Field) bool {
				if ident, ok := v.Type.(*ast.Ident); ok {
					return ident.Name == "error"
				}
				return false
			}) {
				logger.Error("find function", "func", fmt.Sprintf("%#v", funcDecl))
				errList, callFuncs, err := InspectFunction(pass, funcDecl, logger)
				if err != nil {
					logger.Error("inspect function", "err", fmt.Sprintf("%v", err))
				}

				obj := pass.TypesInfo.ObjectOf(funcDecl.Name)
				if obj == nil {
					continue
				}
				name := fmt.Sprintf("%s.%s", obj.Pkg().Name(), obj.Name())
				results[name] = FuncError{
					FuncDecl:  funcDecl,
					Errors:    errList,
					CallFuncs: callFuncs,
				}
			}
		}
	}
	return results, nil
}
