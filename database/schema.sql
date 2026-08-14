-- Bedrock Experience Pack Builder - simple local MySQL 8.0+ schema
-- IDs are stable, application-owned resource identifiers. Creator manifest UUIDs
-- are optional metadata only and are never used as primary/unique identifiers.

CREATE DATABASE IF NOT EXISTS bedrock_experience_maker
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE bedrock_experience_maker;

-- Manually curated catalog entries. Version and manifest fields are optional
-- because provider pages often expose only the current downloadable release.
CREATE TABLE IF NOT EXISTS addons (
  id VARCHAR(63) NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  creator_name VARCHAR(160) NULL,
  description TEXT NULL,
  icon_path VARCHAR(1024) NULL,
  curseforge_url VARCHAR(2048) NULL,
  mcpedl_url VARCHAR(2048) NULL,
  current_version VARCHAR(100) NULL,
  minecraft_version_note VARCHAR(255) NULL,
  manifest_data JSON NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY ix_addons_display_name (display_name)
) ENGINE=InnoDB;

-- Required catalog relationships. An add-on can depend on more than one
-- cataloged add-on, but cannot depend on itself.
CREATE TABLE IF NOT EXISTS addon_dependencies (
  addon_id VARCHAR(63) NOT NULL,
  dependency_id VARCHAR(63) NOT NULL,
  PRIMARY KEY (addon_id, dependency_id),
  CONSTRAINT chk_addon_dependencies_not_self CHECK (addon_id <> dependency_id),
  CONSTRAINT fk_addon_dependencies_addon FOREIGN KEY (addon_id) REFERENCES addons (id) ON DELETE CASCADE,
  CONSTRAINT fk_addon_dependencies_dependency FOREIGN KEY (dependency_id) REFERENCES addons (id) ON DELETE RESTRICT
) ENGINE=InnoDB;

-- A creator's collection. There are no revisions, publishing states, or user
-- accounts in this first iteration.
CREATE TABLE IF NOT EXISTS experience_packs (
  id CHAR(36) NOT NULL,
  slug VARCHAR(160) NOT NULL,
  name VARCHAR(255) NOT NULL,
  creator_name VARCHAR(160) NULL,
  description TEXT NULL,
  setup_notes TEXT NULL,
  icon_path VARCHAR(1024) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_experience_packs_slug (slug)
) ENGINE=InnoDB;

SET @creator_name_exists := (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = 'bedrock_experience_maker'
    AND table_name = 'experience_packs'
    AND column_name = 'creator_name'
);
SET @creator_name_migration := IF(
  @creator_name_exists = 0,
  'ALTER TABLE experience_packs ADD COLUMN creator_name VARCHAR(160) NULL AFTER name',
  'SELECT 1'
);
PREPARE creator_name_statement FROM @creator_name_migration;
EXECUTE creator_name_statement;
DEALLOCATE PREPARE creator_name_statement;

-- `install_order` is explicit: 1 is the first add-on shown in the installation
-- checklist. It is not a compatibility rule and does not block a selection.
CREATE TABLE IF NOT EXISTS experience_pack_addons (
  id CHAR(36) NOT NULL,
  experience_pack_id CHAR(36) NOT NULL,
  addon_id VARCHAR(63) NOT NULL,
  install_order SMALLINT UNSIGNED NOT NULL,
  selected_source ENUM('curseforge', 'mcpedl', 'other') NULL,
  selected_version VARCHAR(100) NULL,
  note TEXT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_pack_addons_order (experience_pack_id, install_order),
  UNIQUE KEY uq_pack_addons_addon (experience_pack_id, addon_id),
  CONSTRAINT fk_pack_addons_pack FOREIGN KEY (experience_pack_id) REFERENCES experience_packs (id) ON DELETE CASCADE,
  CONSTRAINT fk_pack_addons_addon FOREIGN KEY (addon_id) REFERENCES addons (id) ON DELETE RESTRICT
) ENGINE=InnoDB;
