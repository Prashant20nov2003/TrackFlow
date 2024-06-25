package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/track", track)
	http.ListenAndServe(":9876", nil)
}

func track(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	fmt.Println("received", data)
}
