package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type GuestBook struct {
	PostsCount int
	Posts      []string
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/store", storeHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}

func storeHandler(writer http.ResponseWriter, request *http.Request) {
	post := request.FormValue("post")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("data/posts.txt", options, os.FileMode(0600))
	check(err)
	_, err = fmt.Fprintln(file, post)
	check(err)
	err = file.Close()
	check(err)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("templates/create.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	posts := getStrings("data/posts.txt")
	html, err := template.ParseFiles("templates/posts.html")
	check(err)
	guestBook := GuestBook{
		PostsCount: len(posts),
		Posts:      posts,
	}
	err = html.Execute(writer, guestBook)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}
