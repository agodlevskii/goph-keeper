package handlers

import (
	"context"
	"encoding/json"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("uid")
		if err != nil {
			handleHTTPError(w, err, http.StatusUnauthorized)
			return
		}

		uid, err := services.Authorize(cookie.Value)
		if err != nil {
			handleHTTPError(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "uid", uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cid := getClientID(r)

		var req services.AuthReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		token, cid, err := services.Login(h.db, cid, req)
		if err != nil {
			if err.Error() == "invalid username or password" {
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

func (h Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u services.AuthReq
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		if err := services.Register(h.db, u); err != nil {
			if err.Error() == "user with the specified name already exists" {
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
