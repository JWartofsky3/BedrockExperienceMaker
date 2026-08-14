package operations

import (
	"database/sql"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

func RemoveExperiencePackAddon(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		packID := r.PathValue("name")
		addonID := r.PathValue("addon")
		if !authorizeCreator(w, r, db, packID) {
			return
		}
		if _, err := packs.Get(r.Context(), db, packID, false); packs.NotFound(err) {
			packs.WriteJSONError(w, "pack not found", http.StatusNotFound)
			return
		} else if err != nil {
			packs.WriteJSONError(w, "could not load pack", http.StatusInternalServerError)
			return
		}
		var dependent string
		err := db.QueryRowContext(r.Context(), `SELECT d.addon_id FROM addon_dependencies d JOIN experience_pack_addons epa ON epa.experience_pack_id = ? AND epa.addon_id = d.addon_id WHERE d.dependency_id = ? LIMIT 1`, packID, addonID).Scan(&dependent)
		if err == nil {
			packs.WriteJSONError(w, "cannot remove an add-on required by addons/"+dependent, http.StatusBadRequest)
			return
		}
		if err != nil && !packs.NotFound(err) {
			packs.WriteJSONError(w, "could not check add-on dependencies", http.StatusInternalServerError)
			return
		}

		transaction, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			packs.WriteJSONError(w, "could not remove add-on", http.StatusInternalServerError)
			return
		}
		defer transaction.Rollback()
		result, err := transaction.ExecContext(r.Context(), `DELETE FROM experience_pack_addons WHERE experience_pack_id = ? AND addon_id = ?`, packID, addonID)
		if err != nil {
			packs.WriteJSONError(w, "could not remove add-on", http.StatusInternalServerError)
			return
		}
		count, err := result.RowsAffected()
		if err != nil {
			packs.WriteJSONError(w, "could not remove add-on", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			packs.WriteJSONError(w, "add-on is not in this pack", http.StatusNotFound)
			return
		}
		if err := setOrders(r.Context(), transaction, packID, selectedIDs(r.Context(), transaction, packID)); err != nil {
			packs.WriteJSONError(w, "could not update install order", http.StatusInternalServerError)
			return
		}
		if err := transaction.Commit(); err != nil {
			packs.WriteJSONError(w, "could not remove add-on", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
