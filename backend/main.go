package main

import (
	"fmt"
	"net/http"

	"daybook/backend/handlers"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/pages", handlers.PagesHandler)
	http.HandleFunc("/pages/", handlers.PageByIDHandler)

	fmt.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server Failure: ", err)
	}
}
