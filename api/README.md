# Local API

The API follows the standard AIP methods for a top-level `Addon` resource:

- `GET /v1/addons` lists add-ons, with optional `page_size` and `page_token` query parameters.
- `GET /v1/addons/{addon}` gets an add-on.
- `DELETE /v1/addons/{addon}` deletes an add-on and returns `404` when it does not exist. It returns `400` when an experience pack still references it (the REST representation of a failed precondition).

Packs are available at `/v1/packs`. The API supports browsing, getting, creating, and deleting packs, plus adding, removing, and reordering their add-ons. Pack creation is local and unauthenticated; `creator_name` is display metadata only.

1. Import `../database/schema.sql` and `../database/seed_addons.sql` into local MySQL.
2. Set the values from `.env.example` in your shell.
3. From this folder, run `go run ./cmd/server`.

The server listens on `127.0.0.1:8080` by default. Vite proxies `/v1` to that address. The frontend reads only from this API; it does not contain a catalog fallback.
