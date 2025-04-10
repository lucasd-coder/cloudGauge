package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type controller struct{}

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	return r
}
