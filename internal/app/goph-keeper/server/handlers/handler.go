package handlers

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	db   storage.IRepo
	auth services.AuthService
}

func NewHandler(db storage.IRepo) *chi.Mux {
	h := initHandler(db)
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Compress(5, "/*"))

	r.Route("/api/v1/", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", h.Login())
			r.Post("/logout", h.Logout())
			r.Post("/register", h.Register())
		})

		r.With(h.Auth).Route("/storage", func(r chi.Router) {
			r.Route("/binary", func(r chi.Router) {
				r.Get("/", h.GetAllBinaries())
				r.Get("/{id}", h.GetBinaryByID())
				r.Post("/", h.StoreBinary())
			})

			r.Route("/card", func(r chi.Router) {
				r.Get("/", h.GetAllCards())
				r.Get("/{id}", h.GetCardByID())
				r.Post("/", h.StoreCard())
			})

			r.Route("/password", func(r chi.Router) {
				r.Get("/", h.GetAllPasswords())
				r.Get("/{id}", h.GetPasswordByID())
				r.Post("/", h.StorePassword())
			})

			r.Route("/text", func(r chi.Router) {
				r.Get("/", h.GetAllTexts())
				r.Get("/{id}", h.GetTextByID())
				r.Post("/", h.StoreText())
			})
		})
	})

	return r
}

func initHandler(db storage.IRepo) Handler {
	us := services.NewUserService(db)
	ss := services.NewSessionService(db)
	return Handler{
		db:   db,
		auth: services.NewAuthService(ss, us),
	}
}

func handleHTTPError(w http.ResponseWriter, err error, code int) {
	log.Error(err)
	http.Error(w, http.StatusText(code), code)
}
