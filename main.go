package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/", runHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

type codeBody struct {
	Encoded string `json:"code"`
}

func (c *codeBody) Decode() (string, error) {
	b, err := base64.StdEncoding.DecodeString(c.Encoded)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func runHandler(c echo.Context) error {
	encoded := new(codeBody)
	if err := c.Bind(encoded); err != nil {
		return err
	}
	code, err := encoded.Decode()
	if err != nil {
		return err
	}
	log.Println("code below")
	log.Println(code)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	tmpDir, err := ioutil.TempDir(wd, "wasmrun-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir) // TODO: check errr

	// write supporting files
	nodeScriptLoc := filepath.Join(tmpDir, "node_exec_wasm.sh")
	if err := ioutil.WriteFile(
		nodeScriptLoc,
		[]byte(nodeExecWasmSH),
		0777,
	); err != nil {
		return err
	}
	log.Println("wrote", nodeScriptLoc)

	wasmExecScriptLoc := filepath.Join(tmpDir, "wasm_exec.js")
	if err := ioutil.WriteFile(
		wasmExecScriptLoc,
		[]byte(wasmExecJS),
		0777,
	); err != nil {
		return err
	}
	log.Println("wrote", wasmExecScriptLoc)

	mainScriptLoc := filepath.Join(tmpDir, "main.go")
	if err := ioutil.WriteFile(
		mainScriptLoc,
		[]byte(code),
		0777,
	); err != nil {
		return err
	}
	log.Printf("wrote %s (%d bytes)", mainScriptLoc, len(code))

	cmd := exec.Command(
		"go",
		"run",
		fmt.Sprintf(`-exec="%s"`, nodeScriptLoc),
		".",
	)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Dir = filepath.Join(tmpDir)
	log.Println(strings.Join(cmd.Args, " "))
	log.Println("about to run the command inside", cmd.Dir)
	log.Println("env vars", cmd.Env)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("error", err)
		return err
	}
	outStr := string(out)
	log.Println(outStr)
	return c.String(http.StatusOK, outStr)
}

// func buildHandler(w http.ResponseWriter, r *http.Request) {
// 	cmd := exec.Command("go", "build", "-o", "lib.wasm", "./wasm/main.go")
// 	cmd.Env = append(cmd.Env, "GOARCH=wasm")
// 	cmd.Env = append(cmd.Env, "GOOS=js")
// 	if err := cmd.Run(); err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("failure!"))
// 	}
// 	w.WriteHeader(http.StatusOK)
// }
