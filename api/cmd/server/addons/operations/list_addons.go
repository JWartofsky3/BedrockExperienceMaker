package operations

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/addons"
)

const addonColumns = `id, display_name, creator_name, description, icon_path, curseforge_url, mcpedl_url, current_version, minecraft_version_note, manifest_data`
const defaultPageSize = 50
const maxPageSize = 100

func ListAddons(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageSize, err := parsePageSize(r.URL.Query().Get("page_size"))
		if err != nil {
			addons.WriteJSONError(w, "page_size must be a positive integer", http.StatusBadRequest)
			return
		}
		offset, err := parsePageToken(r.URL.Query().Get("page_token"))
		if err != nil {
			addons.WriteJSONError(w, "page_token is invalid", http.StatusBadRequest)
			return
		}

		rows, err := db.QueryContext(r.Context(), `SELECT `+addonColumns+` FROM addons ORDER BY display_name LIMIT ? OFFSET ?`, pageSize+1, offset)
		if err != nil {
			addons.WriteJSONError(w, "could not load add-ons", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := []addons.Addon{}
		hasNextPage := false
		for rows.Next() {
			if len(items) == pageSize {
				hasNextPage = true
				break
			}
			item, err := addons.ScanAddon(rows)
			if err != nil {
				addons.WriteJSONError(w, "could not read add-ons", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			addons.WriteJSONError(w, "could not load add-ons", http.StatusInternalServerError)
			return
		}
		response := struct {
			Addons        []addons.Addon `json:"addons"`
			NextPageToken string  `json:"nextPageToken,omitempty"`
		}{Addons: items}
		if hasNextPage {
			response.NextPageToken = encodePageToken(offset + pageSize)
		}
		addons.WriteJSON(w, response, http.StatusOK)
	}
}

func parsePageSize(value string) (int, error) {
	if value == "" {
		return defaultPageSize, nil
	}
	pageSize, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if pageSize < 1 {
		return 0, errors.New("page size must be positive")
	}
	if pageSize > maxPageSize {
		return maxPageSize, nil
	}
	return pageSize, nil
}

func parsePageToken(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return 0, err
	}
	offset, err := strconv.Atoi(string(decoded))
	if err != nil {
		return 0, err
	}
	if offset < 0 {
		return 0, errors.New("page offset must not be negative")
	}
	return offset, nil
}

func encodePageToken(offset int) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}
