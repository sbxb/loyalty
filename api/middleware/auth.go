package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/services/auth"
)

func AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload, err := auth.GetPayload(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := models.UserAuth{}
		err = json.Unmarshal([]byte(payload), &user)
		if err != nil || user.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.ContextUserKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
