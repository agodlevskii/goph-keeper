package handlers

import (
	"encoding/json"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllTexts(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		ts, err := services.GetAllTexts(db, uid)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(ts); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func GetTextByID(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		id := chi.URLParam(r, "id")

		t, err := services.GetTextByID(db, uid, id)
		if err != nil && err.Error() != "stored text not found" {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if t.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(t); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func StoreText(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)

		var req services.TextReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := services.StoreText(db, uid, req)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
