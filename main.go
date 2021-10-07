package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Post struct {
	UserId uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

const domain string = "https://jsonplaceholder.typicode.com"

func main() {
	resp, err := http.Get(domain + "/posts")
	if err != nil {
		panic(err)
	}
	/*var posts []Post*/
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(os.Stdout.Write(body))
}
