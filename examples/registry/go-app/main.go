package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	fmt.Println("Listening on port 8000...")

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8093", nil)
	if err != nil {
		panic(err)
	}
}
