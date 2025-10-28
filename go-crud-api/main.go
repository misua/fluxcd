package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Description string  `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store struct {
	mu    sync.RWMutex
	items map[int]*Item
	nextID int
}

func NewStore() *Store {
	return &Store{
		items: make(map[int]*Item),
		nextID: 1,
	}
}

func (s *Store) Create(item *Item) *Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	item.ID = s.nextID
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	s.items[item.ID] = item
	s.nextID++
	
	return item
}

func (s *Store) GetAll() []*Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	items := make([]*Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items
}

func (s *Store) GetByID(id int) (*Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	item, exists := s.items[id]
	return item, exists
}

func (s *Store) Update(id int, updated *Item) (*Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	item, exists := s.items[id]
	if !exists {
		return nil, false
	}
	
	item.Name = updated.Name
	item.Description = updated.Description
	item.UpdatedAt = time.Now()
	
	return item, true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	_, exists := s.items[id]
	if !exists {
		return false
	}
	
	delete(s.items, id)
	return true
}

type API struct {
	store *Store
	env   string
}

func NewAPI(store *Store) *API {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	return &API{
		store: store,
		env:   env,
	}
}

func (api *API) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (api *API) respondError(w http.ResponseWriter, status int, message string) {
	api.respondJSON(w, status, map[string]string{"error": message})
}

// GET /health
func (api *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	api.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "healthy",
		"environment": api.env,
		"timestamp": time.Now(),
	})
}

// GET /items
func (api *API) handleGetItems(w http.ResponseWriter, r *http.Request) {
	items := api.store.GetAll()
	api.respondJSON(w, http.StatusOK, items)
}

// GET /items/{id}
func (api *API) handleGetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	
	item, exists := api.store.GetByID(id)
	if !exists {
		api.respondError(w, http.StatusNotFound, "Item not found")
		return
	}
	
	api.respondJSON(w, http.StatusOK, item)
}

// POST /items
func (api *API) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if item.Name == "" {
		api.respondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	
	created := api.store.Create(&item)
	api.respondJSON(w, http.StatusCreated, created)
}

// PUT /items/{id}
func (api *API) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if item.Name == "" {
		api.respondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	
	updated, exists := api.store.Update(id, &item)
	if !exists {
		api.respondError(w, http.StatusNotFound, "Item not found")
		return
	}
	
	api.respondJSON(w, http.StatusOK, updated)
}

// DELETE /items/{id}
func (api *API) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	
	if !api.store.Delete(id) {
		api.respondError(w, http.StatusNotFound, "Item not found")
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	store := NewStore()
	api := NewAPI(store)
	
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	
	// Routes
	r.HandleFunc("/health", api.handleHealth).Methods("GET")
	r.HandleFunc("/items", api.handleGetItems).Methods("GET")
	r.HandleFunc("/items/{id}", api.handleGetItem).Methods("GET")
	r.HandleFunc("/items", api.handleCreateItem).Methods("POST")
	r.HandleFunc("/items/{id}", api.handleUpdateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", api.handleDeleteItem).Methods("DELETE")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting server on port %s (environment: %s)", port, api.env)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
