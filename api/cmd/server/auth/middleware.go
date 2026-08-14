package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
)

const CookieName = "bedrock_session"

func Middleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err == nil && cookie.Value != "" {
			var user User
			err = db.QueryRowContext(r.Context(), `SELECT u.id, u.username FROM sessions s JOIN users u ON u.id = s.user_id WHERE s.token_hash = ? AND s.expires_at > NOW()`, tokenHash(cookie.Value)).Scan(&user.ID, &user.Username)
			if err == nil {
				user.Name = "users/" + user.ID
				r = r.WithContext(WithUser(r.Context(), user))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func RequireUser(w http.ResponseWriter, r *http.Request) (User, bool) {
	user, ok := CurrentUser(r.Context())
	if !ok {
		WriteJSONError(w, "sign in is required", http.StatusUnauthorized)
	}
	return user, ok
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
