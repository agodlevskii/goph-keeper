package handlers

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/auth"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/binary"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/card"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/data"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/password"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/session"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/text"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/user"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	db              any
	authService     auth.Service
	binaryService   binary.Service
	cardService     card.Service
	passwordService password.Service
	textService     text.Service
}

func NewHandler(repoURL string) (*chi.Mux, error) {
	h, err := initHandler(repoURL)
	if err != nil {
		return nil, err
	}

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
				r.Delete("/{id}", h.DeleteBinary())
			})

			r.Route("/card", func(r chi.Router) {
				r.Get("/", h.GetAllCards())
				r.Get("/{id}", h.GetCardByID())
				r.Post("/", h.StoreCard())
				r.Delete("/{id}", h.DeleteCard())
			})

			r.Route("/password", func(r chi.Router) {
				r.Get("/", h.GetAllPasswords())
				r.Get("/{id}", h.GetPasswordByID())
				r.Post("/", h.StorePassword())
				r.Delete("/{id}", h.DeletePassword())
			})

			r.Route("/text", func(r chi.Router) {
				r.Get("/", h.GetAllTexts())
				r.Get("/{id}", h.GetTextByID())
				r.Post("/", h.StoreText())
				r.Delete("/{id}", h.DeleteText())
			})
		})
	})

	return r, nil
}

func initHandler(repoURL string) (Handler, error) {
	dataService, err := data.NewService(repoURL)
	if err != nil {
		return Handler{}, err
	}

	sessionService, err := session.NewService(repoURL)
	if err != nil {
		return Handler{}, err
	}

	userService, err := user.NewService(repoURL)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		authService:     auth.NewService(sessionService, userService),
		binaryService:   binary.NewService(dataService),
		cardService:     card.NewService(dataService),
		passwordService: password.NewService(dataService),
		textService:     text.NewService(dataService),
	}, nil
}

func handleHTTPError(w http.ResponseWriter, err error, code int) {
	log.Error(err)
	http.Error(w, http.StatusText(code), code)
}
