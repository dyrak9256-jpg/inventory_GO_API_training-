package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/dyrak9256-jpg/inventory_GO_API_training-/internal/product"
)

type ProductHandler struct {
	store *product.Store
}

func NewProductHandler(store *product.Store) *ProductHandler {
	return &ProductHandler{store: store}
}

// Create handles POST /products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req product.CreatedProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateCreateRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := h.store.Create(req)
	respondWithJSON(w, http.StatusCreated, p)
}

// GetAll handles GET /products
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	minPriceStr := r.URL.Query().Get("min_price")

	var minPrice float64
	if minPriceStr != "" {
		val, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid min_price query parameter")
			return
		}
		minPrice = val
	}

	products := h.store.GetAll(category, minPrice)
	respondWithJSON(w, http.StatusOK, products)
}

// GetByID handles GET /products/{id}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	p, err := h.store.GetByID(id)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// Update handles PUT /products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var req product.CreatedProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateCreateRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	p, err := h.store.Update(id, req)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// Patch handles PATCH /products/{id}
func (h *ProductHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var req product.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.store.Patch(id, req)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// Delete handles DELETE /products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.store.Delete(id); err != nil {
		mapDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Buy handles POST /products/{id}/buy
func (h *ProductHandler) Buy(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var req product.BuyProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.store.Buy(id, req.Quantity)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// --- Auxiliary functions ---

func validateCreateRequest(req product.CreatedProductRequest) error {
	if req.Name == "" {
		return errors.New("name cannot be empty")
	}
	if req.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}

func parseIDFromVars(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return 0, strconv.ErrSyntax
	}
	return strconv.Atoi(idStr)
}
