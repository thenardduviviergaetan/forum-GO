package main

import (
	"fmt"
	. "forum/pkg/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", ForumHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
