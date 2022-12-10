package handlers

import (
	"encoding/json"
	"errors"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllCards(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		cs, err := services.GetAllCards(db, uid)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(cs); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func GetCardByID(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		id := chi.URLParam(r, "id")

		c, err := services.GetCardByID(db, uid, id)
		if err != nil && err.Error() != "stored card not found" {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if c.ID == "" {
			handleHTTPError(w, errors.New("data not found"), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(c); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func StoreCard(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)

		var req services.CardReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := services.StoreCard(db, uid, req)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}