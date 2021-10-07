package main

import (
	"io"
	"net/http"
	"os"
)

const domain string = "https://jsonplaceholder.typicode.com"

func main() {
	resp, err := http.Get(domain + "/posts")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	os.Stdout.Write(body)
}
