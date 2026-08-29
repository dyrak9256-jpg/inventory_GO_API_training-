package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dyrak9256-jpg/inventory_GO_API_training-/internal/handler"
	"github.com/dyrak9256-jpg/inventory_GO_API_training-/internal/product"
)

func main() {
	store := product.NewStore()
	productHandler := handler.NewProductHandler(store)

	r := mux.NewRouter()

	r.Use(handler.RecoveryMiddleware)
	r.Use(handler.LoggerMiddLeware)

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/products", productHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/products", productHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.GetByID).Methods(http.MethodGet)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.Patch).Methods(http.MethodPatch)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.Delete).Methods(http.MethodDelete)
	api.HandleFunc("/products/{id:[0-9]+}/buy", productHandler.Buy).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server is running on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
