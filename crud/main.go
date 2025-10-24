package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Product model
type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

var (
	productStore = make(map[int]Product)
	nextID       = 1
	storeMu      sync.Mutex
)

// Create Product - POST /products
func createProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	p.ID = nextID
	nextID++
	productStore[p.ID] = p
	storeMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// Get Products - GET /products?category=...
func getProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryFilter := strings.ToLower(r.URL.Query().Get("category"))

	storeMu.Lock()
	var result []Product
	for _, p := range productStore {
		if categoryFilter == "" || strings.ToLower(p.Category) == categoryFilter {
			result = append(result, p)
		}
	}
	storeMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Update Product - PUT /products?id=1
func updateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updated Product
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	p, exists := productStore[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Update fields
	p.Name = updated.Name
	p.Category = updated.Category
	p.Price = updated.Price
	productStore[id] = p

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Delete Product - DELETE /products?id=1
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	if _, exists := productStore[id]; !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	delete(productStore, id)

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createProduct(w, r)
		case http.MethodGet:
			getProducts(w, r)
		case http.MethodPut:
			updateProduct(w, r)
		case http.MethodDelete:
			deleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
