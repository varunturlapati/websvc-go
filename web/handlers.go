package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	
	"github.com/gorilla/mux"
	
	"github.com/varunturlapati/simpleWebSvc/pkg/article"
	"github.com/varunturlapati/simpleWebSvc/pkg/db"
	"github.com/varunturlapati/simpleWebSvc/pkg/db/redisclient"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Articles")
	rc, _ := redisclient.New()
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
	rc, _ := redisclient.New()
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
	var newArticle article.Article
	err := newArticle.UnmarshalBinary(reqBody)
	if err != nil {
		fmt.Printf("req body was not of type article. Err :%v\n", err)
	}
	rc, _ := redisclient.New()
	err = rc.AddEntry(&db.Entry{
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
	rc, _ := redisclient.New()
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
	var updateArticle article.Article
	err := updateArticle.UnmarshalBinary(reqBody)
	if err != nil {
		fmt.Printf("err updating article with id: %s. Err: %v\n", id, err)
	} else {
		rc, _ := redisclient.New()
		err = rc.ChangeEntry(id, &db.Entry{
			Id:    updateArticle.Id,
			Value: updateArticle,
			Ts:    time.Time{},
		})
		if err != nil {
			fmt.Printf("err updating article with id: %s. Err: %v\n", id, err)
		}
	}
}

func HandleRequests() {
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

