package main

import (
	"fmt"
	"net/http"

	"github.com/megapk/fake_http_server/database"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/gorilla/handlers"
)

var db *database.DatabaseConnection

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	link := r.URL.String()
	fmt.Fprintf(w, "My url: %s", link)

	header, _ := json.Marshal(r.Header)
	body, _ := ioutil.ReadAll(r.Body)
	form, _ := json.Marshal(r.Form)
	postForm, _ := json.Marshal(r.PostForm)
	db.AddUrl(&database.Bot{
		Link: link,
		Header: string(header),
		Body: string(body),
		Form: string(form),
		PostForm: string(postForm),
	})
}

func GetUrls(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(db.GetUrls())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	db = database.New("db.sqlite3")
	http.HandleFunc("/get_urls", GetUrls)
	http.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler)))
	http.ListenAndServe(":8080", nil)
}