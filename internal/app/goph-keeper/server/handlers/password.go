package handlers

import (
	"encoding/json"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllPasswords(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		ps, err := services.GetAllPasswords(db, uid)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(ps); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func GetPasswordByID(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		id := chi.URLParam(r, "id")

		p, err := services.GetPasswordByID(db, uid, id)
		if err != nil && err.Error() != "stored password not found" {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if p.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(p); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func StorePassword(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)

		var req services.PasswordReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := services.StorePassword(db, uid, req)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
