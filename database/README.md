# Local database

`schema.sql` is the small, local MySQL 8.0+ schema for the MVP. It creates the `bedrock_experience_maker` database when imported.

## Tables

```text
addons ──< addon_dependencies

users ──< sessions
  │
  └──< experience_packs ──< experience_pack_addons >── addons

addons ──< addon_dependencies
```

- `addons` is the manually curated catalog. Its `id` is the stable resource ID used in `addons/{id}`; it includes the page URLs needed by the UI plus optional current-version, Bedrock-version note, and raw manifest metadata.
- `addon_dependencies` records required catalog add-ons; both add-ons must exist before the relationship can be saved.
- `experience_packs` stores a named collection and its setup notes.
- `experience_pack_addons` joins the two with an explicit installation order and optional source/version choice.
- `users` stores local creator identities with bcrypt password hashes; `sessions` stores hashed, expiring browser-session tokens.

IDs are application-owned, unique resource identifiers and do not need to be UUIDs. Creator-defined manifest UUIDs, pack data, and version metadata can be stored in `addons.manifest_data` but are not used to identify catalog records or enforce compatibility.
