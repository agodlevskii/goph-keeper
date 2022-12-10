package handlers

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

func NewHandler(db storage.IStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Compress(5, "/*"))

	r.Route("/api/v1/", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", Login(db))
			r.Post("/register", Register(db))
		})

		r.With(Auth).Route("/storage", func(r chi.Router) {
			r.Route("/binary", func(r chi.Router) {
				r.Get("/", GetAllBinaries(db))
				r.Get("/{id}", GetBinaryByID(db))
				r.Post("/", StoreBinary(db))
			})

			r.Route("/card", func(r chi.Router) {
				r.Get("/", GetAllCards(db))
				r.Get("/{id}", GetCardByID(db))
				r.Post("/", StoreCard(db))
			})

			r.Route("/password", func(r chi.Router) {
				r.Get("/", GetAllPasswords(db))
				r.Get("/{id}", GetPasswordByID(db))
				r.Post("/", StorePassword(db))
			})

			r.Route("/text", func(r chi.Router) {
				r.Get("/", GetAllTexts(db))
				r.Get("/{id}", GetTextByID(db))
				r.Post("/", StoreText(db))
			})
		})
	})

	return r
}

func handleHTTPError(w http.ResponseWriter, err error, code int) {
	log.Error(err)
	http.Error(w, http.StatusText(code), code)
}
