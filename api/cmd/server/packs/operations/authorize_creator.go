package operations

import (
	"database/sql"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/auth"
	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

func authorizeCreator(w http.ResponseWriter, r *http.Request, db *sql.DB, packID string) bool {
	user, ok := auth.RequireUser(w, r)
	if !ok {
		return false
	}
	creator, err := packs.IsCreator(r.Context(), db, packID, user.ID)
	if packs.NotFound(err) {
		packs.WriteJSONError(w, "pack not found", http.StatusNotFound)
		return false
	}
	if err != nil {
		packs.WriteJSONError(w, "could not load pack", http.StatusInternalServerError)
		return false
	}
	if !creator {
		packs.WriteJSONError(w, "only the creator can change this pack", http.StatusForbidden)
		return false
	}
	return true
}
