package main

import (
	"fmt"
	"net/http"

	"daybook/backend/db"
	"daybook/backend/handlers"
	"daybook/backend/repo"
	"daybook/backend/services"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db.InitDB()
	defer db.Pool.Close()

	pageRepo := repo.NewPostgresPageRepository()
	pageService := services.NewPageService(pageRepo)
	pageHandler := handlers.NewPageHandler(pageService)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("/pages", pageHandler.PagesHandler)
	mux.HandleFunc("/pages/", pageHandler.PageByIDHandler)

	fmt.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", enableCORS(mux))
	if err != nil {
		fmt.Println("Server Failure: ", err)
	}
}
