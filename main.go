package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const domain string = "https://jsonplaceholder.typicode.com"

func getPost(index int, c chan []byte) {
	resp, err := http.Get(domain + fmt.Sprintf("/posts/%d", index))
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	c <- body
}
func CreateFile(data []byte) {
	file, err := ioutil.TempFile("./storage/posts", "post")
	if err != nil {
		panic(err)
	}
	file.Write(data)
}

func main() {
	c := make(chan []byte, 100)
	for i := 1; i <= 100; i++ {
		go getPost(i, c)
	}
	for i := 1; i <= 100; i++ {
		CreateFile(<-c)
	}
}
