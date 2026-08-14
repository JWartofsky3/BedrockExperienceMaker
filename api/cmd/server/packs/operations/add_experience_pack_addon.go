package operations

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

type addExperiencePackAddonRequest struct {
	AddonName string `json:"addonName"`
	Note      string `json:"note"`
}

func AddExperiencePackAddon(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request addExperiencePackAddonRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			packs.WriteJSONError(w, "request body is invalid", http.StatusBadRequest)
			return
		}
		addonID, valid := packs.AddonID(request.AddonName)
		if !valid {
			packs.WriteJSONError(w, "addonName must use the addons/{addon} format", http.StatusBadRequest)
			return
		}
		packID := r.PathValue("name")
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
		var exists bool
		if err := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM addons WHERE id = ?)`, addonID).Scan(&exists); err != nil {
			packs.WriteJSONError(w, "could not load add-on", http.StatusInternalServerError)
			return
		}
		if !exists {
			packs.WriteJSONError(w, "add-on not found", http.StatusNotFound)
			return
		}
		missing, err := missingDependencies(r, db, packID, addonID)
		if err != nil {
			packs.WriteJSONError(w, "could not check add-on dependencies", http.StatusInternalServerError)
			return
		}
		if len(missing) > 0 {
			packs.WriteJSONError(w, "required add-ons must be added first: "+missing[0], http.StatusBadRequest)
			return
		}

		transaction, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			packs.WriteJSONError(w, "could not add add-on", http.StatusInternalServerError)
			return
		}
		defer transaction.Rollback()
		var installOrder int
		if err := transaction.QueryRowContext(r.Context(), `SELECT COALESCE(MAX(install_order), 0) + 1 FROM experience_pack_addons WHERE experience_pack_id = ?`, packID).Scan(&installOrder); err != nil {
			packs.WriteJSONError(w, "could not add add-on", http.StatusInternalServerError)
			return
		}
		id, err := packs.NewID()
		if err != nil {
			packs.WriteJSONError(w, "could not add add-on", http.StatusInternalServerError)
			return
		}
		if _, err := transaction.ExecContext(r.Context(), `INSERT INTO experience_pack_addons (id, experience_pack_id, addon_id, install_order, note) VALUES (?, ?, ?, ?, ?)`, id, packID, addonID, installOrder, request.Note); err != nil {
			packs.WriteJSONError(w, "add-on is already in this pack", http.StatusBadRequest)
			return
		}
		if err := transaction.Commit(); err != nil {
			packs.WriteJSONError(w, "could not add add-on", http.StatusInternalServerError)
			return
		}
		pack, ok := loadPack(w, r, db)
		if ok {
			packs.WriteJSON(w, pack, http.StatusOK)
		}
	}
}

func missingDependencies(r *http.Request, db *sql.DB, packID, addonID string) ([]string, error) {
	rows, err := db.QueryContext(r.Context(), `SELECT d.dependency_id FROM addon_dependencies d LEFT JOIN experience_pack_addons epa ON epa.experience_pack_id = ? AND epa.addon_id = d.dependency_id WHERE d.addon_id = ? AND epa.addon_id IS NULL`, packID, addonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	missing := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		missing = append(missing, "addons/"+id)
	}
	return missing, rows.Err()
}
