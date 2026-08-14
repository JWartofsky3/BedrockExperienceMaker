package operations

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

type reorderExperiencePackAddonsRequest struct {
	AddonNames []string `json:"addonNames"`
}

func ReorderExperiencePackAddons(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request reorderExperiencePackAddonsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			packs.WriteJSONError(w, "request body is invalid", http.StatusBadRequest)
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

		ids := make([]string, 0, len(request.AddonNames))
		seen := map[string]bool{}
		for _, name := range request.AddonNames {
			id, valid := packs.AddonID(name)
			if !valid || seen[id] {
				packs.WriteJSONError(w, "addonNames must contain unique addons/{addon} values", http.StatusBadRequest)
				return
			}
			seen[id] = true
			ids = append(ids, id)
		}
		current := selectedIDs(r.Context(), db, packID)
		if len(ids) != len(current) || !sameAddons(ids, current) {
			packs.WriteJSONError(w, "addonNames must include every add-on in the pack exactly once", http.StatusBadRequest)
			return
		}
		if err := validateDependencyOrder(r.Context(), db, packID, ids); err != nil {
			packs.WriteJSONError(w, "dependencies must be installed before the add-ons that require them", http.StatusBadRequest)
			return
		}
		transaction, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			packs.WriteJSONError(w, "could not reorder add-ons", http.StatusInternalServerError)
			return
		}
		defer transaction.Rollback()
		if err := setOrders(r.Context(), transaction, packID, ids); err != nil {
			packs.WriteJSONError(w, "could not reorder add-ons", http.StatusInternalServerError)
			return
		}
		if err := transaction.Commit(); err != nil {
			packs.WriteJSONError(w, "could not reorder add-ons", http.StatusInternalServerError)
			return
		}
		pack, ok := loadPack(w, r, db)
		if ok {
			packs.WriteJSON(w, pack, http.StatusOK)
		}
	}
}

func selectedIDs(ctx context.Context, queryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, packID string) []string {
	rows, err := queryer.QueryContext(ctx, `SELECT addon_id FROM experience_pack_addons WHERE experience_pack_id = ? ORDER BY install_order`, packID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if rows.Scan(&id) != nil {
			return nil
		}
		ids = append(ids, id)
	}
	return ids
}

func sameAddons(left, right []string) bool {
	values := map[string]bool{}
	for _, value := range left {
		values[value] = true
	}
	for _, value := range right {
		if !values[value] {
			return false
		}
	}
	return true
}

func validateDependencyOrder(ctx context.Context, db *sql.DB, packID string, ids []string) error {
	positions := map[string]int{}
	for index, id := range ids {
		positions[id] = index
	}
	rows, err := db.QueryContext(ctx, `SELECT d.addon_id, d.dependency_id FROM addon_dependencies d JOIN experience_pack_addons dependent ON dependent.experience_pack_id = ? AND dependent.addon_id = d.addon_id JOIN experience_pack_addons dependency ON dependency.experience_pack_id = ? AND dependency.addon_id = d.dependency_id`, packID, packID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var addonID, dependencyID string
		if err := rows.Scan(&addonID, &dependencyID); err != nil {
			return err
		}
		if positions[dependencyID] > positions[addonID] {
			return sql.ErrNoRows
		}
	}
	return rows.Err()
}

func setOrders(ctx context.Context, transaction *sql.Tx, packID string, ids []string) error {
	if _, err := transaction.ExecContext(ctx, `UPDATE experience_pack_addons SET install_order = install_order + 1000 WHERE experience_pack_id = ?`, packID); err != nil {
		return err
	}
	for index, id := range ids {
		if _, err := transaction.ExecContext(ctx, `UPDATE experience_pack_addons SET install_order = ? WHERE experience_pack_id = ? AND addon_id = ?`, index+1, packID, id); err != nil {
			return err
		}
	}
	return nil
}
