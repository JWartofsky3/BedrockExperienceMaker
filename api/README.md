# Local API

The API follows the standard AIP methods for a top-level `Addon` resource:

- `GET /v1/addons` lists add-ons, with optional `page_size` and `page_token` query parameters.
- `GET /v1/addons/{addon}` gets an add-on.
- `DELETE /v1/addons/{addon}` deletes an add-on and returns `404` when it does not exist. It returns `400` when an experience pack still references it (the REST representation of a failed precondition).

Packs are available at `/v1/packs`. Anyone can browse packs and their installed add-ons. Creating a pack and changing its add-ons requires a signed-in creator; only that creator can edit or delete it.

Local authentication uses these cookie-backed endpoints:

- `POST /v1/auth/login` accepts `username` and `password` and returns the current user.
- `POST /v1/auth/logout` clears the current browser session.
- `GET /v1/auth/me` returns the current user or `401` when signed out.

1. Import `../database/schema.sql` and `../database/seed_addons.sql` into local MySQL.
2. Set the values from `.env.example` in your shell.
3. From this folder, run `go run ./cmd/server`.

The server listens on `127.0.0.1:8080` by default. Vite proxies `/v1` to that address. The frontend reads only from this API; it does not contain a catalog fallback.
