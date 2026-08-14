package operations

import (
	"database/sql"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

func DeleteExperiencePack(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := db.ExecContext(r.Context(), `DELETE FROM experience_packs WHERE id = ?`, r.PathValue("name"))
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
