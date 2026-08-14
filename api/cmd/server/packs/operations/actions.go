package operations

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

func HandleAction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value := r.PathValue("name")
		packID, action, found := strings.Cut(value, ":")
		if !found || packID == "" {
			packs.WriteJSONError(w, "pack action was not found", http.StatusNotFound)
			return
		}
		r.SetPathValue("name", packID)
		switch action {
		case "addAddon":
			AddExperiencePackAddon(db)(w, r)
		case "reorderAddons":
			ReorderExperiencePackAddons(db)(w, r)
		default:
			packs.WriteJSONError(w, "pack action was not found", http.StatusNotFound)
		}
	}
}
