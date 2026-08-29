package product

import (
	"sync"
	"time"
)

type Store struct {
	mu     sync.RWMutex
	items  map[int]Product
	nextID int
}

func NewStore() *Store {
	return &Store{
		items:  make(map[int]Product),
		nextID: 1,
	}
}

func (s *Store) Create(req CreateProductRequest) Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	p := Product{
		ID:        s.nextID,
		Name:      req.Name,
		Category:  req.Category,
		Price:     req.Price,
		Stock:     req.Stock,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.items[s.nextID] = p
	s.nextID++
	return p
}

func (s *Store) GetAll(categoryFilter string, minPriceFilter float64) []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Product, 0, len(s.items))
	for _, item := range s.items {
		if categoryFilter != "" && item.Category != categoryFilter {
			continue
		}
		if minPriceFilter > 0 && item.Price < minPriceFilter {
			continue
		}

		result = append(result, item)
	}
	return result
}

func (s *Store) GetByID(id int) (Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.items[id]
	if !ok {
		return Product{}, ErrNotFound
	}

	return p, nil
}

func (s *Store) Update(id int, req CreateProductRequest) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.items[id]
	if !ok {
		return Product{}, ErrNotFound
	}

	p.Name = req.Name
	p.Category = req.Category
	p.Price = req.Price
	p.Stock = req.Stock
	p.UpdatedAt = time.Now().UTC()

	s.items[id] = p
	return p, nil
}

func (s *Store) Patch(id int, req UpdateProductRequest) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.items[id]
	if !ok {
		return Product{}, ErrNotFound
	}

	if req.Name != nil {
		p.Name = *req.Name
	}

	if req.Category != nil {
		p.Category = *req.Category
	}

	if req.Price != nil {
		p.Price = *req.Price
	}

	if req.Stock != nil {
		p.Stock = *req.Stock
	}

	p.UpdatedAt = time.Now().UTC()
	s.items[id] = p
	return p, nil
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}

	delete(s.items, id)
	return nil
}

func (s *Store) Buy(id int, quantity int) (Product, error) {
	if quantity <= 0 {
		return Product{}, ErrInvalidData
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.items[id]
	if !ok {
		return Product{}, ErrNotFound
	}

	if p.Stock < quantity {
		return Product{}, ErrOutOfStock
	}

	p.Stock -= quantity
	p.UpdatedAt = time.Now().UTC()

	s.items[id] = p
	return p, nil
}
