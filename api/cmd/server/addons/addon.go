package addons

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
)

type Addon struct {
	Name                 string          `json:"name"`
	DisplayName          string          `json:"displayName"`
	CreatorName          string          `json:"creatorName,omitempty"`
	Description          string          `json:"description,omitempty"`
	IconPath             string          `json:"iconPath,omitempty"`
	CurseForgeURL        string          `json:"curseforgeUrl,omitempty"`
	MCPEDLURL            string          `json:"mcpedlUrl,omitempty"`
	CurrentVersion       string          `json:"currentVersion,omitempty"`
	MinecraftVersionNote string          `json:"minecraftVersionNote,omitempty"`
	ManifestData         json.RawMessage `json:"manifestData,omitempty"`
	Dependencies         []Dependency    `json:"dependencies,omitempty"`
}

type Dependency struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

func ScanAddon(scanner interface{ Scan(...any) error }) (Addon, error) {
	var item Addon
	var id string
	var creator, description, iconPath, curseForgeURL, mcpedlURL, version, minecraftNote sql.NullString
	var manifestData []byte
	err := scanner.Scan(&id, &item.DisplayName, &creator, &description, &iconPath, &curseForgeURL, &mcpedlURL, &version, &minecraftNote, &manifestData)
	if err != nil {
		return Addon{}, err
	}
	item.Name = "addons/" + id
	item.CreatorName = creator.String
	item.Description = description.String
	item.IconPath = iconPath.String
	item.CurseForgeURL = curseForgeURL.String
	item.MCPEDLURL = mcpedlURL.String
	item.CurrentVersion = version.String
	item.MinecraftVersionNote = minecraftNote.String
	item.ManifestData = manifestData
	return item, nil
}

func PopulateDependencies(ctx context.Context, db *sql.DB, item *Addon) error {
	addonID := strings.TrimPrefix(item.Name, "addons/")
	rows, err := db.QueryContext(ctx, `SELECT d.dependency_id, a.display_name FROM addon_dependencies d JOIN addons a ON a.id = d.dependency_id WHERE d.addon_id = ? ORDER BY a.display_name`, addonID)
	if err != nil {
		return err
	}
	defer rows.Close()

	item.Dependencies = []Dependency{}
	for rows.Next() {
		var dependency Dependency
		var id string
		if err := rows.Scan(&id, &dependency.DisplayName); err != nil {
			return err
		}
		dependency.Name = "addons/" + id
		item.Dependencies = append(item.Dependencies, dependency)
	}
	return rows.Err()
}
