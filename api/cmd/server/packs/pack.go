package packs

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"crypto/rand"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/addons"
)

const addonColumns = `a.id, a.display_name, a.creator_name, a.description, a.icon_path, a.curseforge_url, a.mcpedl_url, a.current_version, a.minecraft_version_note, a.manifest_data`

type ExperiencePack struct {
	Name          string      `json:"name"`
	DisplayName   string      `json:"displayName"`
	CreatorName   string      `json:"creatorName"`
	CreatorUserID string      `json:"creatorUserId,omitempty"`
	Description   string      `json:"description,omitempty"`
	SetupNotes    string      `json:"setupNotes,omitempty"`
	Addons        []PackAddon `json:"addons,omitempty"`
}

type PackAddon struct {
	Addon        addons.Addon `json:"addon"`
	InstallOrder int          `json:"installOrder"`
	Note         string       `json:"note,omitempty"`
}

func Get(ctx context.Context, db *sql.DB, id string, includeAddons bool) (ExperiencePack, error) {
	row := db.QueryRowContext(ctx, `SELECT p.id, p.name, COALESCE(u.username, p.creator_name), p.creator_user_id, p.description, p.setup_notes FROM experience_packs p LEFT JOIN users u ON u.id = p.creator_user_id WHERE p.id = ?`, id)
	var pack ExperiencePack
	var storedID string
	var creator, creatorUserID, description, setupNotes sql.NullString
	if err := row.Scan(&storedID, &pack.DisplayName, &creator, &creatorUserID, &description, &setupNotes); err != nil {
		return ExperiencePack{}, err
	}
	pack.Name = "packs/" + storedID
	pack.CreatorName = creator.String
	pack.CreatorUserID = creatorUserID.String
	pack.Description = description.String
	pack.SetupNotes = setupNotes.String
	if includeAddons {
		items, err := ListAddons(ctx, db, id)
		if err != nil {
			return ExperiencePack{}, err
		}
		pack.Addons = items
	}
	return pack, nil
}

func IsCreator(ctx context.Context, db *sql.DB, packID, userID string) (bool, error) {
	var creatorUserID sql.NullString
	if err := db.QueryRowContext(ctx, `SELECT creator_user_id FROM experience_packs WHERE id = ?`, packID).Scan(&creatorUserID); err != nil {
		return false, err
	}
	return creatorUserID.Valid && creatorUserID.String == userID, nil
}

func ListAddons(ctx context.Context, db *sql.DB, packID string) ([]PackAddon, error) {
	rows, err := db.QueryContext(ctx, `SELECT epa.install_order, epa.note, `+addonColumns+` FROM experience_pack_addons epa JOIN addons a ON a.id = epa.addon_id WHERE epa.experience_pack_id = ? ORDER BY epa.install_order`, packID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []PackAddon{}
	for rows.Next() {
		var item PackAddon
		var id string
		var creator, description, iconPath, curseForgeURL, mcpedlURL, version, minecraftNote, note sql.NullString
		var manifestData []byte
		if err := rows.Scan(&item.InstallOrder, &note, &id, &item.Addon.DisplayName, &creator, &description, &iconPath, &curseForgeURL, &mcpedlURL, &version, &minecraftNote, &manifestData); err != nil {
			return nil, err
		}
		item.Addon.Name = "addons/" + id
		item.Addon.CreatorName = creator.String
		item.Addon.Description = description.String
		item.Addon.IconPath = iconPath.String
		item.Addon.CurseForgeURL = curseForgeURL.String
		item.Addon.MCPEDLURL = mcpedlURL.String
		item.Addon.CurrentVersion = version.String
		item.Addon.MinecraftVersionNote = minecraftNote.String
		item.Addon.ManifestData = manifestData
		item.Note = note.String
		if err := addons.PopulateDependencies(ctx, db, &item.Addon); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func ID(name string) (string, bool) {
	if !strings.HasPrefix(name, "packs/") {
		return "", false
	}
	id := strings.TrimPrefix(name, "packs/")
	return id, id != "" && !strings.Contains(id, "/")
}

func AddonID(name string) (string, bool) {
	if !strings.HasPrefix(name, "addons/") {
		return "", false
	}
	id := strings.TrimPrefix(name, "addons/")
	return id, id != "" && !strings.Contains(id, "/")
}

func NewID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	value := hex.EncodeToString(bytes)
	return value[0:8] + "-" + value[8:12] + "-" + value[12:16] + "-" + value[16:20] + "-" + value[20:], nil
}

func Slug(displayName, id string) string {
	var builder strings.Builder
	lastDash := false
	for _, character := range strings.ToLower(displayName) {
		if (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9') {
			builder.WriteRune(character)
			lastDash = false
		} else if !lastDash {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	value := strings.Trim(builder.String(), "-")
	if value == "" {
		value = "pack"
	}
	if len(value) > 150 {
		value = value[:150]
	}
	return value + "-" + id[:8]
}

func NotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
