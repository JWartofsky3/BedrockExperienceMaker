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

The initial schema is in [database/schema.sql](database/schema.sql), with relationship notes in [database/README.md](database/README.md). It is designed for MySQL 8.0+ and creates only a local `bedrock_experience_maker` database when imported.

The three inspected `.mcaddon` files demonstrate why the catalog separates an **add-on release** from its contained **manifest packs**. Each archive contains resource and behavior packs, and the behavior packs declare dependencies on resource-pack UUIDs and built-in script modules. Auto Miner also has an archive label of `1.0.2` while its behavior-pack header is `1.0.1`.

All primary keys are application-generated UUIDs owned by this application. Manifest header UUIDs and module UUIDs are stored as indexed metadata, never as unique keys: they are useful for dependency matching but cannot be trusted as globally unique creator identifiers. Manifest JSON is retained alongside normalized modules, dependencies, capabilities, and an optional file inventory so the model remains forward-compatible and can later support conflict detection.

### Availability and Bedrock-version policy

The catalog treats provider availability as a current observation, not a promise that every historical release can still be downloaded. A creator may expose only the latest artifact; therefore an unavailable historical source remains in a published plan for context, but must display an availability warning and offer the current project page where possible.

Minecraft version tags, manifest minimum-engine versions, and testing reports are compatibility evidence—not hard validation rules. Since players normally install the current Bedrock release, a mismatch should be surfaced with its source and severity, while still allowing the creator to include the add-on. Experience-pack revisions may record an intended Bedrock version so the app can produce useful warnings.
