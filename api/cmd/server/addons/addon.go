package addons

import (
	"database/sql"
	"encoding/json"
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
