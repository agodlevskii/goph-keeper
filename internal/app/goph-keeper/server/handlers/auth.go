package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/auth"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/user"
)

type UserID string

const uidKey = UserID("uid")

func (h Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("uid")
		if err != nil {
			handleHTTPError(w, err, http.StatusUnauthorized)
			return
		}

		uid, err := h.authService.Authorize(cookie.Value)
		if err != nil {
			handleHTTPError(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), uidKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cid := getClientID(r)

		var req auth.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		token, cid, err := h.authService.Login(r.Context(), cid, req)
		if err != nil {
			if errors.Is(err, auth.ErrWrongCredential) {
				handleHTTPError(w, err, http.StatusUnauthorized)
			} else {
				handleHTTPError(w, err, http.StatusInternalServerError)
			}
			return
		}

		if cid != "" {
			http.SetCookie(w, &http.Cookie{Name: "cid", Value: cid, Path: "/"})
		}

		http.SetCookie(w, &http.Cookie{Name: "uid", Value: token, Path: "/"})
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(token))
	}
}

func (h Handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cid := getClientID(r)
		if loggedOut, err := h.authService.Logout(r.Context(), cid); !loggedOut || err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(""))
	}
}

func (h Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u auth.Request
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		if err := h.authService.Register(r.Context(), u); err != nil {
			if errors.Is(err, user.ErrExists) {
				handleHTTPError(w, err, http.StatusConflict)
			} else {
				handleHTTPError(w, err, http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("User is registered successfully"))
	}
}

func getClientID(r *http.Request) string {
	cid, err := r.Cookie("cid")
	if err != nil {
		return ""
	}
	return cid.Value
}
