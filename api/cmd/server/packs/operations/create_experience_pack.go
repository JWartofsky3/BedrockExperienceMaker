package operations

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/auth"
	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

type createExperiencePackRequest struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	SetupNotes  string `json:"setupNotes"`
}

func CreateExperiencePack(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.RequireUser(w, r)
		if !ok {
			return
		}
		var request createExperiencePackRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			packs.WriteJSONError(w, "request body is invalid", http.StatusBadRequest)
			return
		}
		request.DisplayName = strings.TrimSpace(request.DisplayName)
		if request.DisplayName == "" {
			packs.WriteJSONError(w, "displayName is required", http.StatusBadRequest)
			return
		}

		id, err := packs.NewID()
		if err != nil {
			packs.WriteJSONError(w, "could not create pack", http.StatusInternalServerError)
			return
		}
		_, err = db.ExecContext(r.Context(), `INSERT INTO experience_packs (id, slug, name, creator_name, creator_user_id, description, setup_notes) VALUES (?, ?, ?, ?, ?, ?, ?)`, id, packs.Slug(request.DisplayName, id), request.DisplayName, user.Username, user.ID, request.Description, request.SetupNotes)
		if err != nil {
			packs.WriteJSONError(w, "could not create pack", http.StatusInternalServerError)
			return
		}
		packs.WriteJSON(w, packs.ExperiencePack{
			Name:          "packs/" + id,
			DisplayName:   request.DisplayName,
			CreatorName:   user.Username,
			CreatorUserID: user.ID,
			Description:   request.Description,
			SetupNotes:    request.SetupNotes,
			Addons:        []packs.PackAddon{},
		}, http.StatusCreated)
	}
}
