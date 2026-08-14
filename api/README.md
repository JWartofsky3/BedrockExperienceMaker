# API design

This folder defines the resource-oriented HTTP API that will sit between the React app and the local MySQL database. It is a design contract only; no server is implemented yet.

Base path: `/v1`

All resource names use the application-owned UUID from the database. Minecraft manifest UUIDs never appear in resource names.

## Resources

### Addon

An `Addon` is a manually curated catalog entry.

- Resource type: `bedrockexperiencemaker.local/Addon`
- Resource name: `addons/{addon}`
- Example name: `addons/2e839411-4a0e-4cee-b106-33d5bde91193`

Core fields: `name`, `displayName`, `creatorName`, `description`, `iconPath`, `curseforgeUrl`, `mcpedlUrl`, `currentVersion`, `minecraftVersionNote`, and optional `manifestData`.

### ExperiencePack

An `ExperiencePack` is a named collection with a description and setup notes. It owns its add-on membership resources.

- Resource type: `bedrockexperiencemaker.local/ExperiencePack`
- Resource name: `experiencePacks/{experience_pack}`
- Example name: `experiencePacks/9f026c09-8388-42fb-889a-2fb8ce2d1b2c`

Core fields: `name`, `displayName`, `description`, `setupNotes`, and `iconPath`.

### ExperiencePackAddon

An `ExperiencePackAddon` is an add-on selected for one experience pack. It is a child resource because installation order and the selected source belong to the relationship, not to the catalog add-on itself.

- Resource type: `bedrockexperiencemaker.local/ExperiencePackAddon`
- Resource name: `experiencePacks/{experience_pack}/addons/{addon}`
- Example name: `experiencePacks/9f026c09-8388-42fb-889a-2fb8ce2d1b2c/addons/2e839411-4a0e-4cee-b106-33d5bde91193`

Core fields: `name`, `addon` (for example, `addons/{addon}`), `installOrder`, `selectedSource`, `selectedVersion`, and `note`.

The child resource's final identifier is the catalog add-on ID. This directly matches the MVP rule that an add-on appears at most once in a pack.

## Standard methods

The API uses resource `name` fields for standard `Get` and `Delete` requests, as prescribed by AIP-131 and AIP-135. `List` methods include `pageSize`, `pageToken`, and `nextPageToken`, as prescribed by AIP-132. The pack-creation methods follow AIP-133.

| Method | HTTP route | Purpose |
| --- | --- | --- |
| `ListAddons` | `GET /v1/addons` | Browse cataloged add-ons. Supports optional `page_size`, `page_token`, and `filter`. Default order: `displayName` ascending. |
| `GetAddon` | `GET /v1/{name=addons/*}` | Return one add-on. |
| `ListExperiencePacks` | `GET /v1/experiencePacks` | List experience packs. Supports pagination. |
| `GetExperiencePack` | `GET /v1/{name=experiencePacks/*}` | Return one pack, without embedding child add-ons. |
| `CreateExperiencePack` | `POST /v1/experiencePacks` | Create a pack. The body contains `experiencePack`; optional `experience_pack_id` lets a client choose its UUID. |
| `DeleteExperiencePack` | `DELETE /v1/{name=experiencePacks/*}` | Delete a pack and its child membership records. |
| `ListExperiencePackAddons` | `GET /v1/{parent=experiencePacks/*}/addons` | Return a pack's selected add-ons, ordered by `installOrder` ascending. |
| `GetExperiencePackAddon` | `GET /v1/{name=experiencePacks/*/addons/*}` | Return one selected add-on in a pack. |
| `CreateExperiencePackAddon` | `POST /v1/{parent=experiencePacks/*}/addons` | Add a catalog add-on to a pack. The body contains `experiencePackAddon`; optional `addon_id` is the selected catalog add-on ID. |
| `DeleteExperiencePackAddon` | `DELETE /v1/{name=experiencePacks/*/addons/*}` | Remove an add-on from a pack. |

`ListAddons` may begin with a simple name/creator text filter. No advanced filter language or client-controlled ordering is part of the first implementation.

## Response and error rules

- `Get*` returns the resource directly, not a wrapper object.
- `List*` returns the plural resource field (`addons` or `experiencePacks`) and `nextPageToken`.
- `Delete*` returns HTTP `204 No Content` on success.
- A missing named resource returns `404 Not Found`.
- Invalid UUIDs, page sizes, or invalid resource relationships return `400 Bad Request`.
- A pack request naming an add-on that is already present returns `409 Conflict`.

## Deliberately deferred

- Add-on write endpoints: catalog entries will be inserted manually in the database during the MVP.
- Update endpoints, pack revisions, sharing, accounts, and access control.
- Version enforcement, compatibility enforcement, archive download, and combined archive generation.
