package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const apiPackagePath = "github.com/wallet/api"

var (
	excludedFiles = []string{"FrpcApi.go", "FrpcPeriphery.go", "FrpcRequest.go", "FileServerRouter.go", "LnurlApi.go",
		"LnurlBoltPhoneStore.go",
		"LnurlBoltServerStore.go",
		"LnurlLndLightningServiceApi.go",
		"LnurlPhoneRouter.go",
		"LnurlPostRequest.go",
		"LnurlServerRouter.go",
		"LnurlUtil.go",
		"LndAddrStoreInBolt.go", "" +
			"ApiUtil.go"}
	excludedFunctions = []string{"MakeJsonResult", "MakeJsonErrorResult", "TapMarshalRespString", "B64DecodeToHex"}
)

type APIFunction struct {
	Name       string
	Parameters []*ast.Field
	Results    *ast.FieldList
	Recv       *ast.FieldList
}

func main() {
	functions, err := parseAPIPackage(apiPackagePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Loaded functions:")
	for _, fn := range functions {
		fmt.Printf("- %s\n", fn.Name)
	}

	err = generateHandlers(functions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Handler code generated successfully.")
}

func parseAPIPackage(packagePath string) ([]APIFunction, error) {
	log.Printf("Starting to parse package: %s", packagePath)

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedSyntax |
			packages.NeedTypesInfo | packages.NeedTypesSizes,
		Dir: ".",
	}

	log.Printf("Loading packages...")
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		log.Printf("Error loading packages: %v", err)
		return nil, err
	}

	log.Printf("Loaded %d packages", len(pkgs))
	for i, pkg := range pkgs {
		log.Printf("Package %d: %s", i, pkg.PkgPath)
		log.Printf("  Name: %s", pkg.Name)
		log.Printf("  Files: %v", pkg.GoFiles)
		log.Printf("  Syntax trees: %d", len(pkg.Syntax))
		log.Printf("  Errors: %v", pkg.Errors)
		if pkg.Fset == nil {
			log.Printf("  FileSet is nil")
		} else {
			log.Printf("  FileSet is not nil")
		}
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no package found")
	}

	var functions []APIFunction
	processedFuncs := make(map[string]bool)

	// 如果 packages.Load 没有提供足够的信息，尝试使用 go/parser
	if len(pkgs[0].Syntax) == 0 {
		log.Println("No syntax trees available, attempting to parse directly")
		fset := token.NewFileSet()
		if len(pkgs[0].GoFiles) == 0 {
			return nil, fmt.Errorf("no Go files found in package")
		}
		pkgDir := filepath.Dir(pkgs[0].GoFiles[0])
		parsedPkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing directory: %v", err)
			return nil, err
		}

		for _, pkg := range parsedPkgs {
			for fileName, file := range pkg.Files {
				log.Printf("Processing file: %s", fileName)
				if isExcludedFile(filepath.Base(fileName)) {
					log.Printf("Skipping excluded file: %s", fileName)
					continue
				}

				ast.Inspect(file, func(n ast.Node) bool {
					fn, ok := n.(*ast.FuncDecl)
					if !ok {
						return true
					}
					funcName := fn.Name.Name
					if shouldProcessFunc(funcName) && !processedFuncs[funcName] {
						log.Printf("Adding function: %s", funcName)
						functions = append(functions, APIFunction{
							Name:       funcName,
							Parameters: fn.Type.Params.List,
							Results:    fn.Type.Results,
							Recv:       fn.Recv,
						})
						processedFuncs[funcName] = true
					}
					return true
				})
			}
		}
	} else {
		// 原来的处理逻辑
		for _, pkg := range pkgs {
			if pkg == nil {
				log.Println("Skipping nil package")
				continue
			}

			if pkg.Fset == nil {
				log.Printf("FileSet is nil for package %s, skipping", pkg.PkgPath)
				continue
			}

			log.Printf("Processing package: %s", pkg.PkgPath)

			for _, file := range pkg.Syntax {
				if file == nil {
					log.Println("Nil syntax tree, skipping")
					continue
				}

				filePos := file.Pos()
				if !filePos.IsValid() {
					log.Println("Invalid position, skipping file")
					continue
				}

				tokenFile := pkg.Fset.File(filePos)
				if tokenFile == nil {
					log.Printf("No token file found for position %v, skipping file", filePos)
					continue
				}

				filename := filepath.Base(tokenFile.Name())
				log.Printf("Processing file: %s", filename)

				if isExcludedFile(filename) {
					log.Printf("Skipping excluded file: %s", filename)
					continue
				}

				ast.Inspect(file, func(n ast.Node) bool {
					fn, ok := n.(*ast.FuncDecl)
					if !ok {
						return true
					}
					funcName := fn.Name.Name
					if shouldProcessFunc(funcName) && !processedFuncs[funcName] {
						log.Printf("Adding function: %s", funcName)
						functions = append(functions, APIFunction{
							Name:       funcName,
							Parameters: fn.Type.Params.List,
							Results:    fn.Type.Results,
							Recv:       fn.Recv,
						})
						processedFuncs[funcName] = true
					}
					return true
				})
			}
		}
	}

	log.Printf("Parsed %d functions", len(functions))
	return functions, nil
}

func isExcludedFile(filename string) bool {
	for _, excludedFile := range excludedFiles {
		if filename == excludedFile {
			return true
		}
	}
	return false
}

func shouldProcessFunc(name string) bool {
	if !unicode.IsUpper(rune(name[0])) {
		return false
	}
	if isEmptyFunc(name) || isSpecialFunc(name) {
		return false
	}
	for _, excludedFunc := range excludedFunctions {
		if name == excludedFunc {
			return false
		}
	}
	return true
}

func isEmptyFunc(name string) bool {
	return false
}

func isSpecialFunc(name string) bool {
	return name == "MakeJsonResult" || name == "MakeJsonErrorResult"
}

func generateHandlers(functions []APIFunction) error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return err
	}
	fmt.Println("Current working directory:", cwd)

	cmdPath := filepath.Join(cwd, "./cmd/client")
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		err = os.Mkdir(cmdPath, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
	}

	filePath := filepath.Join(cmdPath, "generated_handlers.go")
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Fprintf(out, `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/wallet/api"
)

`)

	for _, fn := range functions {
		generateHandler(out, fn)
	}

	generateMain(out, functions)

	return nil
}

func generateHandler(out *os.File, fn APIFunction) {
	isMethod := fn.Recv != nil && len(fn.Recv.List) > 0

	fmt.Fprintf(out, "\nfunc handle%s(w http.ResponseWriter, r *http.Request) {\n", fn.Name)
	fmt.Fprintln(out, `    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    fmt.Printf("Received request for /api%s\n", r.URL.Path)

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }

    fmt.Printf("Received raw data for /api%s: %s\n", r.URL.Path, string(body))

    var params struct {`)

	if isMethod {
		receiverType := getTypeString(fn.Recv.List[0].Type)
		fmt.Fprintf(out, "        Receiver %s `json:\"receiver\"`\n", receiverType)
	}

	for _, param := range fn.Parameters {
		paramType := getTypeString(param.Type)
		for _, name := range param.Names {
			capitalizedName := strings.Title(name.Name)
			fmt.Fprintf(out, "        %s %s `json:\"%s\"`\n", capitalizedName, paramType, name.Name)
		}
	}

	fmt.Fprintln(out, `    }

    err = json.Unmarshal(body, &params)
    if err != nil {
        http.Error(w, "Error parsing JSON: " + err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)`)

	var paramNames []string
	if isMethod {
		paramNames = append(paramNames, "params.Receiver")
	}
	for i, param := range fn.Parameters {
		for _, name := range param.Names {
			paramName := "params." + strings.Title(name.Name)
			if i == len(fn.Parameters)-1 && isVariadic(param.Type) {
				paramName += "..."
			}
			paramNames = append(paramNames, paramName)
		}
	}

	if usesSpecialResultHandling(fn) {
		if isMethod {
			fmt.Fprintf(out, "    result := params.Receiver.%s(%s)\n", fn.Name, strings.Join(paramNames[1:], ", "))
		} else {
			fmt.Fprintf(out, "    result := api.%s(%s)\n", fn.Name, strings.Join(paramNames, ", "))
		}
		fmt.Fprintln(out, `    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(result))`)
	} else if fn.Results != nil && len(fn.Results.List) > 0 {
		fmt.Fprintf(out, "    ")
		for i := 0; i < len(fn.Results.List); i++ {
			if i > 0 {
				fmt.Fprintf(out, ", ")
			}
			fmt.Fprintf(out, "result%d", i)
		}
		if isMethod {
			fmt.Fprintf(out, " := params.Receiver.%s(%s)\n", fn.Name, strings.Join(paramNames[1:], ", "))
		} else {
			fmt.Fprintf(out, " := api.%s(%s)\n", fn.Name, strings.Join(paramNames, ", "))
		}

		fmt.Fprintln(out, `    w.Header().Set("Content-Type", "application/json")
    var response map[string]interface{}`)

		lastResultType := getTypeString(fn.Results.List[len(fn.Results.List)-1].Type)
		if lastResultType == "error" {
			fmt.Fprintf(out, `
    if result%d != nil {
        response = map[string]interface{}{
            "code":    500,
            "data":    nil,
            "error":   result%d.Error(),
            "success": false,
        }
    } else {
        response = map[string]interface{}{
            "code":    200,
            "data":    `, len(fn.Results.List)-1, len(fn.Results.List)-1)

			if len(fn.Results.List) > 1 {
				fmt.Fprintf(out, "result0")
			} else {
				fmt.Fprintf(out, "nil")
			}

			fmt.Fprintln(out, `,
            "error":   "",
            "success": true,
        }
    }`)
		} else {
			fmt.Fprintln(out, `
    response = map[string]interface{}{
        "code":    200,
        "data":    map[string]interface{}{`)
			for i := 0; i < len(fn.Results.List); i++ {
				fmt.Fprintf(out, `
            "result%d": result%d,`, i, i)
			}
			fmt.Fprintln(out, `
        },
        "error":   "",
        "success": true,
    }`)
		}

		fmt.Fprintln(out, `
    jsonResponse, _ := json.Marshal(response)
    w.Write(jsonResponse)`)
	} else {
		if isMethod {
			fmt.Fprintf(out, "    params.Receiver.%s(%s)\n", fn.Name, strings.Join(paramNames[1:], ", "))
		} else {
			fmt.Fprintf(out, "    api.%s(%s)\n", fn.Name, strings.Join(paramNames, ", "))
		}
		fmt.Fprintln(out, `    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "code":    200,
        "data":    nil,
        "error":   "",
        "success": true,
    }
    jsonResponse, _ := json.Marshal(response)
    w.Write(jsonResponse)`)
	}

	fmt.Fprintln(out, "}")
}

// 新添加的辅助函数
func isVariadic(expr ast.Expr) bool {
	_, ok := expr.(*ast.Ellipsis)
	return ok
}

func usesSpecialResultHandling(fn APIFunction) bool {
	return strings.Contains(fn.Name, "MakeJsonResult") || strings.Contains(fn.Name, "MakeJsonErrorResult")
}
func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string", "bool", "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
			"byte", "rune", "float32", "float64", "complex64", "complex128",
			"error", "any", "T":
			return t.Name
		default:
			// 检查是否是 proto 类型
			if strings.HasSuffix(t.Name, "pb") || strings.Contains(t.Name, "Proto") {
				return t.Name
			}
			return "api." + t.Name
		}
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			if ident.Name == "api" {
				return "api." + t.Sel.Name
			}
			return ident.Name + "." + t.Sel.Name
		}
		return getTypeString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + getTypeString(t.Elt)
	case *ast.MapType:
		return "map[" + getTypeString(t.Key) + "]" + getTypeString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "[]" + getTypeString(t.Elt) // 将可变参数转换为切片
	default:
		return "interface{}"
	}
}
func generateMain(out *os.File, functions []APIFunction) {
	fmt.Fprintln(out, `
func main() {
    mux := http.NewServeMux()`)

	for _, fn := range functions {
		fmt.Fprintf(out, "    mux.HandleFunc(\"/api/%s\", handle%s)\n", fn.Name, fn.Name)
	}

	fmt.Fprintln(out, `
    fmt.Println("Server is starting...")
    fmt.Println("Available endpoints:")`)

	for _, fn := range functions {
		fmt.Fprintf(out, "    fmt.Println(\"- http://localhost:7047/api/%s\")\n", fn.Name)
	}

	fmt.Fprintln(out, `    fmt.Println("Listening on :7047")
    http.ListenAndServe(":7047", mux)
}`)
}
