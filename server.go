package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Object struct {
	name       string
	longtitude float64
	latitude   float64
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(r.URL.String())
	if err != nil {
		log.Printf("Failed to parse url: %s with error %s\n", url.String(), err.Error())
		w.Write([]byte("Url malformed, please try again!"))
		return
	}
	log.Printf("Request from %s URL: %s\n", r.RemoteAddr, url)
	params, err := parseParams(url.Query())
	if err != nil {
		errMsg := "Error when extracting search params:" + err.Error()
		log.Println(errMsg)
		w.Write([]byte(errMsg))
	}
	topItems := topItems(20, params)
	serializedItems := strings.Join(topItems, "\n")
	w.Write([]byte(serializedItems))
}

func parseParams(rawParams url.Values) (paramsObj Object, err error) {
	if len(rawParams["searchTerm"]) > 1 || len(rawParams["lng"]) > 1 || len(rawParams["lat"]) > 1 {
		err = errors.New("Max 1 search term")
		return
	}
	paramsObj.latitude, err = strconv.ParseFloat(rawParams["lat"][0], 64)
	if err != nil {
		return
	}
	paramsObj.longtitude, err = strconv.ParseFloat(rawParams["lng"][0], 64)
	if err != nil {
		return
	}
	paramsObj.name = strings.ToLower(rawParams["searchTerm"][0])
	return
}

func main() {
	_, err := classifyDatabase()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
