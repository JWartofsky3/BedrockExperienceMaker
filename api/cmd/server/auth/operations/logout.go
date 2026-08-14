package operations

import (
	"database/sql"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/auth"
)

func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(auth.CookieName); err == nil {
			_, _ = db.ExecContext(r.Context(), `DELETE FROM sessions WHERE token_hash = SHA2(?, 256)`, cookie.Value)
		}
		http.SetCookie(w, &http.Cookie{Name: auth.CookieName, Value: "", Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: -1})
		w.WriteHeader(http.StatusNoContent)
	}
}
