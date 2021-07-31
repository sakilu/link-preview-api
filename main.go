package main

import (
	"encoding/json"
	"github.com/badoux/goscraper"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Preview struct {
    Url   	  string   `json:"url"`
    Title         string   `json:"title"`
    Description   string   `json:"description"`
    Images      []string   `json:"images"`
}

func getUrlData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,Authorization,access-control-allow-origin")

	r.ParseForm()
	url, ok := r.Form["url"]
	if !ok {
		http.Error(w, "The url is required", http.StatusBadRequest)
		return
	}
	s, err := goscraper.Scrape(url[0], 5)
	if err != nil {
		http.Error(w, "can't generate preview", http.StatusBadRequest)
		return
	}
	var pvw Preview
	pvw.Url = s.Preview.Link
	pvw.Title = s.Preview.Title
	pvw.Description = s.Preview.Description
	pvw.Images = s.Preview.Images

	json.NewEncoder(w).Encode(pvw)

}

func GetEmptyString(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,Authorization,access-control-allow-origin")
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", getUrlData).Methods("POST")
	router.HandleFunc("/", GetEmptyString).Methods("OPTIONS", "GET")

	router.Use(mux.CORSMethodMiddleware(router))
	log.Fatal(http.ListenAndServe(":4747", router))
}