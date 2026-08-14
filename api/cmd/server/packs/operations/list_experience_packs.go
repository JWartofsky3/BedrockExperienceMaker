package operations

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
)

const defaultPageSize = 50
const maxPageSize = 100

func ListExperiencePacks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageSize, err := parsePageSize(r.URL.Query().Get("page_size"))
		if err != nil {
			packs.WriteJSONError(w, "page_size must be a positive integer", http.StatusBadRequest)
			return
		}
		offset, err := parsePageToken(r.URL.Query().Get("page_token"))
		if err != nil {
			packs.WriteJSONError(w, "page_token is invalid", http.StatusBadRequest)
			return
		}
		rows, err := db.QueryContext(r.Context(), `SELECT id FROM experience_packs ORDER BY created_at DESC LIMIT ? OFFSET ?`, pageSize+1, offset)
		if err != nil {
			packs.WriteJSONError(w, "could not load packs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := []packs.ExperiencePack{}
		hasNextPage := false
		for rows.Next() {
			if len(items) == pageSize {
				hasNextPage = true
				break
			}
			var id string
			if err := rows.Scan(&id); err != nil {
				packs.WriteJSONError(w, "could not read packs", http.StatusInternalServerError)
				return
			}
			item, err := packs.Get(r.Context(), db, id, false)
			if err != nil {
				packs.WriteJSONError(w, "could not read packs", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			packs.WriteJSONError(w, "could not load packs", http.StatusInternalServerError)
			return
		}
		packs.WriteJSON(w, struct {
			ExperiencePacks []packs.ExperiencePack `json:"experiencePacks"`
			NextPageToken   string                 `json:"nextPageToken,omitempty"`
		}{ExperiencePacks: items, NextPageToken: nextPageToken(offset, pageSize, hasNextPage)}, http.StatusOK)
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
	if err != nil || offset < 0 {
		return 0, errors.New("page token is invalid")
	}
	return offset, nil
}

func nextPageToken(offset, pageSize int, hasNextPage bool) string {
	if !hasNextPage {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(offset + pageSize)))
}
