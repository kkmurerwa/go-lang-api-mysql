package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"description"`
	Content string `json:"content"`
}

type User struct {
	ID string `json:"_id"`
	UserName string `json:"user_name"`
	FirstName string `json:"first_name"`
	Password string `json:"password"`
	CreatedOn string `json:"created_on"`
}

var db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_db")

type BaseResponse struct {
	StatusCode bool `json:"statusCode"`
	Message string `json:"Message"`
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
	var articles []Article

	results, err := db.Query("SELECT id, title, description, content FROM articles")

	if err != nil {
		//fmt.Printf("Could not fetch values: %v\n", err.Error())

		response := BaseResponse{StatusCode: true, Message: "Could not create article"}

		json.NewEncoder(w).Encode(response)
	} else {
		for results.Next() {
			var article Article

			err := results.Scan(&article.Id, &article.Title, &article.Desc, &article.Content)

			if err != nil {
				panic("Could not fetch articles")
			}

			articles = append(articles, article)

			fmt.Println(article)
		}


		fmt.Println("Endpoint Hit: returnAllArticles")
		json.NewEncoder(w).Encode(articles)
	}



	results.Close()
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	results, err := db.Query("SELECT id, title, description, content FROM articles WHERE id = ?", key)

	if err != nil {
		fmt.Printf("Could not fetch values: %v\n", err.Error())

		response := BaseResponse{StatusCode: true, Message: "Could not find article"}

		json.NewEncoder(w).Encode(response)
	}

	if err != nil {
		//fmt.Printf("Could not fetch values: %v\n", err.Error())

		response := BaseResponse{StatusCode: true, Message: "Could not create article"}

		json.NewEncoder(w).Encode(response)
	} else {
		if results.Next() {
			var article Article

			err := results.Scan(&article.Id, &article.Title, &article.Desc, &article.Content)

			if err != nil {
				panic("Could not fetch articles")
			}

			fmt.Println(article)

			fmt.Println("Endpoint Hit: returnAllArticles")
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	insert, err := db.Query("INSERT INTO articles(title, description, content, created_on) VALUES (?, ?, ?, NOW())", article.Title, article.Desc, article.Content)

	if err != nil {
		fmt.Printf("Could not insert article: %v\n", err.Error())
	}

	//defer insert.Close()

	fmt.Println("Successfully inserted into articles table")
	response := BaseResponse{StatusCode: true, Message: "Successfully created article"}

	json.NewEncoder(w).Encode(response)

	err = insert.Close()
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	delete, err := db.Query("DELETE FROM articles WHERE id = ?", id)

	if err != nil {
		response := BaseResponse{StatusCode: true, Message: "Could not delete article"}

		json.NewEncoder(w).Encode(response)
	} else {
		response := BaseResponse{StatusCode: true, Message: "Successfully deleted article"}

		json.NewEncoder(w).Encode(response)
	}

	delete.Close()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	// add our new DELETE endpoint here
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	handleRequests()
}