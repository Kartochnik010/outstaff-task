package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func greet(w http.ResponseWriter, r *http.Request) {
	song := r.URL.Query().Get("song")
	if song == "" {
		fmt.Fprintln(w, "query arguments can't be null: song")
		return
	}
	group := r.URL.Query().Get("group")
	if group == "" {
		fmt.Fprintln(w, "query arguments can't be null: group")
		return
	}
	res := map[string]interface{}{
		"releaseDate": "16.07.2006",
		"text":        fmt.Sprintf("%v's song called %v", group, song),
		"link":        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, string(bytes))
}

func main() {
	http.HandleFunc("/info", greet)
	log.Println("Starting music service on port: 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}
}
