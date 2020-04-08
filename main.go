package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/varunturlapati/simpleWebSvc/Db"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	a "github.com/varunturlapati/simpleWebSvc/Article"
	"github.com/varunturlapati/simpleWebSvc/Db/RedisClient"
)

var articles = a.Articles{
	a.Article{
		Id:      "1",
		Title:   "Test Title",
		Desc:    "Test Desc",
		Content: "Test Content",
	},
	a.Article{
		Id:      "2",
		Title:   "Test Title2",
		Desc:    "Test Desc2",
		Content: "Test Content2",
	},
	a.Article{
		Id:      "3",
		Title:   "Test Title3",
		Desc:    "Test Desc3",
		Content: "Test Content3",
	},
	a.Article{
		Id:      "3",
		Title:   "Test Title3",
		Desc:    "Test Desc3",
		Content: "Test Content3",
	},
	a.Article{
		Id:      "4",
		Title:   "Test Title4",
		Desc:    "Test Desc4",
		Content: "Test Content4",
	},
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Articles")
	rc, _ := RedisClient.New()
	res, err := rc.GetAllEntries()
	if err != nil {
		fmt.Printf("error fetching all entries. Err: %v\n", err)
	} else {
		fmt.Printf("Success")
		json.NewEncoder(w).Encode(res)
	}
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("Endpoint Hit: returnSingleArticle with id = " + key)
	rc, _ := RedisClient.New()
	res, err := rc.GetEntry(key)
	if err != nil {
		fmt.Printf("No such article with id = %s. Err: %v\n", key, err)
		json.NewEncoder(w).Encode(nil)
	} else {
		json.NewEncoder(w).Encode(res.Value)
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newArticle a.Article
	err := newArticle.UnmarshalBinary(reqBody)
	if err != nil {
		fmt.Printf("req body was not of type article. Err :%v\n", err)
	}
	rc, _ := RedisClient.New()
	err = rc.AddEntry(&Db.Entry{
		Id:    newArticle.Id,
		Value: newArticle,
		Ts:    time.Time{},
	})

	if err != nil {
		fmt.Printf("Err inserting an article. Err: %v\n", err)
	}

	json.NewEncoder(w).Encode(newArticle)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteArticle")
	vars := mux.Vars(r)
	id := vars["id"]
	/*
		for ind, art := range articles {
			if art.Id == id {
				articles = append(articles[:ind], articles[ind+1:]...)
			}
		}
	*/
	rc, _ := RedisClient.New()
	err := rc.RemoveEntry(id)
	if err != nil {
		fmt.Printf("err removing article with id: %s. Err: %v\n", id, err)
	} else {
		fmt.Printf("removed article with id: %s\n", id)
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateArticle")
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var updateArticle a.Article
	err := updateArticle.UnmarshalBinary(reqBody)
	if err != nil {
		fmt.Printf("err updating article with id: %s. Err: %v\n", id, err)
	} else {
		rc, _ := RedisClient.New()
		err = rc.ChangeEntry(id, &Db.Entry{
			Id:    updateArticle.Id,
			Value: updateArticle,
			Ts:    time.Time{},
		})
		if err != nil {
			fmt.Printf("err updating article with id: %s. Err: %v\n", id, err)
		}
	}
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/articles", getAllArticles)
	myRouter.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle)

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	rClient, _ := RedisClient.New()
	pong, err := rClient.PingPong()
	fmt.Printf("%v and %v\n", pong, err)
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()

}
