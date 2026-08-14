-- Bedrock Experience Pack Builder - local MySQL 8.0+ schema
-- Primary keys are application-generated UUIDs (CHAR(36)).
-- Minecraft manifest UUIDs are deliberately indexed but never UNIQUE.

CREATE DATABASE IF NOT EXISTS bedrock_experience_maker
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE bedrock_experience_maker;

CREATE TABLE creators (
  id CHAR(36) NOT NULL,
  display_name VARCHAR(160) NOT NULL,
  profile_url VARCHAR(2048) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_creators_display_name (display_name)
) ENGINE=InnoDB;

CREATE TABLE addons (
  id CHAR(36) NOT NULL,
  slug VARCHAR(160) NOT NULL,
  name VARCHAR(255) NOT NULL,
  short_description VARCHAR(500) NULL,
  description TEXT NULL,
  icon_path VARCHAR(1024) NULL,
  catalog_status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_addons_slug (slug),
  KEY ix_addons_name (name)
) ENGINE=InnoDB;

CREATE TABLE addon_creators (
  addon_id CHAR(36) NOT NULL,
  creator_id CHAR(36) NOT NULL,
  credit_role ENUM('author', 'contributor', 'publisher') NOT NULL DEFAULT 'author',
  PRIMARY KEY (addon_id, creator_id),
  CONSTRAINT fk_addon_creators_addon FOREIGN KEY (addon_id) REFERENCES addons (id) ON DELETE CASCADE,
  CONSTRAINT fk_addon_creators_creator FOREIGN KEY (creator_id) REFERENCES creators (id) ON DELETE RESTRICT
) ENGINE=InnoDB;

CREATE TABLE tags (
  id CHAR(36) NOT NULL,
  slug VARCHAR(80) NOT NULL,
  label VARCHAR(100) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_tags_slug (slug)
) ENGINE=InnoDB;

CREATE TABLE addon_tags (
  addon_id CHAR(36) NOT NULL,
  tag_id CHAR(36) NOT NULL,
  PRIMARY KEY (addon_id, tag_id),
  CONSTRAINT fk_addon_tags_addon FOREIGN KEY (addon_id) REFERENCES addons (id) ON DELETE CASCADE,
  CONSTRAINT fk_addon_tags_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Provider/project pages, e.g. the CurseForge and MCPEDL page for one add-on.
CREATE TABLE addon_sources (
  id CHAR(36) NOT NULL,
  addon_id CHAR(36) NOT NULL,
  provider ENUM('curseforge', 'mcpedl', 'creator_site', 'other') NOT NULL,
  provider_project_id VARCHAR(128) NULL,
  project_slug VARCHAR(255) NULL,
  project_url VARCHAR(2048) NOT NULL,
  embed_url VARCHAR(2048) NULL,
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_addon_sources_provider_project (provider, provider_project_id),
  UNIQUE KEY uq_addon_sources_url (project_url(255)),
  KEY ix_addon_sources_addon (addon_id),
  CONSTRAINT fk_addon_sources_addon FOREIGN KEY (addon_id) REFERENCES addons (id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- A curator's logical release. Do not infer this solely from a manifest header.
CREATE TABLE addon_releases (
  id CHAR(36) NOT NULL,
  addon_id CHAR(36) NOT NULL,
  version_label VARCHAR(100) NOT NULL,
  release_notes TEXT NULL,
  published_at DATETIME NULL,
  release_status ENUM('current', 'historical', 'unknown') NOT NULL DEFAULT 'current',
  is_recommended BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_addon_releases_version (addon_id, version_label),
  KEY ix_addon_releases_addon (addon_id),
  CONSTRAINT fk_addon_releases_addon FOREIGN KEY (addon_id) REFERENCES addons (id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- A downloadable artifact or project-page reference for a particular release.
-- It records provenance; the application never stores or serves the archive itself.
CREATE TABLE addon_release_sources (
  id CHAR(36) NOT NULL,
  addon_release_id CHAR(36) NOT NULL,
  addon_source_id CHAR(36) NOT NULL,
  source_version_label VARCHAR(100) NULL,
  provider_file_id VARCHAR(128) NULL,
  page_url VARCHAR(2048) NOT NULL,
  direct_download_url VARCHAR(2048) NULL,
  archive_filename VARCHAR(512) NULL,
  archive_size_bytes BIGINT UNSIGNED NULL,
  archive_sha256 BINARY(32) NULL,
  availability_status ENUM('available', 'project_page_only', 'unavailable', 'unknown') NOT NULL DEFAULT 'unknown',
  availability_checked_at DATETIME NULL,
  inspected_at DATETIME NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_release_source_release_pair (id, addon_release_id),
  UNIQUE KEY uq_release_source_provider_file (addon_source_id, provider_file_id),
  KEY ix_release_sources_release (addon_release_id),
  CONSTRAINT fk_release_sources_release FOREIGN KEY (addon_release_id) REFERENCES addon_releases (id) ON DELETE CASCADE,
  CONSTRAINT fk_release_sources_source FOREIGN KEY (addon_source_id) REFERENCES addon_sources (id) ON DELETE RESTRICT
) ENGINE=InnoDB;

-- One .mcaddon can contain any number of behavior/resource/skin/world packs.
CREATE TABLE release_packs (
  id CHAR(36) NOT NULL,
  addon_release_id CHAR(36) NOT NULL,
  inspected_from_release_source_id CHAR(36) NULL,
  archive_path VARCHAR(1024) NOT NULL,
  pack_kind ENUM('behavior', 'resource', 'skin', 'world_template', 'other') NOT NULL,
  manifest_format_version SMALLINT UNSIGNED NULL,
  manifest_header_name VARCHAR(255) NOT NULL,
  manifest_header_description TEXT NULL,
  manifest_header_uuid CHAR(36) NULL,
  manifest_version VARCHAR(100) NULL,
  min_engine_version VARCHAR(100) NULL,
  pack_scope ENUM('world', 'global', 'any') NULL,
  product_type VARCHAR(100) NULL,
  manifest_json JSON NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_release_packs_path (addon_release_id, archive_path),
  KEY ix_release_packs_release (addon_release_id),
  KEY ix_release_packs_header_uuid (manifest_header_uuid),
  CONSTRAINT fk_release_packs_release FOREIGN KEY (addon_release_id) REFERENCES addon_releases (id) ON DELETE CASCADE,
  CONSTRAINT fk_release_packs_inspected_source FOREIGN KEY (inspected_from_release_source_id) REFERENCES addon_release_sources (id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE release_pack_modules (
  id CHAR(36) NOT NULL,
  release_pack_id CHAR(36) NOT NULL,
  module_index SMALLINT UNSIGNED NOT NULL,
  module_type VARCHAR(64) NOT NULL,
  module_uuid CHAR(36) NULL,
  module_version VARCHAR(100) NULL,
  description TEXT NULL,
  language VARCHAR(64) NULL,
  entry_path VARCHAR(1024) NULL,
  module_json JSON NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_release_pack_modules_position (release_pack_id, module_index),
  KEY ix_release_pack_modules_uuid (module_uuid),
  CONSTRAINT fk_release_pack_modules_pack FOREIGN KEY (release_pack_id) REFERENCES release_packs (id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Dependencies can target a creator-defined pack UUID or a built-in module name.
CREATE TABLE release_pack_dependencies (
  id CHAR(36) NOT NULL,
  release_pack_id CHAR(36) NOT NULL,
  dependency_index SMALLINT UNSIGNED NOT NULL,
  dependency_kind ENUM('pack_uuid', 'module_name') NOT NULL,
  dependency_uuid CHAR(36) NULL,
  dependency_module_name VARCHAR(255) NULL,
  required_version VARCHAR(100) NULL,
  dependency_json JSON NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_release_pack_dependencies_position (release_pack_id, dependency_index),
  KEY ix_dependencies_uuid (dependency_uuid),
  KEY ix_dependencies_module (dependency_module_name),
  CONSTRAINT fk_release_pack_dependencies_pack FOREIGN KEY (release_pack_id) REFERENCES release_packs (id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE release_pack_capabilities (
  release_pack_id CHAR(36) NOT NULL,
  capability VARCHAR(100) NOT NULL,
  PRIMARY KEY (release_pack_id, capability),
  CONSTRAINT fk_release_pack_capabilities_pack FOREIGN KEY (release_pack_id) REFERENCES release_packs (id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Inventory only: this supports future conflict detection without retaining add-on files.
CREATE TABLE release_pack_files (
  id CHAR(36) NOT NULL,
  release_pack_id CHAR(36) NOT NULL,
  relative_path VARCHAR(1024) NOT NULL,
  size_bytes BIGINT UNSIGNED NOT NULL,
  content_sha256 BINARY(32) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_release_pack_files_path (release_pack_id, relative_path),
  KEY ix_release_pack_files_path (relative_path(255)),
  CONSTRAINT fk_release_pack_files_pack FOREIGN KEY (release_pack_id) REFERENCES release_packs (id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE release_notes (
  id CHAR(36) NOT NULL,
  addon_release_id CHAR(36) NOT NULL,
  note_kind ENUM('warning', 'requirement', 'activation', 'compatibility', 'general') NOT NULL,
  severity ENUM('info', 'warning', 'critical') NOT NULL DEFAULT 'info',
  body TEXT NOT NULL,
  sort_order SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY ix_release_notes_release (addon_release_id),
  CONSTRAINT fk_release_notes_release FOREIGN KEY (addon_release_id) REFERENCES addon_releases (id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE experiments (
  id CHAR(36) NOT NULL,
  code VARCHAR(100) NOT NULL,
  label VARCHAR(255) NOT NULL,
  description TEXT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_experiments_code (code)
) ENGINE=InnoDB;

CREATE TABLE addon_release_experiments (
  addon_release_id CHAR(36) NOT NULL,
  experiment_id CHAR(36) NOT NULL,
  requirement ENUM('required', 'recommended') NOT NULL DEFAULT 'required',
  PRIMARY KEY (addon_release_id, experiment_id),
  CONSTRAINT fk_release_experiments_release FOREIGN KEY (addon_release_id) REFERENCES addon_releases (id) ON DELETE CASCADE,
  CONSTRAINT fk_release_experiments_experiment FOREIGN KEY (experiment_id) REFERENCES experiments (id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Compatibility is advisory. Version tags and manifest requirements are evidence,
-- not a hard install block, because Bedrock players normally use the current game build.
CREATE TABLE addon_release_compatibility_evidence (
  id CHAR(36) NOT NULL,
  addon_release_id CHAR(36) NOT NULL,
  minecraft_version VARCHAR(100) NOT NULL,
  evidence_type ENUM('manifest_minimum', 'provider_tag', 'creator_claim', 'curator_test', 'user_report') NOT NULL,
  assessment ENUM('unknown', 'likely_compatible', 'verified_compatible', 'warning', 'reported_incompatible') NOT NULL DEFAULT 'unknown',
  note TEXT NULL,
  observed_at DATETIME NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY ix_release_compatibility_release (addon_release_id),
  KEY ix_release_compatibility_version (minecraft_version),
  CONSTRAINT fk_release_compatibility_release FOREIGN KEY (addon_release_id) REFERENCES addon_releases (id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE release_compatibility_rules (
  id CHAR(36) NOT NULL,
  left_addon_release_id CHAR(36) NOT NULL,
  right_addon_release_id CHAR(36) NOT NULL,
  rule_type ENUM('incompatible', 'requires', 'recommended_before', 'recommended_after', 'warning') NOT NULL,
  severity ENUM('info', 'warning', 'critical') NOT NULL DEFAULT 'warning',
  note TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_compatibility_rule (left_addon_release_id, right_addon_release_id, rule_type),
  CONSTRAINT fk_compatibility_left_release FOREIGN KEY (left_addon_release_id) REFERENCES addon_releases (id) ON DELETE CASCADE,
  CONSTRAINT fk_compatibility_right_release FOREIGN KEY (right_addon_release_id) REFERENCES addon_releases (id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Stable identity for a creator's collection. No account model is needed locally yet.
CREATE TABLE experience_packs (
  id CHAR(36) NOT NULL,
  slug VARCHAR(160) NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_by_name VARCHAR(160) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_experience_packs_slug (slug)
) ENGINE=InnoDB;

-- A revision makes exported/shared install plans reproducible after catalog entries change.
CREATE TABLE experience_pack_revisions (
  id CHAR(36) NOT NULL,
  experience_pack_id CHAR(36) NOT NULL,
  revision_number INT UNSIGNED NOT NULL,
  description TEXT NULL,
  icon_path VARCHAR(1024) NULL,
  setup_notes TEXT NULL,
  target_bedrock_version VARCHAR(100) NULL,
  status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
  published_at DATETIME NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_experience_pack_revisions_number (experience_pack_id, revision_number),
  CONSTRAINT fk_experience_pack_revisions_pack FOREIGN KEY (experience_pack_id) REFERENCES experience_packs (id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE experience_pack_items (
  id CHAR(36) NOT NULL,
  experience_pack_revision_id CHAR(36) NOT NULL,
  addon_release_id CHAR(36) NOT NULL,
  addon_release_source_id CHAR(36) NULL,
  display_order SMALLINT UNSIGNED NOT NULL,
  is_required BOOLEAN NOT NULL DEFAULT TRUE,
  creator_note TEXT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_experience_pack_items_order (experience_pack_revision_id, display_order),
  UNIQUE KEY uq_experience_pack_items_release (experience_pack_revision_id, addon_release_id),
  CONSTRAINT fk_experience_pack_items_revision FOREIGN KEY (experience_pack_revision_id) REFERENCES experience_pack_revisions (id) ON DELETE CASCADE,
  CONSTRAINT fk_experience_pack_items_release FOREIGN KEY (addon_release_id) REFERENCES addon_releases (id) ON DELETE RESTRICT,
  CONSTRAINT fk_experience_pack_items_source_release FOREIGN KEY (addon_release_source_id, addon_release_id) REFERENCES addon_release_sources (id, addon_release_id) ON DELETE RESTRICT
) ENGINE=InnoDB;

CREATE TABLE experience_pack_setup_steps (
  id CHAR(36) NOT NULL,
  experience_pack_revision_id CHAR(36) NOT NULL,
  step_order SMALLINT UNSIGNED NOT NULL,
  title VARCHAR(255) NOT NULL,
  instruction TEXT NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_experience_pack_setup_steps_order (experience_pack_revision_id, step_order),
  CONSTRAINT fk_experience_pack_setup_steps_revision FOREIGN KEY (experience_pack_revision_id) REFERENCES experience_pack_revisions (id) ON DELETE CASCADE
) ENGINE=InnoDB;
