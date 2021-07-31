package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"io"
	"os"
	"net/http"
	"github.com/badoux/goscraper"
	"github.com/gorilla/handlers"
)

type Preview struct {
    Url   	  string   `json:"url"`
    Title         string   `json:"title"`
    Description   string   `json:"description"`
    Images      []string   `json:"images"`
}

func getUrlData(w http.ResponseWriter, r *http.Request) {
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

func getUrl(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "")
}
func GetEmptyString(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "")
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", getUrlData).Methods("POST")
	router.HandleFunc("/", getUrl).Methods("GET")
	router.HandleFunc("/", GetEmptyString).Methods("OPTIONS")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{GetOrigins()})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	http.ListenAndServe(GetPort(), handlers.CORS(originsOk, headersOk, methodsOk)(router))
}

func GetOrigins() string {
	var origin = os.Getenv("ORIGIN_ALLOWED")
	if origin == "" {
		origin = "*"
		fmt.Println("INFO: No ORIGIN_ALLOWED environment variable detected, defaulting to " + origin)
	}
	return origin
}
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
