package operations

import (
	"database/sql"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/auth"
	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

func DeleteExperiencePack(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.RequireUser(w, r)
		if !ok {
			return
		}
		packID := r.PathValue("name")
		owner, err := packs.IsCreator(r.Context(), db, packID, user.ID)
		if packs.NotFound(err) {
			packs.WriteJSONError(w, "pack not found", http.StatusNotFound)
			return
		}
		if err != nil {
			packs.WriteJSONError(w, "could not load pack", http.StatusInternalServerError)
			return
		}
		if !owner {
			packs.WriteJSONError(w, "only the creator can change this pack", http.StatusForbidden)
			return
		}
		result, err := db.ExecContext(r.Context(), `DELETE FROM experience_packs WHERE id = ?`, packID)
		if err != nil {
			packs.WriteJSONError(w, "could not delete pack", http.StatusInternalServerError)
			return
		}
		count, err := result.RowsAffected()
		if err != nil {
			packs.WriteJSONError(w, "could not delete pack", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			packs.WriteJSONError(w, "pack not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
