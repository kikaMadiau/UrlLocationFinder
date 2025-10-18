package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/find", finderHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
