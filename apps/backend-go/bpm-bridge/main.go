package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("BPM Bridge Service starting on :8080...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
// Force rebuild 1
