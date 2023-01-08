package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"

	"github.com/go-chi/chi/v5"
)

func (h Handler) DeletePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		if err := h.passwordService.DeletePassword(r.Context(), uid, id); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(""))
	}
}

func (h Handler) GetAllPasswords() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		ps, err := h.passwordService.GetAllPasswords(r.Context(), uid)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		if err = json.NewEncoder(w).Encode(ps); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
		}
	}
}

func (h Handler) GetPasswordByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		p, err := h.passwordService.GetPasswordByID(r.Context(), uid, id)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		if p.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(p); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
		}
	}
}

func (h Handler) StorePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)

		var req models.PasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := h.passwordService.StorePassword(r.Context(), uid, req)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
