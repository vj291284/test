package main

import (
	"fmt"
	//"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var INDEX_HTML []byte

func main() {

	fmt.Println("Starting server on http://localhost:3000/")
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/redirect", CreateRedirectHandler)
	http.ListenAndServe(":3000", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get/")
	w.Write(INDEX_HTML)

}

func CreateRedirectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	r.PostFormValue("RecordingID")
	log.Println(r.PostFormValue("RecordingID"))
	log.Println("Creating Redirect", r.Form)

	fmt.Fprintln(w, "Submiting request to recorder manager")
}

func init() {
	INDEX_HTML, _ = ioutil.ReadFile("./index.html")
}
