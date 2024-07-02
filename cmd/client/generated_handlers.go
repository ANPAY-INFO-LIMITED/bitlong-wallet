package main

import (
	"encoding/json"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/wallet/api"
	"github.com/wallet/models"
	"io/ioutil"
	"net/http"
)

func handleGetApiVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetApiVersion()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleNewVersionTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.NewVersionTag()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSetPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Path    string `json:"path"`
		Network string `json:"network"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SetPath(params.Path, params.Network)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetPath()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleFileTestConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.FileTestConfig()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleReadConfigFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.ReadConfigFile()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleReadConfigFile1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.ReadConfigFile1()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleReadConfigFile2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.ReadConfigFile2()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCreateDir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.CreateDir()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCreateDir2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.CreateDir2()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleVisit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.Visit()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCreateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Dir      string `json:"dir"`
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CreateFile(params.Dir, params.Filename, params.Content)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleReadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		FilePath string `json:"filePath"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ReadFile(params.FilePath)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCopyFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		SrcPath  string `json:"srcPath"`
		DestPath string `json:"destPath"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CopyFile(params.SrcPath, params.DestPath)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		FilePath string `json:"filePath"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.DeleteFile(params.FilePath)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGenerateKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Mnemonic string `json:"mnemonic"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GenerateKeys(params.Mnemonic)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetPublicKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetPublicKey()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNPublicKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetNPublicKey()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetJsonPublicKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetJsonPublicKey()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSignMess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Message string `json:"message"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SignMess(params.Message)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleRouterForKeyService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.RouterForKeyService()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleLitdStopDaemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.LitdStopDaemon()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleLitdLocalStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.LitdLocalStop()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSubServerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SubServerStatus()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTapdStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetTapdStatus()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetLitStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetLitStatus()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNewAddress_P2TR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetNewAddress_P2TR()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNewAddress_P2WKH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetNewAddress_P2WKH()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNewAddress_NP2WKH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetNewAddress_NP2WKH()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleStoreAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Name           string `json:"name"`
		Address        string `json:"address"`
		Balance        int    `json:"balance"`
		AddressType    string `json:"addressType"`
		DerivationPath string `json:"derivationPath"`
		IsInternal     bool   `json:"isInternal"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.StoreAddr(params.Name, params.Address, params.Balance, params.AddressType, params.DerivationPath, params.IsInternal)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleRemoveAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.RemoveAddr(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.QueryAddr(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryAllAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.QueryAllAddr()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNonZeroBalanceAddresses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetNonZeroBalanceAddresses()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleUpdateAllAddressesByGNZBA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.UpdateAllAddressesByGNZBA()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAccountsString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAccountsString()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAccounts()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAddressTypeToDerivationPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AddressType string `json:"addressType"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AddressTypeToDerivationPath(params.AddressType)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetPathByAddressType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AddressType string `json:"addressType"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetPathByAddressType(params.AddressType)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockWrap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		BlockHash string `json:"blockHash"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetBlockWrap(params.BlockHash)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockInfoByHeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Height int64 `json:"height"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetBlockInfoByHeight(params.Height)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetWalletBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetWalletBalance()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessGetWalletBalanceResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		WalletBalanceResponse *lnrpc.WalletBalanceResponse `json:"walletBalanceResponse"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ProcessGetWalletBalanceResult(params.WalletBalanceResponse)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCalculateImportedTapAddressBalanceAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ListAddressesResponse *walletrpc.ListAddressesResponse `json:"listAddressesResponse"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CalculateImportedTapAddressBalanceAmount(params.ListAddressesResponse)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetInfoOfLnd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetInfoOfLnd()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetIdentityPubkey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetIdentityPubkey()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNewAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetNewAddress()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAddInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Value int64  `json:"value"`
		Memo  string `json:"memo"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AddInvoice(params.Value, params.Memo)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListInvoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListInvoices()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSimplifyInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Invoice *lnrpc.ListInvoiceResponse `json:"invoice"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SimplifyInvoice(params.Invoice)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleLookupInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Rhash string `json:"rhash"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.LookupInvoice(params.Rhash)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAbandonChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AbandonChannel()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleBatchOpenChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.BatchOpenChannel()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleChannelAcceptor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ChannelAcceptor()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleChannelBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ChannelBalance()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCheckMacaroonPermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CheckMacaroonPermissions()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCloseChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		FundingTxidStr string `json:"fundingTxidStr"`
		OutputIndex    int    `json:"outputIndex"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CloseChannel(params.FundingTxidStr, params.OutputIndex)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleClosedChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ClosedChannels()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDecodePayReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		PayReq string `json:"payReq"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.DecodePayReq(params.PayReq)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleExportAllChannelBackups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ExportAllChannelBackups()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleExportChannelBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ExportChannelBackup()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChanInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ChanId string `json:"chanId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetChanInfo(params.ChanId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleOpenChannelSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		NodePubkey         string `json:"nodePubkey"`
		LocalFundingAmount int64  `json:"localFundingAmount"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.OpenChannelSync(params.NodePubkey, params.LocalFundingAmount)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleOpenChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		NodePubkey         string `json:"nodePubkey"`
		LocalFundingAmount int64  `json:"localFundingAmount"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.OpenChannel(params.NodePubkey, params.LocalFundingAmount)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListChannels()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handlePendingChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.PendingChannels()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChannelState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ChanPoint string `json:"chanPoint"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetChannelState(params.ChanPoint)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChannelInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ChanPoint string `json:"chanPoint"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetChannelInfo(params.ChanPoint)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleRestoreChannelBackups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.RestoreChannelBackups()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSubscribeChannelBackups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SubscribeChannelBackups()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSubscribeChannelEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SubscribeChannelEvents()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSubscribeChannelGraph(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SubscribeChannelGraph()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleUpdateChannelPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.UpdateChannelPolicy()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleVerifyChanBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.VerifyChanBackup()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleConnectPeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Pubkey string `json:"pubkey"`
		Host   string `json:"host"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ConnectPeer(params.Pubkey, params.Host)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleEstimateFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Addr   string `json:"addr"`
		Amount int64  `json:"amount"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.EstimateFee(params.Addr, params.Amount)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendPaymentSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Invoice string `json:"invoice"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SendPaymentSync(params.Invoice)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendPaymentSync0amt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Invoice string `json:"invoice"`
		Amt     int64  `json:"amt"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SendPaymentSync0amt(params.Invoice, params.Amt)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendCoins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Addr    string `json:"addr"`
		Amount  int64  `json:"amount"`
		FeeRate int64  `json:"feeRate"`
		SendAll bool   `json:"sendAll"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SendCoins(params.Addr, params.Amount, params.FeeRate, params.SendAll)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendMany(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		JsonAddr string `json:"jsonAddr"`
		FeeRate  int64  `json:"feeRate"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SendMany(params.JsonAddr, params.FeeRate)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendAllCoins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Addr string `json:"addr"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SendAllCoins(params.Addr)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleLndStopDaemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.LndStopDaemon()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListPermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListPermissions()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendPaymentV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Invoice  string `json:"invoice"`
		Feelimit int64  `json:"feelimit"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SendPaymentV2(params.Invoice, params.Feelimit)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleTrackPaymentV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Payhash string `json:"payhash"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.TrackPaymentV2(params.Payhash)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendToRouteV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Payhash []byte       `json:"payhash"`
		Route   *lnrpc.Route `json:"route"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SendToRouteV2(params.Payhash, params.Route)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleEstimateRouteFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Dest   string `json:"dest"`
		Amtsat int64  `json:"amtsat"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.EstimateRouteFee(params.Dest, params.Amtsat)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetStateForSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetStateForSubscribe()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetState()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGenSeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GenSeed()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleInitWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Seed     string `json:"seed"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.InitWallet(params.Seed, params.Password)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleUnlockWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.UnlockWallet(params.Password)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ChangePassword(params.CurrentPassword, params.NewPassword)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAddresses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListAddresses()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAddressesAndGetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ListAddressesAndGetResponse()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListAccounts()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleFindAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.FindAccount(params.Name)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListLeases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListLeases()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListSweeps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListSweeps()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListUnspent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListUnspent()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleNextAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.NextAddr()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSetServerHost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Server string `json:"server"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SetServerHost(params.Server)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetServerHost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetServerHost()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.Login(params.Username, params.Password)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.Refresh(params.Username, params.Password)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendPostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Url         string `json:"url"`
		Token       string `json:"token"`
		RequestBody []byte `json:"requestBody"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.SendPostRequest(params.Url, params.Token, params.RequestBody)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSimplifyTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address   string                              `json:"address"`
		Responses *api.GetAddressTransactionsResponse `json:"responses"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SimplifyTransactions(params.Address, params.Responses)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressInfoByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAddressInfoByMempool(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAddressTransactions(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressTransferOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAddressTransferOut(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressTransferOutResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAddressTransferOutResult(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressTransactionsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAddressTransactionsByMempool(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressTransactionsChainByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetAddressTransactionsChainByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressTransactionsMempoolByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetAddressTransactionsMempoolByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressUTXOByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetAddressUTXOByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAddressValidationByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetAddressValidationByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockByMempoolByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockByMempoolByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockHeaderByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockHeaderByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockHeightByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockHeightByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockTimestampByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockTimestampByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockRawByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockRawByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockStatusByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockStatusByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockTipHeightByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetBlockTipHeightByMempool()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleBlockTipHeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.BlockTipHeight()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockTipHashByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockTipHashByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockTransactionIDByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockTransactionIDByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockTransactionIDsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockTransactionIDsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockTransactionsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockTransactionsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlocksByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlocksByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlocksBulkByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlocksBulkByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMempoolBlocksFeesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMempoolBlocksFeesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetRecommendedFeesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetRecommendedFeesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetDifficultyAdjustmentByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetDifficultyAdjustmentByMempool()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNetworkStatsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetNetworkStatsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNodesSlashChannelsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetNodesSlashChannelsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNodesInCountryByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetNodesInCountryByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNodeStatsPerCountryByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetNodeStatsPerCountryByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetISPNodesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetISPNodesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNodeStatsPerISPByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetNodeStatsPerISPByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTop100NodesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTop100NodesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTop100NodesbyLiquidityByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTop100NodesbyLiquidityByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTop100NodesbyConnectivityByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTop100NodesbyConnectivityByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTop100OldestNodesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTop100OldestNodesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetNodeStatsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetNodeStatsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetHistoricalNodeStatsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetHistoricalNodeStatsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChannelByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetChannelByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChannelsfromTXIDByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetChannelsfromTXIDByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChannelsfromNodePubkeyByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetChannelsfromNodePubkeyByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChannelGeodataByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetChannelGeodataByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChannelGeodataforNodeByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetChannelGeodataforNodeByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMempoolByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMempoolByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMempoolTransactionIDsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMempoolTransactionIDsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMempoolRecentByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMempoolRecentByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMempoolRBFTransactionsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMempoolRBFTransactionsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMempoolFullRBFTransactionsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMempoolFullRBFTransactionsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMiningPoolsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMiningPoolsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMiningPoolByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMiningPoolByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMiningPoolHashratesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMiningPoolHashratesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMiningPoolHashrateByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMiningPoolHashrateByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMiningPoolBlocksByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetMiningPoolBlocksByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetHashrateByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetHashrateByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetDifficultyAdjustmentsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetDifficultyAdjustmentsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetRewardStatsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetRewardStatsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockFeesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockFeesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockRewardsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockRewardsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockFeeratesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockFeeratesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockSizesandWeightsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockSizesandWeightsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockPredictionsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockPredictionsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockAuditScoreByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockAuditScoreByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlocksAuditScoresByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlocksAuditScoresByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetBlockAuditSummaryByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetBlockAuditSummaryByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetChildrenPayforParentByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetChildrenPayforParentByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Transaction string `json:"transaction"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetTransactionByMempool(params.Transaction)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionHexByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionHexByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionMerkleblockProofByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionMerkleblockProofByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionMerkleProofByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionMerkleProofByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionOutspendByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionOutspendByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionOutspendsByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionOutspendsByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionRawByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionRawByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionRBFHistoryByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionRBFHistoryByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionStatusByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionStatusByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionTimesByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.GetTransactionTimesByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handlePostTransactionByMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.PostTransactionByMempool()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleStartLitd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.StartLitd()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleStartLnd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.StartLnd()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleStartTapd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.StartTapd()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetUserOwnIssuanceHistoryInfos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetUserOwnIssuanceHistoryInfos(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetIssuanceTransactionFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetIssuanceTransactionFee(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMintTransactionFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token  string `json:"token"`
		Id     int    `json:"id"`
		Number int    `json:"number"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetMintTransactionFee(params.Token, params.Id, params.Number)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetLocalIssuanceTransactionFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		FeeRate int `json:"feeRate"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetLocalIssuanceTransactionFee(params.FeeRate)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetLocalIssuanceTransactionByteSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetLocalIssuanceTransactionByteSize()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetIssuanceTransactionCalculatedFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetIssuanceTransactionCalculatedFee(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMintTransactionCalculatedFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token  string `json:"token"`
		Id     int    `json:"id"`
		Number int    `json:"number"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetMintTransactionCalculatedFee(params.Token, params.Id, params.Number)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetIssuanceTransactionByteSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetIssuanceTransactionByteSize()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTapdMintAssetAndFinalizeTransactionByteSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetTapdMintAssetAndFinalizeTransactionByteSize()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTapdSendReservedAssetTransactionByteSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetTapdSendReservedAssetTransactionByteSize()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMintTransactionByteSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetMintTransactionByteSize()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetServerOwnSetFairLaunchInfos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetServerOwnSetFairLaunchInfos(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessOwnSetFairLaunchResponseToIssuanceHistoryInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		FairLaunchInfos *[]models.FairLaunchInfo `json:"fairLaunchInfos"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ProcessOwnSetFairLaunchResponseToIssuanceHistoryInfo(params.FairLaunchInfos)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetServerFeeRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetServerFeeRate(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetServerQueryMint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token  string `json:"token"`
		Id     int    `json:"id"`
		Number int    `json:"number"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetServerQueryMint(params.Token, params.Id, params.Number)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetServerIssuanceHistoryInfos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetServerIssuanceHistoryInfos(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetLocalTapdIssuanceHistoryInfos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetLocalTapdIssuanceHistoryInfos()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTimestampByBatchTxidWithGetTransactionsResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		TransactionDetails *lnrpc.TransactionDetails `json:"transactionDetails"`
		BatchTxid          string                    `json:"batchTxid"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetTimestampByBatchTxidWithGetTransactionsResponse(params.TransactionDetails, params.BatchTxid)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionByBatchTxid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		TransactionDetails *lnrpc.TransactionDetails `json:"transactionDetails"`
		BatchTxid          string                    `json:"batchTxid"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetTransactionByBatchTxid(params.TransactionDetails, params.BatchTxid)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetIdByBatchTxidWithListAssetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ListAssetResponse *taprpc.ListAssetResponse `json:"listAssetResponse"`
		BatchTxid         string                    `json:"batchTxid"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAssetIdByBatchTxidWithListAssetResponse(params.ListAssetResponse, params.BatchTxid)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetIdByOutpointAndNameWithListAssetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ListAssetResponse *taprpc.ListAssetResponse `json:"listAssetResponse"`
		Outpoint          string                    `json:"outpoint"`
		Name              string                    `json:"name"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAssetIdByOutpointAndNameWithListAssetResponse(params.ListAssetResponse, params.Outpoint, params.Name)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetsByOutpointWithListAssetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ListAssetResponse *taprpc.ListAssetResponse `json:"listAssetResponse"`
		Outpoint          string                    `json:"outpoint"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAssetsByOutpointWithListAssetResponse(params.ListAssetResponse, params.Outpoint)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetImageByImageData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ImageData string `json:"imageData"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetImageByImageData(params.ImageData)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetOwnSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetOwnSet(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetRate(params.Token)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetQueryMint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token            string `json:"token"`
		FairLaunchInfoId string `json:"FairLaunchInfoId"`
		MintedNumber     int    `json:"MintedNumber"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAssetQueryMint(params.Token, params.FairLaunchInfoId, params.MintedNumber)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendGetReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Url         string `json:"url"`
		Token       string `json:"token"`
		RequestBody []byte `json:"requestBody"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.SendGetReq(params.Url, params.Token, params.RequestBody)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAnchorVirtualPsbts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		VirtualPsbts []string `json:"virtualPsbts"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AnchorVirtualPsbts(params.VirtualPsbts)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleFundVirtualPsbt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		IsPsbtNotRaw bool     `json:"isPsbtNotRaw"`
		Psbt         []string `json:"psbt"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.FundVirtualPsbt(params.IsPsbtNotRaw, params.Psbt...)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleNextInternalKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		KeyFamily int `json:"keyFamily"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.NextInternalKey(params.KeyFamily)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleNextScriptKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		KeyFamily int `json:"keyFamily"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.NextScriptKey(params.KeyFamily)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProveAssetOwnership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId   string `json:"assetId"`
		ScriptKey string `json:"scriptKey"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProveAssetOwnership(params.AssetId, params.ScriptKey)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleRemoveUTXOLease(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.RemoveUTXOLease()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSignVirtualPsbt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		FundedPsbt string `json:"fundedPsbt"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SignVirtualPsbt(params.FundedPsbt)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleVerifyAssetOwnership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ProofWithWitness string `json:"proofWithWitness"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.VerifyAssetOwnership(params.ProofWithWitness)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSimplifyAssetsTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SimplifyAssetsTransfer()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSimplifyAssetsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Assets []*taprpc.Asset `json:"assets"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SimplifyAssetsList(params.Assets)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncUniverseFullSpecified(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		UniverseHost string `json:"universeHost"`
		Id           string `json:"id"`
		ProofType    string `json:"proofType"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncUniverseFullSpecified(params.UniverseHost, params.Id, params.ProofType)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncAssetIssuance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncAssetIssuance(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncAssetTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncAssetTransfer(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncAssetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SyncAssetAll(params.Id)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncAssetAllSlice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Ids []string `json:"ids"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SyncAssetAllSlice(params.Ids)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncAssetAllWithAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Ids []string `json:"ids"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SyncAssetAllWithAssets(params.Ids...)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetBalances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetBalances()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetGroupBalances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetGroupBalances()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetIdByAssetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetBalance *[]api.AssetBalance `json:"assetBalance"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetIdByAssetBalance(params.AssetBalance)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncAllAssetsByAssetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncAllAssetsByAssetBalance()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetsIdSlice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetsIdSlice()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetKeysTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AssetKeysTransfer(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeavesSpecified(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id        string `json:"id"`
		ProofType string `json:"proofType"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AssetLeavesSpecified(params.Id, params.ProofType)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessAssetTransferLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Response *universerpc.AssetLeafResponse `json:"response"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessAssetTransferLeave(params.Response)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeavesTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AssetLeavesTransfer(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeavesTransfer_ONLY_FOR_TEST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AssetLeavesTransfer_ONLY_FOR_TEST(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessAssetIssuanceLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Response *universerpc.AssetLeafResponse `json:"response"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessAssetIssuanceLeave(params.Response)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetInfoByIssuanceLeaf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetInfoByIssuanceLeaf(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDecodeRawProofByte(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		RawProof []byte `json:"rawProof"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.DecodeRawProofByte(params.RawProof)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDecodeRawProofString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Proof string `json:"proof"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.DecodeRawProofString(params.Proof)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Response *taprpc.DecodeProofResponse `json:"response"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessProof(params.Response)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDecodeRawProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Proof string `json:"proof"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.DecodeRawProof(params.Proof)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessListAllAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Response *taprpc.ListAssetResponse `json:"response"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessListAllAssets(params.Response)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetList()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessListAllAssetsSimplified(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Result *[]api.ListAllAsset `json:"result"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessListAllAssetsSimplified(params.Result)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetListSimplified(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetListSimplified()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetIdByListAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetIdByListAll()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncUniverseFullIssuanceByIdSlice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Ids []string `json:"ids"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncUniverseFullIssuanceByIdSlice(params.Ids)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncUniverseFullTransferByIdSlice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Ids []string `json:"ids"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncUniverseFullTransferByIdSlice(params.Ids)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncUniverseFullNoSlice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncUniverseFullNoSlice()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleOutpointToAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Outpoint string `json:"outpoint"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.OutpointToAddress(params.Outpoint)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleTransactionAndIndexToAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Transaction string `json:"transaction"`
		IndexStr    string `json:"indexStr"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.TransactionAndIndexToAddress(params.Transaction, params.IndexStr)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleTransactionAndIndexToValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Transaction string `json:"transaction"`
		IndexStr    string `json:"indexStr"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.TransactionAndIndexToValue(params.Transaction, params.IndexStr)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCompareScriptKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ScriptKey1 string `json:"scriptKey1"`
		ScriptKey2 string `json:"scriptKey2"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CompareScriptKey(params.ScriptKey1, params.ScriptKey2)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetHoldInfosIncludeSpent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetHoldInfosIncludeSpent(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetHoldInfosExcludeSpent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetHoldInfosExcludeSpent(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetHoldInfosIncludeSpentSlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetHoldInfosIncludeSpentSlow(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAddressIsSpent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AddressIsSpent(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAddressIsSpentAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AddressIsSpentAll(params.Address)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleOutpointToTransactionInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Outpoint string `json:"outpoint"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.OutpointToTransactionInfo(params.Outpoint)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetTransactionInfos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetTransactionInfos(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncAllAssetByList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncAllAssetByList()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetInfoById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetInfoById(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetHoldInfosExcludeSpentSlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetHoldInfosExcludeSpentSlow(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetTransactionInfoSlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetTransactionInfoSlow(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetIDAndTransferScriptKeyToOutpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id        string `json:"id"`
		ScriptKey string `json:"scriptKey"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AssetIDAndTransferScriptKeyToOutpoint(params.Id, params.ScriptKey)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAllAssetListWithoutProcession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAllAssetListWithoutProcession()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListBatchesAndGetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ListBatchesAndGetResponse()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionsAndGetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetTransactionsAndGetResponse()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionsExcludeLabelTapdAssetMinting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetTransactionsExcludeLabelTapdAssetMinting()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleExcludeLabelIsTapdAssetMinting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Response *lnrpc.TransactionDetails `json:"response"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ExcludeLabelIsTapdAssetMinting(params.Response)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAssetAndGetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ListAssetAndGetResponse()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAssetAndGetResponseByFlags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		WithWitness   bool `json:"withWitness"`
		IncludeSpent  bool `json:"includeSpent"`
		IncludeLeased bool `json:"includeLeased"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ListAssetAndGetResponseByFlags(params.WithWitness, params.IncludeSpent, params.IncludeLeased)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListBatchesAndGetCustomResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ListBatchesAndGetCustomResponse()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAssetAndGetCustomResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ListAssetAndGetCustomResponse()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionsAndGetCustomResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetTransactionsAndGetCustomResponse()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeafKeysIssuance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.AssetLeafKeysIssuance(params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeavesIssuance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.AssetLeavesIssuance(params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionsWhoseLabelIsTapdAssetMinting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetTransactionsWhoseLabelIsTapdAssetMinting()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetTransactionsWhoseLabelIsNotTapdAssetMinting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetTransactionsWhoseLabelIsNotTapdAssetMinting()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDecodeTransactionsWhoseLabelIsNotTapdAssetMinting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		RawTransactions []string `json:"rawTransactions"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.DecodeTransactionsWhoseLabelIsNotTapdAssetMinting(params.RawTransactions)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleRawTransactionHexSliceToRequestBodyRawString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		RawTransactions []string `json:"rawTransactions"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.RawTransactionHexSliceToRequestBodyRawString(params.RawTransactions)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handlePostCallBitcoindToDecodeRawTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Token           string   `json:"token"`
		RawTransactions []string `json:"rawTransactions"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.PostCallBitcoindToDecodeRawTransaction(params.Token, params.RawTransactions)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessDecodedTransactionsData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		DecodedRawTransactions *[]api.PostDecodeRawTransactionResponse `json:"decodedRawTransactions"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessDecodedTransactionsData(params.DecodedRawTransactions)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAndDecodeTransactionsWhoseLabelIsNotTapdAssetMinting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.GetAndDecodeTransactionsWhoseLabelIsNotTapdAssetMinting()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCancelBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CancelBatch()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleFinalizeBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		FeeRate int `json:"feeRate"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.FinalizeBatch(params.FeeRate)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListBatches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListBatches()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleMintAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Name                   string    `json:"name"`
		AssetTypeIsCollectible bool      `json:"assetTypeIsCollectible"`
		AssetMetaData          *api.Meta `json:"assetMetaData"`
		Amount                 int       `json:"amount"`
		NewGroupedAsset        bool      `json:"newGroupedAsset"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.MintAsset(params.Name, params.AssetTypeIsCollectible, params.AssetMetaData, params.Amount, params.NewGroupedAsset)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAddGroupAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Name                   string    `json:"name"`
		AssetTypeIsCollectible bool      `json:"assetTypeIsCollectible"`
		AssetMetaData          *api.Meta `json:"assetMetaData"`
		Amount                 int       `json:"amount"`
		GroupKey               string    `json:"groupKey"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AddGroupAsset(params.Name, params.AssetTypeIsCollectible, params.AssetMetaData, params.Amount, params.GroupKey)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleNewMeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Description string `json:"description"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.NewMeta(params.Description)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleLoadImageByByte(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Receiver *api.Meta `json:"receiver"`
		Image    []byte    `json:"image"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := params.Receiver.LoadImageByByte(params.Image)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleLoadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Receiver *api.Meta `json:"receiver"`
		File     string    `json:"file"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := params.Receiver.LoadImage(params.File)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleToJsonStr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Receiver *api.Meta `json:"receiver"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := params.Receiver.ToJsonStr()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetMetaFromStr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Receiver *api.Meta `json:"receiver"`
		MetaStr  string    `json:"metaStr"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	params.Receiver.GetMetaFromStr(params.MetaStr)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSaveImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Receiver *api.Meta `json:"receiver"`
		Dir      string    `json:"dir"`
		Name     string    `json:"name"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := params.Receiver.SaveImage(params.Dir, params.Name)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Receiver *api.Meta `json:"receiver"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := params.Receiver.GetImage()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleFetchAssetMeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Receiver *api.Meta `json:"receiver"`
		IsHash   bool      `json:"isHash"`
		Data     string    `json:"data"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := params.Receiver.FetchAssetMeta(params.IsHash, params.Data)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleImportProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		ProofFile    string `json:"proofFile"`
		GenesisPoint string `json:"genesisPoint"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ImportProof(params.ProofFile, params.GenesisPoint)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAddrReceives(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AddrReceives(params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleBurnAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetIdStr   string `json:"AssetIdStr"`
		AmountToBurn int64  `json:"amountToBurn"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.BurnAsset(params.AssetIdStr, params.AmountToBurn)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDebugLevel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.DebugLevel()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDecodeAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Addr string `json:"addr"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.DecodeAddr(params.Addr)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDecodeProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		RawProof string `json:"rawProof"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.DecodeProof(params.RawProof)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleExportProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.ExportProof()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetInfoOfTap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetInfoOfTap()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		WithWitness   bool `json:"withWitness"`
		IncludeSpent  bool `json:"includeSpent"`
		IncludeLeased bool `json:"includeLeased"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListAssets(params.WithWitness, params.IncludeSpent, params.IncludeLeased)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListSimpleAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		WithWitness   bool `json:"withWitness"`
		IncludeSpent  bool `json:"includeSpent"`
		IncludeLeased bool `json:"includeLeased"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListSimpleAssets(params.WithWitness, params.IncludeSpent, params.IncludeLeased)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleFindAssetByAssetName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetName string `json:"assetName"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.FindAssetByAssetName(params.AssetName)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListGroups()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryAssetTransfers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.QueryAssetTransfers(params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListUtxos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		IncludeLeased bool `json:"includeLeased"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListUtxos(params.IncludeLeased)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleNewAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
		Amt     int    `json:"amt"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.NewAddr(params.AssetId, params.Amt)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryAddrs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.QueryAddrs(params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSendAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		JsonAddrs string `json:"jsonAddrs"`
		FeeRate   int64  `json:"feeRate"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SendAssets(params.JsonAddrs, params.FeeRate)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSubscribeReceiveAssetEventNtfns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SubscribeReceiveAssetEventNtfns()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSubscribeSendAssetEventNtfns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SubscribeSendAssetEventNtfns()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleVerifyProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.VerifyProof()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleTapStopDaemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.TapStopDaemon()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessListBalancesResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Response *taprpc.ListBalancesResponse `json:"response"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessListBalancesResponse(params.Response)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleProcessListBalancesByGroupKeyResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Response *taprpc.ListBalancesResponse `json:"response"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ProcessListBalancesByGroupKeyResponse(params.Response)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListBalances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListBalances()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListBalancesByGroupKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListBalancesByGroupKey()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleCheckAssetIssuanceIsLocal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.CheckAssetIssuanceIsLocal(params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAssetsProcessed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		WithWitness   bool `json:"withWitness"`
		IncludeSpent  bool `json:"includeSpent"`
		IncludeLeased bool `json:"includeLeased"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.ListAssetsProcessed(params.WithWitness, params.IncludeSpent, params.IncludeLeased)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListAssetsAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListAssetsAll()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListNFTGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListNFTGroups()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListNFTAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListNFTAssets()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryAllNFTByGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.QueryAllNFTByGroup()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAddFederationServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.AddFederationServer()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeafKeysAndGetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId   string                `json:"assetId"`
		ProofType universerpc.ProofType `json:"proofType"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.AssetLeafKeysAndGetResponse(params.AssetId, params.ProofType)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeafKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id        string `json:"id"`
		ProofType string `json:"proofType"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AssetLeafKeys(params.Id, params.ProofType)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeaves(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.AssetLeaves(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleGetAssetInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.GetAssetInfo(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetRoots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.AssetRoots()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDeleteAssetRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.DeleteAssetRoot()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleDeleteFederationServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.DeleteFederationServer()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleUniverseInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.UniverseInfo()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleInsertProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.InsertProof()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleListFederationServers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.ListFederationServers()
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleMultiverseRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.MultiverseRoot()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryAssetRoots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.QueryAssetRoots(params.Id)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryAssetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		AssetId string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.QueryAssetStats(params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.QueryEvents()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryFederationSyncConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.QueryFederationSyncConfig()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleQueryProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.QueryProof()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSetFederationSyncConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.SetFederationSyncConfig()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleSyncUniverse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		UniverseHost string `json:"universeHost"`
		AssetId      string `json:"assetId"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0 := api.SyncUniverse(params.UniverseHost, params.AssetId)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	response = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{

			"result0": result0,
		},
		"error":   "",
		"success": true,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleUniverseStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	api.UniverseStats()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    200,
		"data":    nil,
		"error":   "",
		"success": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func handleAssetLeavesAndGetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var params struct {
		IsGroup   bool                  `json:"isGroup"`
		Id        string                `json:"id"`
		ProofType universerpc.ProofType `json:"proofType"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed parameters for /api%s: %+v\n", r.URL.Path, params)
	result0, result1 := api.AssetLeavesAndGetResponse(params.IsGroup, params.Id, params.ProofType)
	w.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if result1 != nil {
		response = map[string]interface{}{
			"code":    500,
			"data":    nil,
			"error":   result1.Error(),
			"success": false,
		}
	} else {
		response = map[string]interface{}{
			"code":    200,
			"data":    result0,
			"error":   "",
			"success": true,
		}
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/GetApiVersion", handleGetApiVersion)
	mux.HandleFunc("/api/NewVersionTag", handleNewVersionTag)
	mux.HandleFunc("/api/SetPath", handleSetPath)
	mux.HandleFunc("/api/GetPath", handleGetPath)
	mux.HandleFunc("/api/FileTestConfig", handleFileTestConfig)
	mux.HandleFunc("/api/ReadConfigFile", handleReadConfigFile)
	mux.HandleFunc("/api/ReadConfigFile1", handleReadConfigFile1)
	mux.HandleFunc("/api/ReadConfigFile2", handleReadConfigFile2)
	mux.HandleFunc("/api/CreateDir", handleCreateDir)
	mux.HandleFunc("/api/CreateDir2", handleCreateDir2)
	mux.HandleFunc("/api/Visit", handleVisit)
	mux.HandleFunc("/api/CreateFile", handleCreateFile)
	mux.HandleFunc("/api/ReadFile", handleReadFile)
	mux.HandleFunc("/api/CopyFile", handleCopyFile)
	mux.HandleFunc("/api/DeleteFile", handleDeleteFile)
	mux.HandleFunc("/api/GenerateKeys", handleGenerateKeys)
	mux.HandleFunc("/api/GetPublicKey", handleGetPublicKey)
	mux.HandleFunc("/api/GetNPublicKey", handleGetNPublicKey)
	mux.HandleFunc("/api/GetJsonPublicKey", handleGetJsonPublicKey)
	mux.HandleFunc("/api/SignMess", handleSignMess)
	mux.HandleFunc("/api/RouterForKeyService", handleRouterForKeyService)
	mux.HandleFunc("/api/LitdStopDaemon", handleLitdStopDaemon)
	mux.HandleFunc("/api/LitdLocalStop", handleLitdLocalStop)
	mux.HandleFunc("/api/SubServerStatus", handleSubServerStatus)
	mux.HandleFunc("/api/GetTapdStatus", handleGetTapdStatus)
	mux.HandleFunc("/api/GetLitStatus", handleGetLitStatus)
	mux.HandleFunc("/api/GetNewAddress_P2TR", handleGetNewAddress_P2TR)
	mux.HandleFunc("/api/GetNewAddress_P2WKH", handleGetNewAddress_P2WKH)
	mux.HandleFunc("/api/GetNewAddress_NP2WKH", handleGetNewAddress_NP2WKH)
	mux.HandleFunc("/api/StoreAddr", handleStoreAddr)
	mux.HandleFunc("/api/RemoveAddr", handleRemoveAddr)
	mux.HandleFunc("/api/QueryAddr", handleQueryAddr)
	mux.HandleFunc("/api/QueryAllAddr", handleQueryAllAddr)
	mux.HandleFunc("/api/GetNonZeroBalanceAddresses", handleGetNonZeroBalanceAddresses)
	mux.HandleFunc("/api/UpdateAllAddressesByGNZBA", handleUpdateAllAddressesByGNZBA)
	mux.HandleFunc("/api/GetAllAccountsString", handleGetAllAccountsString)
	mux.HandleFunc("/api/GetAllAccounts", handleGetAllAccounts)
	mux.HandleFunc("/api/AddressTypeToDerivationPath", handleAddressTypeToDerivationPath)
	mux.HandleFunc("/api/GetPathByAddressType", handleGetPathByAddressType)
	mux.HandleFunc("/api/GetBlockWrap", handleGetBlockWrap)
	mux.HandleFunc("/api/GetBlockInfoByHeight", handleGetBlockInfoByHeight)
	mux.HandleFunc("/api/GetWalletBalance", handleGetWalletBalance)
	mux.HandleFunc("/api/ProcessGetWalletBalanceResult", handleProcessGetWalletBalanceResult)
	mux.HandleFunc("/api/CalculateImportedTapAddressBalanceAmount", handleCalculateImportedTapAddressBalanceAmount)
	mux.HandleFunc("/api/GetInfoOfLnd", handleGetInfoOfLnd)
	mux.HandleFunc("/api/GetIdentityPubkey", handleGetIdentityPubkey)
	mux.HandleFunc("/api/GetNewAddress", handleGetNewAddress)
	mux.HandleFunc("/api/AddInvoice", handleAddInvoice)
	mux.HandleFunc("/api/ListInvoices", handleListInvoices)
	mux.HandleFunc("/api/SimplifyInvoice", handleSimplifyInvoice)
	mux.HandleFunc("/api/LookupInvoice", handleLookupInvoice)
	mux.HandleFunc("/api/AbandonChannel", handleAbandonChannel)
	mux.HandleFunc("/api/BatchOpenChannel", handleBatchOpenChannel)
	mux.HandleFunc("/api/ChannelAcceptor", handleChannelAcceptor)
	mux.HandleFunc("/api/ChannelBalance", handleChannelBalance)
	mux.HandleFunc("/api/CheckMacaroonPermissions", handleCheckMacaroonPermissions)
	mux.HandleFunc("/api/CloseChannel", handleCloseChannel)
	mux.HandleFunc("/api/ClosedChannels", handleClosedChannels)
	mux.HandleFunc("/api/DecodePayReq", handleDecodePayReq)
	mux.HandleFunc("/api/ExportAllChannelBackups", handleExportAllChannelBackups)
	mux.HandleFunc("/api/ExportChannelBackup", handleExportChannelBackup)
	mux.HandleFunc("/api/GetChanInfo", handleGetChanInfo)
	mux.HandleFunc("/api/OpenChannelSync", handleOpenChannelSync)
	mux.HandleFunc("/api/OpenChannel", handleOpenChannel)
	mux.HandleFunc("/api/ListChannels", handleListChannels)
	mux.HandleFunc("/api/PendingChannels", handlePendingChannels)
	mux.HandleFunc("/api/GetChannelState", handleGetChannelState)
	mux.HandleFunc("/api/GetChannelInfo", handleGetChannelInfo)
	mux.HandleFunc("/api/RestoreChannelBackups", handleRestoreChannelBackups)
	mux.HandleFunc("/api/SubscribeChannelBackups", handleSubscribeChannelBackups)
	mux.HandleFunc("/api/SubscribeChannelEvents", handleSubscribeChannelEvents)
	mux.HandleFunc("/api/SubscribeChannelGraph", handleSubscribeChannelGraph)
	mux.HandleFunc("/api/UpdateChannelPolicy", handleUpdateChannelPolicy)
	mux.HandleFunc("/api/VerifyChanBackup", handleVerifyChanBackup)
	mux.HandleFunc("/api/ConnectPeer", handleConnectPeer)
	mux.HandleFunc("/api/EstimateFee", handleEstimateFee)
	mux.HandleFunc("/api/SendPaymentSync", handleSendPaymentSync)
	mux.HandleFunc("/api/SendPaymentSync0amt", handleSendPaymentSync0amt)
	mux.HandleFunc("/api/SendCoins", handleSendCoins)
	mux.HandleFunc("/api/SendMany", handleSendMany)
	mux.HandleFunc("/api/SendAllCoins", handleSendAllCoins)
	mux.HandleFunc("/api/LndStopDaemon", handleLndStopDaemon)
	mux.HandleFunc("/api/ListPermissions", handleListPermissions)
	mux.HandleFunc("/api/SendPaymentV2", handleSendPaymentV2)
	mux.HandleFunc("/api/TrackPaymentV2", handleTrackPaymentV2)
	mux.HandleFunc("/api/SendToRouteV2", handleSendToRouteV2)
	mux.HandleFunc("/api/EstimateRouteFee", handleEstimateRouteFee)
	mux.HandleFunc("/api/GetStateForSubscribe", handleGetStateForSubscribe)
	mux.HandleFunc("/api/GetState", handleGetState)
	mux.HandleFunc("/api/GenSeed", handleGenSeed)
	mux.HandleFunc("/api/InitWallet", handleInitWallet)
	mux.HandleFunc("/api/UnlockWallet", handleUnlockWallet)
	mux.HandleFunc("/api/ChangePassword", handleChangePassword)
	mux.HandleFunc("/api/ListAddresses", handleListAddresses)
	mux.HandleFunc("/api/ListAddressesAndGetResponse", handleListAddressesAndGetResponse)
	mux.HandleFunc("/api/ListAccounts", handleListAccounts)
	mux.HandleFunc("/api/FindAccount", handleFindAccount)
	mux.HandleFunc("/api/ListLeases", handleListLeases)
	mux.HandleFunc("/api/ListSweeps", handleListSweeps)
	mux.HandleFunc("/api/ListUnspent", handleListUnspent)
	mux.HandleFunc("/api/NextAddr", handleNextAddr)
	mux.HandleFunc("/api/SetServerHost", handleSetServerHost)
	mux.HandleFunc("/api/GetServerHost", handleGetServerHost)
	mux.HandleFunc("/api/Login", handleLogin)
	mux.HandleFunc("/api/Refresh", handleRefresh)
	mux.HandleFunc("/api/SendPostRequest", handleSendPostRequest)
	mux.HandleFunc("/api/SimplifyTransactions", handleSimplifyTransactions)
	mux.HandleFunc("/api/GetAddressInfoByMempool", handleGetAddressInfoByMempool)
	mux.HandleFunc("/api/GetAddressTransactions", handleGetAddressTransactions)
	mux.HandleFunc("/api/GetAddressTransferOut", handleGetAddressTransferOut)
	mux.HandleFunc("/api/GetAddressTransferOutResult", handleGetAddressTransferOutResult)
	mux.HandleFunc("/api/GetAddressTransactionsByMempool", handleGetAddressTransactionsByMempool)
	mux.HandleFunc("/api/GetAddressTransactionsChainByMempool", handleGetAddressTransactionsChainByMempool)
	mux.HandleFunc("/api/GetAddressTransactionsMempoolByMempool", handleGetAddressTransactionsMempoolByMempool)
	mux.HandleFunc("/api/GetAddressUTXOByMempool", handleGetAddressUTXOByMempool)
	mux.HandleFunc("/api/GetAddressValidationByMempool", handleGetAddressValidationByMempool)
	mux.HandleFunc("/api/GetBlockByMempoolByMempool", handleGetBlockByMempoolByMempool)
	mux.HandleFunc("/api/GetBlockHeaderByMempool", handleGetBlockHeaderByMempool)
	mux.HandleFunc("/api/GetBlockHeightByMempool", handleGetBlockHeightByMempool)
	mux.HandleFunc("/api/GetBlockTimestampByMempool", handleGetBlockTimestampByMempool)
	mux.HandleFunc("/api/GetBlockRawByMempool", handleGetBlockRawByMempool)
	mux.HandleFunc("/api/GetBlockStatusByMempool", handleGetBlockStatusByMempool)
	mux.HandleFunc("/api/GetBlockTipHeightByMempool", handleGetBlockTipHeightByMempool)
	mux.HandleFunc("/api/BlockTipHeight", handleBlockTipHeight)
	mux.HandleFunc("/api/GetBlockTipHashByMempool", handleGetBlockTipHashByMempool)
	mux.HandleFunc("/api/GetBlockTransactionIDByMempool", handleGetBlockTransactionIDByMempool)
	mux.HandleFunc("/api/GetBlockTransactionIDsByMempool", handleGetBlockTransactionIDsByMempool)
	mux.HandleFunc("/api/GetBlockTransactionsByMempool", handleGetBlockTransactionsByMempool)
	mux.HandleFunc("/api/GetBlocksByMempool", handleGetBlocksByMempool)
	mux.HandleFunc("/api/GetBlocksBulkByMempool", handleGetBlocksBulkByMempool)
	mux.HandleFunc("/api/GetMempoolBlocksFeesByMempool", handleGetMempoolBlocksFeesByMempool)
	mux.HandleFunc("/api/GetRecommendedFeesByMempool", handleGetRecommendedFeesByMempool)
	mux.HandleFunc("/api/GetDifficultyAdjustmentByMempool", handleGetDifficultyAdjustmentByMempool)
	mux.HandleFunc("/api/GetNetworkStatsByMempool", handleGetNetworkStatsByMempool)
	mux.HandleFunc("/api/GetNodesSlashChannelsByMempool", handleGetNodesSlashChannelsByMempool)
	mux.HandleFunc("/api/GetNodesInCountryByMempool", handleGetNodesInCountryByMempool)
	mux.HandleFunc("/api/GetNodeStatsPerCountryByMempool", handleGetNodeStatsPerCountryByMempool)
	mux.HandleFunc("/api/GetISPNodesByMempool", handleGetISPNodesByMempool)
	mux.HandleFunc("/api/GetNodeStatsPerISPByMempool", handleGetNodeStatsPerISPByMempool)
	mux.HandleFunc("/api/GetTop100NodesByMempool", handleGetTop100NodesByMempool)
	mux.HandleFunc("/api/GetTop100NodesbyLiquidityByMempool", handleGetTop100NodesbyLiquidityByMempool)
	mux.HandleFunc("/api/GetTop100NodesbyConnectivityByMempool", handleGetTop100NodesbyConnectivityByMempool)
	mux.HandleFunc("/api/GetTop100OldestNodesByMempool", handleGetTop100OldestNodesByMempool)
	mux.HandleFunc("/api/GetNodeStatsByMempool", handleGetNodeStatsByMempool)
	mux.HandleFunc("/api/GetHistoricalNodeStatsByMempool", handleGetHistoricalNodeStatsByMempool)
	mux.HandleFunc("/api/GetChannelByMempool", handleGetChannelByMempool)
	mux.HandleFunc("/api/GetChannelsfromTXIDByMempool", handleGetChannelsfromTXIDByMempool)
	mux.HandleFunc("/api/GetChannelsfromNodePubkeyByMempool", handleGetChannelsfromNodePubkeyByMempool)
	mux.HandleFunc("/api/GetChannelGeodataByMempool", handleGetChannelGeodataByMempool)
	mux.HandleFunc("/api/GetChannelGeodataforNodeByMempool", handleGetChannelGeodataforNodeByMempool)
	mux.HandleFunc("/api/GetMempoolByMempool", handleGetMempoolByMempool)
	mux.HandleFunc("/api/GetMempoolTransactionIDsByMempool", handleGetMempoolTransactionIDsByMempool)
	mux.HandleFunc("/api/GetMempoolRecentByMempool", handleGetMempoolRecentByMempool)
	mux.HandleFunc("/api/GetMempoolRBFTransactionsByMempool", handleGetMempoolRBFTransactionsByMempool)
	mux.HandleFunc("/api/GetMempoolFullRBFTransactionsByMempool", handleGetMempoolFullRBFTransactionsByMempool)
	mux.HandleFunc("/api/GetMiningPoolsByMempool", handleGetMiningPoolsByMempool)
	mux.HandleFunc("/api/GetMiningPoolByMempool", handleGetMiningPoolByMempool)
	mux.HandleFunc("/api/GetMiningPoolHashratesByMempool", handleGetMiningPoolHashratesByMempool)
	mux.HandleFunc("/api/GetMiningPoolHashrateByMempool", handleGetMiningPoolHashrateByMempool)
	mux.HandleFunc("/api/GetMiningPoolBlocksByMempool", handleGetMiningPoolBlocksByMempool)
	mux.HandleFunc("/api/GetHashrateByMempool", handleGetHashrateByMempool)
	mux.HandleFunc("/api/GetDifficultyAdjustmentsByMempool", handleGetDifficultyAdjustmentsByMempool)
	mux.HandleFunc("/api/GetRewardStatsByMempool", handleGetRewardStatsByMempool)
	mux.HandleFunc("/api/GetBlockFeesByMempool", handleGetBlockFeesByMempool)
	mux.HandleFunc("/api/GetBlockRewardsByMempool", handleGetBlockRewardsByMempool)
	mux.HandleFunc("/api/GetBlockFeeratesByMempool", handleGetBlockFeeratesByMempool)
	mux.HandleFunc("/api/GetBlockSizesandWeightsByMempool", handleGetBlockSizesandWeightsByMempool)
	mux.HandleFunc("/api/GetBlockPredictionsByMempool", handleGetBlockPredictionsByMempool)
	mux.HandleFunc("/api/GetBlockAuditScoreByMempool", handleGetBlockAuditScoreByMempool)
	mux.HandleFunc("/api/GetBlocksAuditScoresByMempool", handleGetBlocksAuditScoresByMempool)
	mux.HandleFunc("/api/GetBlockAuditSummaryByMempool", handleGetBlockAuditSummaryByMempool)
	mux.HandleFunc("/api/GetChildrenPayforParentByMempool", handleGetChildrenPayforParentByMempool)
	mux.HandleFunc("/api/GetTransactionByMempool", handleGetTransactionByMempool)
	mux.HandleFunc("/api/GetTransactionHexByMempool", handleGetTransactionHexByMempool)
	mux.HandleFunc("/api/GetTransactionMerkleblockProofByMempool", handleGetTransactionMerkleblockProofByMempool)
	mux.HandleFunc("/api/GetTransactionMerkleProofByMempool", handleGetTransactionMerkleProofByMempool)
	mux.HandleFunc("/api/GetTransactionOutspendByMempool", handleGetTransactionOutspendByMempool)
	mux.HandleFunc("/api/GetTransactionOutspendsByMempool", handleGetTransactionOutspendsByMempool)
	mux.HandleFunc("/api/GetTransactionRawByMempool", handleGetTransactionRawByMempool)
	mux.HandleFunc("/api/GetTransactionRBFHistoryByMempool", handleGetTransactionRBFHistoryByMempool)
	mux.HandleFunc("/api/GetTransactionStatusByMempool", handleGetTransactionStatusByMempool)
	mux.HandleFunc("/api/GetTransactionTimesByMempool", handleGetTransactionTimesByMempool)
	mux.HandleFunc("/api/PostTransactionByMempool", handlePostTransactionByMempool)
	mux.HandleFunc("/api/StartLitd", handleStartLitd)
	mux.HandleFunc("/api/StartLnd", handleStartLnd)
	mux.HandleFunc("/api/StartTapd", handleStartTapd)
	mux.HandleFunc("/api/GetUserOwnIssuanceHistoryInfos", handleGetUserOwnIssuanceHistoryInfos)
	mux.HandleFunc("/api/GetIssuanceTransactionFee", handleGetIssuanceTransactionFee)
	mux.HandleFunc("/api/GetMintTransactionFee", handleGetMintTransactionFee)
	mux.HandleFunc("/api/GetLocalIssuanceTransactionFee", handleGetLocalIssuanceTransactionFee)
	mux.HandleFunc("/api/GetLocalIssuanceTransactionByteSize", handleGetLocalIssuanceTransactionByteSize)
	mux.HandleFunc("/api/GetIssuanceTransactionCalculatedFee", handleGetIssuanceTransactionCalculatedFee)
	mux.HandleFunc("/api/GetMintTransactionCalculatedFee", handleGetMintTransactionCalculatedFee)
	mux.HandleFunc("/api/GetIssuanceTransactionByteSize", handleGetIssuanceTransactionByteSize)
	mux.HandleFunc("/api/GetTapdMintAssetAndFinalizeTransactionByteSize", handleGetTapdMintAssetAndFinalizeTransactionByteSize)
	mux.HandleFunc("/api/GetTapdSendReservedAssetTransactionByteSize", handleGetTapdSendReservedAssetTransactionByteSize)
	mux.HandleFunc("/api/GetMintTransactionByteSize", handleGetMintTransactionByteSize)
	mux.HandleFunc("/api/GetServerOwnSetFairLaunchInfos", handleGetServerOwnSetFairLaunchInfos)
	mux.HandleFunc("/api/ProcessOwnSetFairLaunchResponseToIssuanceHistoryInfo", handleProcessOwnSetFairLaunchResponseToIssuanceHistoryInfo)
	mux.HandleFunc("/api/GetServerFeeRate", handleGetServerFeeRate)
	mux.HandleFunc("/api/GetServerQueryMint", handleGetServerQueryMint)
	mux.HandleFunc("/api/GetServerIssuanceHistoryInfos", handleGetServerIssuanceHistoryInfos)
	mux.HandleFunc("/api/GetLocalTapdIssuanceHistoryInfos", handleGetLocalTapdIssuanceHistoryInfos)
	mux.HandleFunc("/api/GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos", handleGetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos)
	mux.HandleFunc("/api/GetTimestampByBatchTxidWithGetTransactionsResponse", handleGetTimestampByBatchTxidWithGetTransactionsResponse)
	mux.HandleFunc("/api/GetTransactionByBatchTxid", handleGetTransactionByBatchTxid)
	mux.HandleFunc("/api/GetAssetIdByBatchTxidWithListAssetResponse", handleGetAssetIdByBatchTxidWithListAssetResponse)
	mux.HandleFunc("/api/GetAssetIdByOutpointAndNameWithListAssetResponse", handleGetAssetIdByOutpointAndNameWithListAssetResponse)
	mux.HandleFunc("/api/GetAssetsByOutpointWithListAssetResponse", handleGetAssetsByOutpointWithListAssetResponse)
	mux.HandleFunc("/api/GetImageByImageData", handleGetImageByImageData)
	mux.HandleFunc("/api/GetOwnSet", handleGetOwnSet)
	mux.HandleFunc("/api/GetRate", handleGetRate)
	mux.HandleFunc("/api/GetAssetQueryMint", handleGetAssetQueryMint)
	mux.HandleFunc("/api/SendGetReq", handleSendGetReq)
	mux.HandleFunc("/api/AnchorVirtualPsbts", handleAnchorVirtualPsbts)
	mux.HandleFunc("/api/FundVirtualPsbt", handleFundVirtualPsbt)
	mux.HandleFunc("/api/NextInternalKey", handleNextInternalKey)
	mux.HandleFunc("/api/NextScriptKey", handleNextScriptKey)
	mux.HandleFunc("/api/ProveAssetOwnership", handleProveAssetOwnership)
	mux.HandleFunc("/api/RemoveUTXOLease", handleRemoveUTXOLease)
	mux.HandleFunc("/api/SignVirtualPsbt", handleSignVirtualPsbt)
	mux.HandleFunc("/api/VerifyAssetOwnership", handleVerifyAssetOwnership)
	mux.HandleFunc("/api/SimplifyAssetsTransfer", handleSimplifyAssetsTransfer)
	mux.HandleFunc("/api/SimplifyAssetsList", handleSimplifyAssetsList)
	mux.HandleFunc("/api/SyncUniverseFullSpecified", handleSyncUniverseFullSpecified)
	mux.HandleFunc("/api/SyncAssetIssuance", handleSyncAssetIssuance)
	mux.HandleFunc("/api/SyncAssetTransfer", handleSyncAssetTransfer)
	mux.HandleFunc("/api/SyncAssetAll", handleSyncAssetAll)
	mux.HandleFunc("/api/SyncAssetAllSlice", handleSyncAssetAllSlice)
	mux.HandleFunc("/api/SyncAssetAllWithAssets", handleSyncAssetAllWithAssets)
	mux.HandleFunc("/api/GetAllAssetBalances", handleGetAllAssetBalances)
	mux.HandleFunc("/api/GetAllAssetGroupBalances", handleGetAllAssetGroupBalances)
	mux.HandleFunc("/api/GetAllAssetIdByAssetBalance", handleGetAllAssetIdByAssetBalance)
	mux.HandleFunc("/api/SyncAllAssetsByAssetBalance", handleSyncAllAssetsByAssetBalance)
	mux.HandleFunc("/api/GetAllAssetsIdSlice", handleGetAllAssetsIdSlice)
	mux.HandleFunc("/api/AssetKeysTransfer", handleAssetKeysTransfer)
	mux.HandleFunc("/api/AssetLeavesSpecified", handleAssetLeavesSpecified)
	mux.HandleFunc("/api/ProcessAssetTransferLeave", handleProcessAssetTransferLeave)
	mux.HandleFunc("/api/AssetLeavesTransfer", handleAssetLeavesTransfer)
	mux.HandleFunc("/api/AssetLeavesTransfer_ONLY_FOR_TEST", handleAssetLeavesTransfer_ONLY_FOR_TEST)
	mux.HandleFunc("/api/ProcessAssetIssuanceLeave", handleProcessAssetIssuanceLeave)
	mux.HandleFunc("/api/GetAssetInfoByIssuanceLeaf", handleGetAssetInfoByIssuanceLeaf)
	mux.HandleFunc("/api/DecodeRawProofByte", handleDecodeRawProofByte)
	mux.HandleFunc("/api/DecodeRawProofString", handleDecodeRawProofString)
	mux.HandleFunc("/api/ProcessProof", handleProcessProof)
	mux.HandleFunc("/api/DecodeRawProof", handleDecodeRawProof)
	mux.HandleFunc("/api/ProcessListAllAssets", handleProcessListAllAssets)
	mux.HandleFunc("/api/GetAllAssetList", handleGetAllAssetList)
	mux.HandleFunc("/api/ProcessListAllAssetsSimplified", handleProcessListAllAssetsSimplified)
	mux.HandleFunc("/api/GetAllAssetListSimplified", handleGetAllAssetListSimplified)
	mux.HandleFunc("/api/GetAllAssetIdByListAll", handleGetAllAssetIdByListAll)
	mux.HandleFunc("/api/SyncUniverseFullIssuanceByIdSlice", handleSyncUniverseFullIssuanceByIdSlice)
	mux.HandleFunc("/api/SyncUniverseFullTransferByIdSlice", handleSyncUniverseFullTransferByIdSlice)
	mux.HandleFunc("/api/SyncUniverseFullNoSlice", handleSyncUniverseFullNoSlice)
	mux.HandleFunc("/api/OutpointToAddress", handleOutpointToAddress)
	mux.HandleFunc("/api/TransactionAndIndexToAddress", handleTransactionAndIndexToAddress)
	mux.HandleFunc("/api/TransactionAndIndexToValue", handleTransactionAndIndexToValue)
	mux.HandleFunc("/api/CompareScriptKey", handleCompareScriptKey)
	mux.HandleFunc("/api/GetAssetHoldInfosIncludeSpent", handleGetAssetHoldInfosIncludeSpent)
	mux.HandleFunc("/api/GetAssetHoldInfosExcludeSpent", handleGetAssetHoldInfosExcludeSpent)
	mux.HandleFunc("/api/GetAssetHoldInfosIncludeSpentSlow", handleGetAssetHoldInfosIncludeSpentSlow)
	mux.HandleFunc("/api/AddressIsSpent", handleAddressIsSpent)
	mux.HandleFunc("/api/AddressIsSpentAll", handleAddressIsSpentAll)
	mux.HandleFunc("/api/OutpointToTransactionInfo", handleOutpointToTransactionInfo)
	mux.HandleFunc("/api/GetAssetTransactionInfos", handleGetAssetTransactionInfos)
	mux.HandleFunc("/api/SyncAllAssetByList", handleSyncAllAssetByList)
	mux.HandleFunc("/api/GetAssetInfoById", handleGetAssetInfoById)
	mux.HandleFunc("/api/GetAssetHoldInfosExcludeSpentSlow", handleGetAssetHoldInfosExcludeSpentSlow)
	mux.HandleFunc("/api/GetAssetTransactionInfoSlow", handleGetAssetTransactionInfoSlow)
	mux.HandleFunc("/api/AssetIDAndTransferScriptKeyToOutpoint", handleAssetIDAndTransferScriptKeyToOutpoint)
	mux.HandleFunc("/api/GetAllAssetListWithoutProcession", handleGetAllAssetListWithoutProcession)
	mux.HandleFunc("/api/ListBatchesAndGetResponse", handleListBatchesAndGetResponse)
	mux.HandleFunc("/api/GetTransactionsAndGetResponse", handleGetTransactionsAndGetResponse)
	mux.HandleFunc("/api/GetTransactionsExcludeLabelTapdAssetMinting", handleGetTransactionsExcludeLabelTapdAssetMinting)
	mux.HandleFunc("/api/ExcludeLabelIsTapdAssetMinting", handleExcludeLabelIsTapdAssetMinting)
	mux.HandleFunc("/api/ListAssetAndGetResponse", handleListAssetAndGetResponse)
	mux.HandleFunc("/api/ListAssetAndGetResponseByFlags", handleListAssetAndGetResponseByFlags)
	mux.HandleFunc("/api/ListBatchesAndGetCustomResponse", handleListBatchesAndGetCustomResponse)
	mux.HandleFunc("/api/ListAssetAndGetCustomResponse", handleListAssetAndGetCustomResponse)
	mux.HandleFunc("/api/GetTransactionsAndGetCustomResponse", handleGetTransactionsAndGetCustomResponse)
	mux.HandleFunc("/api/AssetLeafKeysIssuance", handleAssetLeafKeysIssuance)
	mux.HandleFunc("/api/AssetLeavesIssuance", handleAssetLeavesIssuance)
	mux.HandleFunc("/api/GetTransactionsWhoseLabelIsTapdAssetMinting", handleGetTransactionsWhoseLabelIsTapdAssetMinting)
	mux.HandleFunc("/api/GetTransactionsWhoseLabelIsNotTapdAssetMinting", handleGetTransactionsWhoseLabelIsNotTapdAssetMinting)
	mux.HandleFunc("/api/DecodeTransactionsWhoseLabelIsNotTapdAssetMinting", handleDecodeTransactionsWhoseLabelIsNotTapdAssetMinting)
	mux.HandleFunc("/api/RawTransactionHexSliceToRequestBodyRawString", handleRawTransactionHexSliceToRequestBodyRawString)
	mux.HandleFunc("/api/PostCallBitcoindToDecodeRawTransaction", handlePostCallBitcoindToDecodeRawTransaction)
	mux.HandleFunc("/api/ProcessDecodedTransactionsData", handleProcessDecodedTransactionsData)
	mux.HandleFunc("/api/GetAndDecodeTransactionsWhoseLabelIsNotTapdAssetMinting", handleGetAndDecodeTransactionsWhoseLabelIsNotTapdAssetMinting)
	mux.HandleFunc("/api/CancelBatch", handleCancelBatch)
	mux.HandleFunc("/api/FinalizeBatch", handleFinalizeBatch)
	mux.HandleFunc("/api/ListBatches", handleListBatches)
	mux.HandleFunc("/api/MintAsset", handleMintAsset)
	mux.HandleFunc("/api/AddGroupAsset", handleAddGroupAsset)
	mux.HandleFunc("/api/NewMeta", handleNewMeta)
	mux.HandleFunc("/api/LoadImageByByte", handleLoadImageByByte)
	mux.HandleFunc("/api/LoadImage", handleLoadImage)
	mux.HandleFunc("/api/ToJsonStr", handleToJsonStr)
	mux.HandleFunc("/api/GetMetaFromStr", handleGetMetaFromStr)
	mux.HandleFunc("/api/SaveImage", handleSaveImage)
	mux.HandleFunc("/api/GetImage", handleGetImage)
	mux.HandleFunc("/api/FetchAssetMeta", handleFetchAssetMeta)
	mux.HandleFunc("/api/ImportProof", handleImportProof)
	mux.HandleFunc("/api/AddrReceives", handleAddrReceives)
	mux.HandleFunc("/api/BurnAsset", handleBurnAsset)
	mux.HandleFunc("/api/DebugLevel", handleDebugLevel)
	mux.HandleFunc("/api/DecodeAddr", handleDecodeAddr)
	mux.HandleFunc("/api/DecodeProof", handleDecodeProof)
	mux.HandleFunc("/api/ExportProof", handleExportProof)
	mux.HandleFunc("/api/GetInfoOfTap", handleGetInfoOfTap)
	mux.HandleFunc("/api/ListAssets", handleListAssets)
	mux.HandleFunc("/api/ListSimpleAssets", handleListSimpleAssets)
	mux.HandleFunc("/api/FindAssetByAssetName", handleFindAssetByAssetName)
	mux.HandleFunc("/api/ListGroups", handleListGroups)
	mux.HandleFunc("/api/QueryAssetTransfers", handleQueryAssetTransfers)
	mux.HandleFunc("/api/ListUtxos", handleListUtxos)
	mux.HandleFunc("/api/NewAddr", handleNewAddr)
	mux.HandleFunc("/api/QueryAddrs", handleQueryAddrs)
	mux.HandleFunc("/api/SendAssets", handleSendAssets)
	mux.HandleFunc("/api/SubscribeReceiveAssetEventNtfns", handleSubscribeReceiveAssetEventNtfns)
	mux.HandleFunc("/api/SubscribeSendAssetEventNtfns", handleSubscribeSendAssetEventNtfns)
	mux.HandleFunc("/api/VerifyProof", handleVerifyProof)
	mux.HandleFunc("/api/TapStopDaemon", handleTapStopDaemon)
	mux.HandleFunc("/api/ProcessListBalancesResponse", handleProcessListBalancesResponse)
	mux.HandleFunc("/api/ProcessListBalancesByGroupKeyResponse", handleProcessListBalancesByGroupKeyResponse)
	mux.HandleFunc("/api/ListBalances", handleListBalances)
	mux.HandleFunc("/api/ListBalancesByGroupKey", handleListBalancesByGroupKey)
	mux.HandleFunc("/api/CheckAssetIssuanceIsLocal", handleCheckAssetIssuanceIsLocal)
	mux.HandleFunc("/api/ListAssetsProcessed", handleListAssetsProcessed)
	mux.HandleFunc("/api/ListAssetsAll", handleListAssetsAll)
	mux.HandleFunc("/api/ListNFTGroups", handleListNFTGroups)
	mux.HandleFunc("/api/ListNFTAssets", handleListNFTAssets)
	mux.HandleFunc("/api/QueryAllNFTByGroup", handleQueryAllNFTByGroup)
	mux.HandleFunc("/api/AddFederationServer", handleAddFederationServer)
	mux.HandleFunc("/api/AssetLeafKeysAndGetResponse", handleAssetLeafKeysAndGetResponse)
	mux.HandleFunc("/api/AssetLeafKeys", handleAssetLeafKeys)
	mux.HandleFunc("/api/AssetLeaves", handleAssetLeaves)
	mux.HandleFunc("/api/GetAssetInfo", handleGetAssetInfo)
	mux.HandleFunc("/api/AssetRoots", handleAssetRoots)
	mux.HandleFunc("/api/DeleteAssetRoot", handleDeleteAssetRoot)
	mux.HandleFunc("/api/DeleteFederationServer", handleDeleteFederationServer)
	mux.HandleFunc("/api/UniverseInfo", handleUniverseInfo)
	mux.HandleFunc("/api/InsertProof", handleInsertProof)
	mux.HandleFunc("/api/ListFederationServers", handleListFederationServers)
	mux.HandleFunc("/api/MultiverseRoot", handleMultiverseRoot)
	mux.HandleFunc("/api/QueryAssetRoots", handleQueryAssetRoots)
	mux.HandleFunc("/api/QueryAssetStats", handleQueryAssetStats)
	mux.HandleFunc("/api/QueryEvents", handleQueryEvents)
	mux.HandleFunc("/api/QueryFederationSyncConfig", handleQueryFederationSyncConfig)
	mux.HandleFunc("/api/QueryProof", handleQueryProof)
	mux.HandleFunc("/api/SetFederationSyncConfig", handleSetFederationSyncConfig)
	mux.HandleFunc("/api/SyncUniverse", handleSyncUniverse)
	mux.HandleFunc("/api/UniverseStats", handleUniverseStats)
	mux.HandleFunc("/api/AssetLeavesAndGetResponse", handleAssetLeavesAndGetResponse)

	fmt.Println("Server is starting...")
	fmt.Println("Available endpoints:")
	fmt.Println("- http://localhost:7047/api/GetApiVersion")
	fmt.Println("- http://localhost:7047/api/NewVersionTag")
	fmt.Println("- http://localhost:7047/api/SetPath")
	fmt.Println("- http://localhost:7047/api/GetPath")
	fmt.Println("- http://localhost:7047/api/FileTestConfig")
	fmt.Println("- http://localhost:7047/api/ReadConfigFile")
	fmt.Println("- http://localhost:7047/api/ReadConfigFile1")
	fmt.Println("- http://localhost:7047/api/ReadConfigFile2")
	fmt.Println("- http://localhost:7047/api/CreateDir")
	fmt.Println("- http://localhost:7047/api/CreateDir2")
	fmt.Println("- http://localhost:7047/api/Visit")
	fmt.Println("- http://localhost:7047/api/CreateFile")
	fmt.Println("- http://localhost:7047/api/ReadFile")
	fmt.Println("- http://localhost:7047/api/CopyFile")
	fmt.Println("- http://localhost:7047/api/DeleteFile")
	fmt.Println("- http://localhost:7047/api/GenerateKeys")
	fmt.Println("- http://localhost:7047/api/GetPublicKey")
	fmt.Println("- http://localhost:7047/api/GetNPublicKey")
	fmt.Println("- http://localhost:7047/api/GetJsonPublicKey")
	fmt.Println("- http://localhost:7047/api/SignMess")
	fmt.Println("- http://localhost:7047/api/RouterForKeyService")
	fmt.Println("- http://localhost:7047/api/LitdStopDaemon")
	fmt.Println("- http://localhost:7047/api/LitdLocalStop")
	fmt.Println("- http://localhost:7047/api/SubServerStatus")
	fmt.Println("- http://localhost:7047/api/GetTapdStatus")
	fmt.Println("- http://localhost:7047/api/GetLitStatus")
	fmt.Println("- http://localhost:7047/api/GetNewAddress_P2TR")
	fmt.Println("- http://localhost:7047/api/GetNewAddress_P2WKH")
	fmt.Println("- http://localhost:7047/api/GetNewAddress_NP2WKH")
	fmt.Println("- http://localhost:7047/api/StoreAddr")
	fmt.Println("- http://localhost:7047/api/RemoveAddr")
	fmt.Println("- http://localhost:7047/api/QueryAddr")
	fmt.Println("- http://localhost:7047/api/QueryAllAddr")
	fmt.Println("- http://localhost:7047/api/GetNonZeroBalanceAddresses")
	fmt.Println("- http://localhost:7047/api/UpdateAllAddressesByGNZBA")
	fmt.Println("- http://localhost:7047/api/GetAllAccountsString")
	fmt.Println("- http://localhost:7047/api/GetAllAccounts")
	fmt.Println("- http://localhost:7047/api/AddressTypeToDerivationPath")
	fmt.Println("- http://localhost:7047/api/GetPathByAddressType")
	fmt.Println("- http://localhost:7047/api/GetBlockWrap")
	fmt.Println("- http://localhost:7047/api/GetBlockInfoByHeight")
	fmt.Println("- http://localhost:7047/api/GetWalletBalance")
	fmt.Println("- http://localhost:7047/api/ProcessGetWalletBalanceResult")
	fmt.Println("- http://localhost:7047/api/CalculateImportedTapAddressBalanceAmount")
	fmt.Println("- http://localhost:7047/api/GetInfoOfLnd")
	fmt.Println("- http://localhost:7047/api/GetIdentityPubkey")
	fmt.Println("- http://localhost:7047/api/GetNewAddress")
	fmt.Println("- http://localhost:7047/api/AddInvoice")
	fmt.Println("- http://localhost:7047/api/ListInvoices")
	fmt.Println("- http://localhost:7047/api/SimplifyInvoice")
	fmt.Println("- http://localhost:7047/api/LookupInvoice")
	fmt.Println("- http://localhost:7047/api/AbandonChannel")
	fmt.Println("- http://localhost:7047/api/BatchOpenChannel")
	fmt.Println("- http://localhost:7047/api/ChannelAcceptor")
	fmt.Println("- http://localhost:7047/api/ChannelBalance")
	fmt.Println("- http://localhost:7047/api/CheckMacaroonPermissions")
	fmt.Println("- http://localhost:7047/api/CloseChannel")
	fmt.Println("- http://localhost:7047/api/ClosedChannels")
	fmt.Println("- http://localhost:7047/api/DecodePayReq")
	fmt.Println("- http://localhost:7047/api/ExportAllChannelBackups")
	fmt.Println("- http://localhost:7047/api/ExportChannelBackup")
	fmt.Println("- http://localhost:7047/api/GetChanInfo")
	fmt.Println("- http://localhost:7047/api/OpenChannelSync")
	fmt.Println("- http://localhost:7047/api/OpenChannel")
	fmt.Println("- http://localhost:7047/api/ListChannels")
	fmt.Println("- http://localhost:7047/api/PendingChannels")
	fmt.Println("- http://localhost:7047/api/GetChannelState")
	fmt.Println("- http://localhost:7047/api/GetChannelInfo")
	fmt.Println("- http://localhost:7047/api/RestoreChannelBackups")
	fmt.Println("- http://localhost:7047/api/SubscribeChannelBackups")
	fmt.Println("- http://localhost:7047/api/SubscribeChannelEvents")
	fmt.Println("- http://localhost:7047/api/SubscribeChannelGraph")
	fmt.Println("- http://localhost:7047/api/UpdateChannelPolicy")
	fmt.Println("- http://localhost:7047/api/VerifyChanBackup")
	fmt.Println("- http://localhost:7047/api/ConnectPeer")
	fmt.Println("- http://localhost:7047/api/EstimateFee")
	fmt.Println("- http://localhost:7047/api/SendPaymentSync")
	fmt.Println("- http://localhost:7047/api/SendPaymentSync0amt")
	fmt.Println("- http://localhost:7047/api/SendCoins")
	fmt.Println("- http://localhost:7047/api/SendMany")
	fmt.Println("- http://localhost:7047/api/SendAllCoins")
	fmt.Println("- http://localhost:7047/api/LndStopDaemon")
	fmt.Println("- http://localhost:7047/api/ListPermissions")
	fmt.Println("- http://localhost:7047/api/SendPaymentV2")
	fmt.Println("- http://localhost:7047/api/TrackPaymentV2")
	fmt.Println("- http://localhost:7047/api/SendToRouteV2")
	fmt.Println("- http://localhost:7047/api/EstimateRouteFee")
	fmt.Println("- http://localhost:7047/api/GetStateForSubscribe")
	fmt.Println("- http://localhost:7047/api/GetState")
	fmt.Println("- http://localhost:7047/api/GenSeed")
	fmt.Println("- http://localhost:7047/api/InitWallet")
	fmt.Println("- http://localhost:7047/api/UnlockWallet")
	fmt.Println("- http://localhost:7047/api/ChangePassword")
	fmt.Println("- http://localhost:7047/api/ListAddresses")
	fmt.Println("- http://localhost:7047/api/ListAddressesAndGetResponse")
	fmt.Println("- http://localhost:7047/api/ListAccounts")
	fmt.Println("- http://localhost:7047/api/FindAccount")
	fmt.Println("- http://localhost:7047/api/ListLeases")
	fmt.Println("- http://localhost:7047/api/ListSweeps")
	fmt.Println("- http://localhost:7047/api/ListUnspent")
	fmt.Println("- http://localhost:7047/api/NextAddr")
	fmt.Println("- http://localhost:7047/api/SetServerHost")
	fmt.Println("- http://localhost:7047/api/GetServerHost")
	fmt.Println("- http://localhost:7047/api/Login")
	fmt.Println("- http://localhost:7047/api/Refresh")
	fmt.Println("- http://localhost:7047/api/SendPostRequest")
	fmt.Println("- http://localhost:7047/api/SimplifyTransactions")
	fmt.Println("- http://localhost:7047/api/GetAddressInfoByMempool")
	fmt.Println("- http://localhost:7047/api/GetAddressTransactions")
	fmt.Println("- http://localhost:7047/api/GetAddressTransferOut")
	fmt.Println("- http://localhost:7047/api/GetAddressTransferOutResult")
	fmt.Println("- http://localhost:7047/api/GetAddressTransactionsByMempool")
	fmt.Println("- http://localhost:7047/api/GetAddressTransactionsChainByMempool")
	fmt.Println("- http://localhost:7047/api/GetAddressTransactionsMempoolByMempool")
	fmt.Println("- http://localhost:7047/api/GetAddressUTXOByMempool")
	fmt.Println("- http://localhost:7047/api/GetAddressValidationByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockByMempoolByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockHeaderByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockHeightByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockTimestampByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockRawByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockStatusByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockTipHeightByMempool")
	fmt.Println("- http://localhost:7047/api/BlockTipHeight")
	fmt.Println("- http://localhost:7047/api/GetBlockTipHashByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockTransactionIDByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockTransactionIDsByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockTransactionsByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlocksByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlocksBulkByMempool")
	fmt.Println("- http://localhost:7047/api/GetMempoolBlocksFeesByMempool")
	fmt.Println("- http://localhost:7047/api/GetRecommendedFeesByMempool")
	fmt.Println("- http://localhost:7047/api/GetDifficultyAdjustmentByMempool")
	fmt.Println("- http://localhost:7047/api/GetNetworkStatsByMempool")
	fmt.Println("- http://localhost:7047/api/GetNodesSlashChannelsByMempool")
	fmt.Println("- http://localhost:7047/api/GetNodesInCountryByMempool")
	fmt.Println("- http://localhost:7047/api/GetNodeStatsPerCountryByMempool")
	fmt.Println("- http://localhost:7047/api/GetISPNodesByMempool")
	fmt.Println("- http://localhost:7047/api/GetNodeStatsPerISPByMempool")
	fmt.Println("- http://localhost:7047/api/GetTop100NodesByMempool")
	fmt.Println("- http://localhost:7047/api/GetTop100NodesbyLiquidityByMempool")
	fmt.Println("- http://localhost:7047/api/GetTop100NodesbyConnectivityByMempool")
	fmt.Println("- http://localhost:7047/api/GetTop100OldestNodesByMempool")
	fmt.Println("- http://localhost:7047/api/GetNodeStatsByMempool")
	fmt.Println("- http://localhost:7047/api/GetHistoricalNodeStatsByMempool")
	fmt.Println("- http://localhost:7047/api/GetChannelByMempool")
	fmt.Println("- http://localhost:7047/api/GetChannelsfromTXIDByMempool")
	fmt.Println("- http://localhost:7047/api/GetChannelsfromNodePubkeyByMempool")
	fmt.Println("- http://localhost:7047/api/GetChannelGeodataByMempool")
	fmt.Println("- http://localhost:7047/api/GetChannelGeodataforNodeByMempool")
	fmt.Println("- http://localhost:7047/api/GetMempoolByMempool")
	fmt.Println("- http://localhost:7047/api/GetMempoolTransactionIDsByMempool")
	fmt.Println("- http://localhost:7047/api/GetMempoolRecentByMempool")
	fmt.Println("- http://localhost:7047/api/GetMempoolRBFTransactionsByMempool")
	fmt.Println("- http://localhost:7047/api/GetMempoolFullRBFTransactionsByMempool")
	fmt.Println("- http://localhost:7047/api/GetMiningPoolsByMempool")
	fmt.Println("- http://localhost:7047/api/GetMiningPoolByMempool")
	fmt.Println("- http://localhost:7047/api/GetMiningPoolHashratesByMempool")
	fmt.Println("- http://localhost:7047/api/GetMiningPoolHashrateByMempool")
	fmt.Println("- http://localhost:7047/api/GetMiningPoolBlocksByMempool")
	fmt.Println("- http://localhost:7047/api/GetHashrateByMempool")
	fmt.Println("- http://localhost:7047/api/GetDifficultyAdjustmentsByMempool")
	fmt.Println("- http://localhost:7047/api/GetRewardStatsByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockFeesByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockRewardsByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockFeeratesByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockSizesandWeightsByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockPredictionsByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockAuditScoreByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlocksAuditScoresByMempool")
	fmt.Println("- http://localhost:7047/api/GetBlockAuditSummaryByMempool")
	fmt.Println("- http://localhost:7047/api/GetChildrenPayforParentByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionHexByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionMerkleblockProofByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionMerkleProofByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionOutspendByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionOutspendsByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionRawByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionRBFHistoryByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionStatusByMempool")
	fmt.Println("- http://localhost:7047/api/GetTransactionTimesByMempool")
	fmt.Println("- http://localhost:7047/api/PostTransactionByMempool")
	fmt.Println("- http://localhost:7047/api/StartLitd")
	fmt.Println("- http://localhost:7047/api/StartLnd")
	fmt.Println("- http://localhost:7047/api/StartTapd")
	fmt.Println("- http://localhost:7047/api/GetUserOwnIssuanceHistoryInfos")
	fmt.Println("- http://localhost:7047/api/GetIssuanceTransactionFee")
	fmt.Println("- http://localhost:7047/api/GetMintTransactionFee")
	fmt.Println("- http://localhost:7047/api/GetLocalIssuanceTransactionFee")
	fmt.Println("- http://localhost:7047/api/GetLocalIssuanceTransactionByteSize")
	fmt.Println("- http://localhost:7047/api/GetIssuanceTransactionCalculatedFee")
	fmt.Println("- http://localhost:7047/api/GetMintTransactionCalculatedFee")
	fmt.Println("- http://localhost:7047/api/GetIssuanceTransactionByteSize")
	fmt.Println("- http://localhost:7047/api/GetTapdMintAssetAndFinalizeTransactionByteSize")
	fmt.Println("- http://localhost:7047/api/GetTapdSendReservedAssetTransactionByteSize")
	fmt.Println("- http://localhost:7047/api/GetMintTransactionByteSize")
	fmt.Println("- http://localhost:7047/api/GetServerOwnSetFairLaunchInfos")
	fmt.Println("- http://localhost:7047/api/ProcessOwnSetFairLaunchResponseToIssuanceHistoryInfo")
	fmt.Println("- http://localhost:7047/api/GetServerFeeRate")
	fmt.Println("- http://localhost:7047/api/GetServerQueryMint")
	fmt.Println("- http://localhost:7047/api/GetServerIssuanceHistoryInfos")
	fmt.Println("- http://localhost:7047/api/GetLocalTapdIssuanceHistoryInfos")
	fmt.Println("- http://localhost:7047/api/GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos")
	fmt.Println("- http://localhost:7047/api/GetTimestampByBatchTxidWithGetTransactionsResponse")
	fmt.Println("- http://localhost:7047/api/GetTransactionByBatchTxid")
	fmt.Println("- http://localhost:7047/api/GetAssetIdByBatchTxidWithListAssetResponse")
	fmt.Println("- http://localhost:7047/api/GetAssetIdByOutpointAndNameWithListAssetResponse")
	fmt.Println("- http://localhost:7047/api/GetAssetsByOutpointWithListAssetResponse")
	fmt.Println("- http://localhost:7047/api/GetImageByImageData")
	fmt.Println("- http://localhost:7047/api/GetOwnSet")
	fmt.Println("- http://localhost:7047/api/GetRate")
	fmt.Println("- http://localhost:7047/api/GetAssetQueryMint")
	fmt.Println("- http://localhost:7047/api/SendGetReq")
	fmt.Println("- http://localhost:7047/api/AnchorVirtualPsbts")
	fmt.Println("- http://localhost:7047/api/FundVirtualPsbt")
	fmt.Println("- http://localhost:7047/api/NextInternalKey")
	fmt.Println("- http://localhost:7047/api/NextScriptKey")
	fmt.Println("- http://localhost:7047/api/ProveAssetOwnership")
	fmt.Println("- http://localhost:7047/api/RemoveUTXOLease")
	fmt.Println("- http://localhost:7047/api/SignVirtualPsbt")
	fmt.Println("- http://localhost:7047/api/VerifyAssetOwnership")
	fmt.Println("- http://localhost:7047/api/SimplifyAssetsTransfer")
	fmt.Println("- http://localhost:7047/api/SimplifyAssetsList")
	fmt.Println("- http://localhost:7047/api/SyncUniverseFullSpecified")
	fmt.Println("- http://localhost:7047/api/SyncAssetIssuance")
	fmt.Println("- http://localhost:7047/api/SyncAssetTransfer")
	fmt.Println("- http://localhost:7047/api/SyncAssetAll")
	fmt.Println("- http://localhost:7047/api/SyncAssetAllSlice")
	fmt.Println("- http://localhost:7047/api/SyncAssetAllWithAssets")
	fmt.Println("- http://localhost:7047/api/GetAllAssetBalances")
	fmt.Println("- http://localhost:7047/api/GetAllAssetGroupBalances")
	fmt.Println("- http://localhost:7047/api/GetAllAssetIdByAssetBalance")
	fmt.Println("- http://localhost:7047/api/SyncAllAssetsByAssetBalance")
	fmt.Println("- http://localhost:7047/api/GetAllAssetsIdSlice")
	fmt.Println("- http://localhost:7047/api/AssetKeysTransfer")
	fmt.Println("- http://localhost:7047/api/AssetLeavesSpecified")
	fmt.Println("- http://localhost:7047/api/ProcessAssetTransferLeave")
	fmt.Println("- http://localhost:7047/api/AssetLeavesTransfer")
	fmt.Println("- http://localhost:7047/api/AssetLeavesTransfer_ONLY_FOR_TEST")
	fmt.Println("- http://localhost:7047/api/ProcessAssetIssuanceLeave")
	fmt.Println("- http://localhost:7047/api/GetAssetInfoByIssuanceLeaf")
	fmt.Println("- http://localhost:7047/api/DecodeRawProofByte")
	fmt.Println("- http://localhost:7047/api/DecodeRawProofString")
	fmt.Println("- http://localhost:7047/api/ProcessProof")
	fmt.Println("- http://localhost:7047/api/DecodeRawProof")
	fmt.Println("- http://localhost:7047/api/ProcessListAllAssets")
	fmt.Println("- http://localhost:7047/api/GetAllAssetList")
	fmt.Println("- http://localhost:7047/api/ProcessListAllAssetsSimplified")
	fmt.Println("- http://localhost:7047/api/GetAllAssetListSimplified")
	fmt.Println("- http://localhost:7047/api/GetAllAssetIdByListAll")
	fmt.Println("- http://localhost:7047/api/SyncUniverseFullIssuanceByIdSlice")
	fmt.Println("- http://localhost:7047/api/SyncUniverseFullTransferByIdSlice")
	fmt.Println("- http://localhost:7047/api/SyncUniverseFullNoSlice")
	fmt.Println("- http://localhost:7047/api/OutpointToAddress")
	fmt.Println("- http://localhost:7047/api/TransactionAndIndexToAddress")
	fmt.Println("- http://localhost:7047/api/TransactionAndIndexToValue")
	fmt.Println("- http://localhost:7047/api/CompareScriptKey")
	fmt.Println("- http://localhost:7047/api/GetAssetHoldInfosIncludeSpent")
	fmt.Println("- http://localhost:7047/api/GetAssetHoldInfosExcludeSpent")
	fmt.Println("- http://localhost:7047/api/GetAssetHoldInfosIncludeSpentSlow")
	fmt.Println("- http://localhost:7047/api/AddressIsSpent")
	fmt.Println("- http://localhost:7047/api/AddressIsSpentAll")
	fmt.Println("- http://localhost:7047/api/OutpointToTransactionInfo")
	fmt.Println("- http://localhost:7047/api/GetAssetTransactionInfos")
	fmt.Println("- http://localhost:7047/api/SyncAllAssetByList")
	fmt.Println("- http://localhost:7047/api/GetAssetInfoById")
	fmt.Println("- http://localhost:7047/api/GetAssetHoldInfosExcludeSpentSlow")
	fmt.Println("- http://localhost:7047/api/GetAssetTransactionInfoSlow")
	fmt.Println("- http://localhost:7047/api/AssetIDAndTransferScriptKeyToOutpoint")
	fmt.Println("- http://localhost:7047/api/GetAllAssetListWithoutProcession")
	fmt.Println("- http://localhost:7047/api/ListBatchesAndGetResponse")
	fmt.Println("- http://localhost:7047/api/GetTransactionsAndGetResponse")
	fmt.Println("- http://localhost:7047/api/GetTransactionsExcludeLabelTapdAssetMinting")
	fmt.Println("- http://localhost:7047/api/ExcludeLabelIsTapdAssetMinting")
	fmt.Println("- http://localhost:7047/api/ListAssetAndGetResponse")
	fmt.Println("- http://localhost:7047/api/ListAssetAndGetResponseByFlags")
	fmt.Println("- http://localhost:7047/api/ListBatchesAndGetCustomResponse")
	fmt.Println("- http://localhost:7047/api/ListAssetAndGetCustomResponse")
	fmt.Println("- http://localhost:7047/api/GetTransactionsAndGetCustomResponse")
	fmt.Println("- http://localhost:7047/api/AssetLeafKeysIssuance")
	fmt.Println("- http://localhost:7047/api/AssetLeavesIssuance")
	fmt.Println("- http://localhost:7047/api/GetTransactionsWhoseLabelIsTapdAssetMinting")
	fmt.Println("- http://localhost:7047/api/GetTransactionsWhoseLabelIsNotTapdAssetMinting")
	fmt.Println("- http://localhost:7047/api/DecodeTransactionsWhoseLabelIsNotTapdAssetMinting")
	fmt.Println("- http://localhost:7047/api/RawTransactionHexSliceToRequestBodyRawString")
	fmt.Println("- http://localhost:7047/api/PostCallBitcoindToDecodeRawTransaction")
	fmt.Println("- http://localhost:7047/api/ProcessDecodedTransactionsData")
	fmt.Println("- http://localhost:7047/api/GetAndDecodeTransactionsWhoseLabelIsNotTapdAssetMinting")
	fmt.Println("- http://localhost:7047/api/CancelBatch")
	fmt.Println("- http://localhost:7047/api/FinalizeBatch")
	fmt.Println("- http://localhost:7047/api/ListBatches")
	fmt.Println("- http://localhost:7047/api/MintAsset")
	fmt.Println("- http://localhost:7047/api/AddGroupAsset")
	fmt.Println("- http://localhost:7047/api/NewMeta")
	fmt.Println("- http://localhost:7047/api/LoadImageByByte")
	fmt.Println("- http://localhost:7047/api/LoadImage")
	fmt.Println("- http://localhost:7047/api/ToJsonStr")
	fmt.Println("- http://localhost:7047/api/GetMetaFromStr")
	fmt.Println("- http://localhost:7047/api/SaveImage")
	fmt.Println("- http://localhost:7047/api/GetImage")
	fmt.Println("- http://localhost:7047/api/FetchAssetMeta")
	fmt.Println("- http://localhost:7047/api/ImportProof")
	fmt.Println("- http://localhost:7047/api/AddrReceives")
	fmt.Println("- http://localhost:7047/api/BurnAsset")
	fmt.Println("- http://localhost:7047/api/DebugLevel")
	fmt.Println("- http://localhost:7047/api/DecodeAddr")
	fmt.Println("- http://localhost:7047/api/DecodeProof")
	fmt.Println("- http://localhost:7047/api/ExportProof")
	fmt.Println("- http://localhost:7047/api/GetInfoOfTap")
	fmt.Println("- http://localhost:7047/api/ListAssets")
	fmt.Println("- http://localhost:7047/api/ListSimpleAssets")
	fmt.Println("- http://localhost:7047/api/FindAssetByAssetName")
	fmt.Println("- http://localhost:7047/api/ListGroups")
	fmt.Println("- http://localhost:7047/api/QueryAssetTransfers")
	fmt.Println("- http://localhost:7047/api/ListUtxos")
	fmt.Println("- http://localhost:7047/api/NewAddr")
	fmt.Println("- http://localhost:7047/api/QueryAddrs")
	fmt.Println("- http://localhost:7047/api/SendAssets")
	fmt.Println("- http://localhost:7047/api/SubscribeReceiveAssetEventNtfns")
	fmt.Println("- http://localhost:7047/api/SubscribeSendAssetEventNtfns")
	fmt.Println("- http://localhost:7047/api/VerifyProof")
	fmt.Println("- http://localhost:7047/api/TapStopDaemon")
	fmt.Println("- http://localhost:7047/api/ProcessListBalancesResponse")
	fmt.Println("- http://localhost:7047/api/ProcessListBalancesByGroupKeyResponse")
	fmt.Println("- http://localhost:7047/api/ListBalances")
	fmt.Println("- http://localhost:7047/api/ListBalancesByGroupKey")
	fmt.Println("- http://localhost:7047/api/CheckAssetIssuanceIsLocal")
	fmt.Println("- http://localhost:7047/api/ListAssetsProcessed")
	fmt.Println("- http://localhost:7047/api/ListAssetsAll")
	fmt.Println("- http://localhost:7047/api/ListNFTGroups")
	fmt.Println("- http://localhost:7047/api/ListNFTAssets")
	fmt.Println("- http://localhost:7047/api/QueryAllNFTByGroup")
	fmt.Println("- http://localhost:7047/api/AddFederationServer")
	fmt.Println("- http://localhost:7047/api/AssetLeafKeysAndGetResponse")
	fmt.Println("- http://localhost:7047/api/AssetLeafKeys")
	fmt.Println("- http://localhost:7047/api/AssetLeaves")
	fmt.Println("- http://localhost:7047/api/GetAssetInfo")
	fmt.Println("- http://localhost:7047/api/AssetRoots")
	fmt.Println("- http://localhost:7047/api/DeleteAssetRoot")
	fmt.Println("- http://localhost:7047/api/DeleteFederationServer")
	fmt.Println("- http://localhost:7047/api/UniverseInfo")
	fmt.Println("- http://localhost:7047/api/InsertProof")
	fmt.Println("- http://localhost:7047/api/ListFederationServers")
	fmt.Println("- http://localhost:7047/api/MultiverseRoot")
	fmt.Println("- http://localhost:7047/api/QueryAssetRoots")
	fmt.Println("- http://localhost:7047/api/QueryAssetStats")
	fmt.Println("- http://localhost:7047/api/QueryEvents")
	fmt.Println("- http://localhost:7047/api/QueryFederationSyncConfig")
	fmt.Println("- http://localhost:7047/api/QueryProof")
	fmt.Println("- http://localhost:7047/api/SetFederationSyncConfig")
	fmt.Println("- http://localhost:7047/api/SyncUniverse")
	fmt.Println("- http://localhost:7047/api/UniverseStats")
	fmt.Println("- http://localhost:7047/api/AssetLeavesAndGetResponse")
	fmt.Println("Listening on :7047")
	http.ListenAndServe(":7047", mux)
}
