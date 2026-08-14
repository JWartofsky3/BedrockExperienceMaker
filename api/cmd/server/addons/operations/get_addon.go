package operations

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/addons"
)

func GetAddon(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRowContext(r.Context(), `SELECT `+addonColumns+` FROM addons WHERE id = ?`, r.PathValue("name"))
		item, err := addons.ScanAddon(row)
		if errors.Is(err, sql.ErrNoRows) {
			addons.WriteJSONError(w, "add-on not found", http.StatusNotFound)
			return
		}
		if err != nil {
			addons.WriteJSONError(w, "could not load add-on", http.StatusInternalServerError)
			return
		}
		if err := addons.PopulateDependencies(r.Context(), db, &item); err != nil {
			addons.WriteJSONError(w, "could not load add-on dependencies", http.StatusInternalServerError)
			return
		}
		addons.WriteJSON(w, item, http.StatusOK)
	}
}
