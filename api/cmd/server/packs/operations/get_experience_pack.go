package operations

import (
	"database/sql"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

func GetExperiencePack(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pack, ok := loadPack(w, r, db)
		if !ok {
			return
		}
		packs.WriteJSON(w, pack, http.StatusOK)
	}
}

func loadPack(w http.ResponseWriter, r *http.Request, db *sql.DB) (packs.ExperiencePack, bool) {
	id := r.PathValue("name")
	pack, err := packs.Get(r.Context(), db, id, true)
	if packs.NotFound(err) {
		packs.WriteJSONError(w, "pack not found", http.StatusNotFound)
		return packs.ExperiencePack{}, false
	}
	if err != nil {
		packs.WriteJSONError(w, "could not load pack", http.StatusInternalServerError)
		return packs.ExperiencePack{}, false
	}
	return pack, true
}
