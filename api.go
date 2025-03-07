package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type application struct {
	addr string
	db   *sql.DB
}

func NewApplication(addr string, db *sql.DB) *application {
	return &application{
		addr: addr,
		db:   db,
	}
}

func (a *application) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.HandleFunc("/shorten", a.shortenURL).Methods("POST")
	subrouter.HandleFunc("/{shortcode}", a.resolveURL).Methods("GET")

	fmt.Println("Listening on", a.addr)

	return http.ListenAndServe(a.addr, router)
}

type urlBody struct {
	Value string `json:"value"`
}

func (a *application) shortenURL(w http.ResponseWriter, r *http.Request) {
	var url urlBody
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Failed to decode url", http.StatusBadRequest)
		return
	}

	if url.Value == "" {
		http.Error(w, "Url missing", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode()

	query := `INSERT INTO urls (long_url, short_url) VALUES($1,$2) RETURNING id`

	var id int
	err := a.db.QueryRow(query, url.Value, shortCode).Scan(&id)
	if err != nil {
		http.Error(w, "Error generating url", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message":   "Generated URL successfully",
		"short_url": shortCode,
		"id":        id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}

type return_url struct {
	Value string `json:"value"`
}

func (a *application) resolveURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortcode"]

	query := `SELECT long_url FROM urls WHERE short_url = $1`

	var res return_url
	err := a.db.QueryRow(query, shortCode).Scan(&res.Value)
	if err != nil {
		http.Error(w, "Error searching url", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}

func generateShortCode() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(b)[:6]
}
