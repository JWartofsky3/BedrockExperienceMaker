package operations

import (
	"database/sql"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/auth"
)

func GetCurrentUser(_ *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.RequireUser(w, r)
		if !ok {
			return
		}
		auth.WriteJSON(w, user, http.StatusOK)
	}
}
