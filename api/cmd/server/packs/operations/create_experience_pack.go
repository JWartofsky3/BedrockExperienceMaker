package operations

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

type createExperiencePackRequest struct {
	DisplayName string `json:"displayName"`
	CreatorName string `json:"creatorName"`
	Description string `json:"description"`
	SetupNotes  string `json:"setupNotes"`
}

func CreateExperiencePack(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request createExperiencePackRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			packs.WriteJSONError(w, "request body is invalid", http.StatusBadRequest)
			return
		}
		request.DisplayName = strings.TrimSpace(request.DisplayName)
		request.CreatorName = strings.TrimSpace(request.CreatorName)
		if request.DisplayName == "" || request.CreatorName == "" {
			packs.WriteJSONError(w, "displayName and creatorName are required", http.StatusBadRequest)
			return
		}

		id, err := packs.NewID()
		if err != nil {
			packs.WriteJSONError(w, "could not create pack", http.StatusInternalServerError)
			return
		}
		_, err = db.ExecContext(r.Context(), `INSERT INTO experience_packs (id, slug, name, creator_name, description, setup_notes) VALUES (?, ?, ?, ?, ?, ?)`, id, packs.Slug(request.DisplayName, id), request.DisplayName, request.CreatorName, request.Description, request.SetupNotes)
		if err != nil {
			packs.WriteJSONError(w, "could not create pack", http.StatusInternalServerError)
			return
		}
		packs.WriteJSON(w, packs.ExperiencePack{
			Name:        "packs/" + id,
			DisplayName: request.DisplayName,
			CreatorName: request.CreatorName,
			Description: request.Description,
			SetupNotes:  request.SetupNotes,
			Addons:      []packs.PackAddon{},
		}, http.StatusCreated)
	}
}
