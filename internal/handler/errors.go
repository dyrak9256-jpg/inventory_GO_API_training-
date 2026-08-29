package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dyrak9256-jpg/inventory_GO_API_training-/internal/product"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

func mapDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, product.ErrNotFound):
		respondWithError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, product.ErrOutOfStock):
		respondWithError(w, http.StatusConflict, err.Error())
	case errors.Is(err, product.ErrInvalidData):
		respondWithError(w, http.StatusBadRequest, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "internal server error")
	}
}
