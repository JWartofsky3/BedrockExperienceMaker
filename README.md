# Bedrock Experience Pack Builder

Curate Minecraft Bedrock add-ons into clear, shareable experience packs. This is a collection builder and install guide—not a modloader or an add-on hosting service.

## Product boundaries

- The application links to official provider pages and does not host, proxy, cache, or redistribute add-on archives.
- Players import individual `.mcaddon` files into Minecraft Bedrock and activate the relevant behavior and resource packs in their world.
- The initial experience is aimed at Windows and mobile Bedrock players; console installation is not a primary workflow.

## MVP functional requirements

1. **Browse add-ons**
   - Search by name, creator, tag, or category.
   - Filter by provider, Minecraft version, and content type.
   - Show a compact card with the add-on icon, title, creator, version, category, and important warnings.

2. **View an add-on**
   - Display its description, Bedrock versions, included behavior/resource packs, required experiments, caveats, source links, and cataloged releases.
   - Present a manifest-derived technical summary without exposing or serving the original archive.

3. **Open provider pages**
   - Offer the CurseForge and MCPEDL provider pages where available.
   - Attempt an iframe preview only when the provider permits embedding.
   - Always include an **Open official page** action, since either provider can block iframe embedding with browser security headers.

4. **Create an experience pack**
   - Provide a name, icon, description, and setup notes.
   - Add cataloged add-ons and select a recommended source/release for each one.
   - Remove, reorder, and annotate add-ons in the collection.

5. **Review compatibility and setup guidance**
   - Summarize required Bedrock experiments, activation instructions, known conflicts, dependencies, warnings, and the recommended order.
   - Generate a per-add-on installation checklist.

6. **View and share an experience pack**
   - Show the complete collection, chosen releases, source links, instructions, and warnings on a dedicated detail page.
   - Let players open each official provider page.
   - Export a small JSON install plan and/or printable checklist.

7. **Maintain a local catalog**
   - The MVP uses seed data maintained in the application: add-on metadata, provider pages, releases, manifest-derived details, compatibility notes, and warnings.
   - The MVP does not scrape providers or synchronize their catalogs automatically.

## Future: combined experience-pack export

A possible later capability is to download the selected add-ons and build a single `Experience Pack.mcaddon` archive. This is **not part of the MVP** and must be treated as an advanced, opt-in export workflow.

- The selected collection order must be persisted as a merge priority. Add-ons at the top are applied last and take precedence when files conflict.
- The catalog should retain each release's manifest identity, pack type, version, and a file inventory/checksum so a future exporter can detect conflicts before producing an archive.
- The exporter would need to handle behavior-pack and resource-pack files separately, preserve or intentionally reconcile manifest dependencies, and report every overridden file.
- A combined archive does not automatically make the packs compatible or activate them within a Minecraft world; the generated instructions still need to explain required world activation and experiments.
- Before implementation, we must confirm author permissions and provider terms. Downloading and repackaging third-party add-ons may be prohibited even if a direct download is technically possible.

## Local MySQL data model

The MVP schema is in [database/schema.sql](database/schema.sql), with brief relationship notes in [database/README.md](database/README.md). It is designed for MySQL 8.0+ and creates only a local `bedrock_experience_maker` database when imported.

It has a manually curated `addons` catalog, dependencies, experience packs with their ordered add-ons, plus local `users` and `sessions`. Packs display a creator name and are editable only by the signed-in creator. Source URLs, a current version, a Bedrock-version note, and raw manifest data are optional add-on metadata; no versions, manifest data, or compatibility records are required to browse add-ons or build a pack. Add-on IDs are stable, application-owned resource identifiers used directly in `addons/{id}`; they do not need to be UUIDs.

Creator-defined manifest UUIDs can be retained in raw manifest data but are never used to identify catalog records or enforce uniqueness.

### Availability and Bedrock-version policy

The catalog assumes a provider may expose only the newest artifact. It records an optional current version and always links to the provider page; historical-release tracking can be added later if it becomes useful.

Minecraft version tags and manifest minimum-engine versions are optional notes. Since players normally install the current Bedrock release, a mismatch should be displayed as a warning, never used to prevent adding an add-on to a pack.

## Run the add-on catalog

1. Configure the local MySQL Windows service as `MySQL`.
2. Run `npm run seed` from the repository root. It prompts for the local MySQL root password, then applies the schema, seed catalog, example experience, and local API account. It resets the add-ons in the seeded example pack to its documented order.
3. Run `npm run backend` from the repository root.
4. Run `npm run dev` from the repository root and open the URL Vite prints, normally `http://localhost:5173`. Use the seeded local creator account to sign in before creating or editing packs.

Anyone can browse an experience and use **Download JSON** to save its pack data, including the full JSON records for every installed add-on. Only the signed-in creator can change or delete a pack.

The Vite dev server proxies `GET /v1/addons` to the API at `http://127.0.0.1:8080`. The add-on page renders only the records returned by MySQL and shows an error if the API is unavailable.

`npm run full` starts MySQL, applies the local seed data, opens the API in a second Command Prompt window, and starts Vite in the current terminal. It prompts for the MySQL root password each time so that the seed data is refreshed. Close both terminal windows to stop the app.
