package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/lint"
	"io"
	"io/ioutil"
	"net/http"
	"nix_education_beginner_project/models"
	"time"
)

const domain string = "https://jsonplaceholder.typicode.com"
const userID int = 7

var db *sql.DB

func addComment(comment models.Comment) {
	db.Exec("INSERT  INTO `beginner`.`comments`(post_id,name,email,body)" + fmt.Sprintf("VALUES(%v,'%v','%v','%v') ", comment.PostId, comment.Name, comment.Email, comment.Body))
}

func getComments(index uint) {
	var comments []models.Comment

	resp, err := http.Get(domain + fmt.Sprintf("/comments?postId=%d", index))
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &comments)
	if err != nil {
		panic(err)
	}

	for _, comment := range comments {
		go addComment(comment)
	}
}

func getPost(index int, c chan []byte) {
	var posts []models.Post

	resp, err := http.Get(domain + fmt.Sprintf("/posts?userId=%d", index))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &posts)
	if err != nil {
		panic(err)
	}

	for _, post := range posts {
		go getComments(post.Id)
	}

	c <- body
}

func createFile(data []byte) {
	file, err := ioutil.TempFile("./storage/posts", "post")
	if err != nil {
		panic(err)
	}
	file.Write(data)
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/beginner")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(50 * time.Millisecond)
	db.SetMaxIdleConns(100)

	c := make(chan []byte, 100)
	for i := 1; i <= 100; i++ {
		go getPost(userID, c)

	}
	for i := 1; i <= 100; i++ {
		createFile(<-c)
	}

}
