package handlers

import (
	"encoding/json"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllBinaries(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		bs, err := services.GetAllBinaries(db, uid)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(bs); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func GetBinaryByID(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		id := chi.URLParam(r, "id")

		b, err := services.GetBinaryByID(db, uid, id)
		if err != nil && err.Error() != "stored binary not found" {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if b.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(b); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func StoreBinary(db storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)

		var req services.BinaryReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := services.StoreBinary(db, uid, req)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
