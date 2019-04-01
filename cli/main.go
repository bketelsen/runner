package main

import (
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type Code struct {
	Code string `json:"code`
}

func main() {
	var filename string
	var server string
	flag.StringVar(&filename, "file", "sample.go", "the go file to execute on the server")
	flag.StringVar(&server, "server", "http://localhost:8080", "the server to execute the code on")
	flag.Parse()
	codeStr, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	code := Code{Code: base64.StdEncoding.EncodeToString(codeStr)}
	req := gorequest.New()
	req = req.Post(server).
		Send(code).
		Set("Content-Type", "application/json")
	curlCmd, err := req.AsCurlCommand()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(curlCmd)
	resp, body, errs := req.End()
	if len(errs) > 0 {
		log.Fatal(errs)
	} else if resp.StatusCode != http.StatusOK {
		log.Println("error response!")
		log.Println(resp)
	}
	log.Println(body)
}
