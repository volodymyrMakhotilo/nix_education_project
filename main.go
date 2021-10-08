package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"nix_education_beginner_project/models"
)

const domain string = "https://jsonplaceholder.typicode.com"
const userID int = 7

var db *sql.DB

func AddComment(comment models.Comment) {
	log.Println("Entered addComments")
	rows, err := db.Query("INSERT  INTO `beginner`.`comments`(post_id,name,email,body)" + fmt.Sprintf("VALUES(%v,'%v','%v','%v') ", comment.PostId, comment.Name, comment.Email, comment.Body))
	defer rows.Close()
	if err != nil {
		panic(err)
	}
}

func getComments(index uint) {
	log.Println("Entered getComments")

	var comments []models.Comment

	resp, err := http.Get(domain + fmt.Sprintf("/comments?postId=%d", index))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &comments)
	if err != nil {
		panic(err)
	}

	for _, comment := range comments {
		go AddComment(comment)
	}
}

func getPost(index int, c chan []byte) {
	log.Println("Entered getPost")

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

func CreateFile(data []byte) {
	file, err := ioutil.TempFile("./storage/posts", "post")
	if err != nil {
		panic(err)
	}
	file.Write(data)
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/beginner")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	c := make(chan []byte, 100)
	for i := 1; i <= 100; i++ {
		go getPost(userID, c)

	}
	for i := 1; i <= 100; i++ {
		CreateFile(<-c)
	}
}
