package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type RequestBody struct {
	URL string `json:"url"`
}

func finderHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.URL == "" {
		http.Error(w, "Missing 'url' field in JSON", http.StatusBadRequest)
		return
	}

	jsonData, err := UrlLocationFinder(reqBody.URL)
	if err != nil {
		http.Error(w, "Error locating URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
