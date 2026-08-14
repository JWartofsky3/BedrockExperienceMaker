USE bedrock_experience_maker;

INSERT INTO addons (
  id, display_name, creator_name, description, icon_path, curseforge_url, mcpedl_url,
  current_version, minecraft_version_note, manifest_data
) VALUES
(
  'cocoas-monsters',
  'Cocoa''s Monsters',
  'CocoaWarrior',
  'Adds unique monsters to a Bedrock world, including the Ogre, Yeti, Demon, and Lizord.',
  'https://media.forgecdn.net/attachments/description/1099487/description_3780e350-1695-46d9-a875-d29ccdbdd5ef.png',
  'https://www.curseforge.com/minecraft-bedrock/addons/cocoas-monsters',
  'https://mcpedl.com/cocoas-monsters/',
  '1.1.0',
  'Manifest minimum engine version: 1.21.30 (behavior pack).',
  JSON_OBJECT('packs', JSON_ARRAY('behavior', 'resource'), 'script_modules', JSON_ARRAY('@minecraft/server', '@minecraft/server-ui'))
),
(
  'cocoas-birds',
  'Cocoa''s Birds',
  'CocoaWarrior',
  'Adds birds through paired behavior and resource packs.',
  'https://media.forgecdn.net/attachments/description/1099487/description_888e8584-6f14-4718-961e-91c30ba279a5.png',
  'https://www.curseforge.com/minecraft-bedrock/addons/cocoas-birds',
  'https://mcpedl.com/cocoas-birds/',
  '1.0.0',
  'Manifest minimum engine version: 1.21.30 (behavior pack).',
  JSON_OBJECT('packs', JSON_ARRAY('behavior', 'resource'))
),
(
  'auto-miner',
  'Auto Miner',
  'CocoaWarrior',
  'Adds an automated mining experience through paired behavior and resource packs.',
  'https://r2.mcpedl.com/submissions/247167/images/cocoawarriors-auto-miner_2.png',
  'https://www.curseforge.com/minecraft-bedrock/addons/auto-miner',
  'https://mcpedl.com/auto-miner/',
  '1.0.2',
  'Manifest minimum engine version: 1.21.30. The behavior-pack header reports 1.0.1.',
  JSON_OBJECT('packs', JSON_ARRAY('behavior', 'resource'), 'script_modules', JSON_ARRAY('@minecraft/server'))
),
(
  'cocoas-structures',
  'Cocoa''s Structures',
  'CocoaWarrior',
  'Adds world structures designed to be used with Cocoa''s Monsters.',
  'https://media.forgecdn.net/attachments/description/1099487/description_1844ad0f-323d-43f5-bcd5-225c70632c72.PNG',
  NULL,
  NULL,
  NULL,
  NULL,
  JSON_OBJECT('dependencies', JSON_ARRAY('addons/cocoas-monsters'))
),
(
  'treecapitators',
  'Treecapitators',
  'CocoaWarrior',
  'Adds Treecapitator tools that cut trees faster.',
  'https://media.forgecdn.net/avatars/thumbnails/1195/393/256/256/638770982001723258.png',
  NULL,
  NULL,
  '1.0.1',
  'Manifest minimum engine version: 1.21.30.',
  JSON_OBJECT('packs', JSON_ARRAY('behavior', 'resource'), 'script_modules', JSON_ARRAY('@minecraft/server'))
)
ON DUPLICATE KEY UPDATE
  display_name = VALUES(display_name),
  creator_name = VALUES(creator_name),
  description = VALUES(description),
  icon_path = VALUES(icon_path),
  curseforge_url = VALUES(curseforge_url),
  mcpedl_url = VALUES(mcpedl_url),
  current_version = VALUES(current_version),
  minecraft_version_note = VALUES(minecraft_version_note),
  manifest_data = VALUES(manifest_data);

INSERT INTO addon_dependencies (addon_id, dependency_id)
VALUES ('cocoas-structures', 'cocoas-monsters')
ON DUPLICATE KEY UPDATE dependency_id = VALUES(dependency_id);

INSERT INTO experience_packs (id, slug, name, creator_name, description, setup_notes)
VALUES (
  'b2f3a0ca-8b33-4e01-935e-1d0dd410e8d9',
  'cocoas-world-adventure',
  'Cocoa''s World Adventure',
  'CocoaWarrior',
  'A world-focused experience with Cocoa''s creatures, structures, and birds.',
  'Install Cocoa''s Monsters before Cocoa''s Structures, then activate each selected behavior and resource pack in the world settings.'
)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  creator_name = VALUES(creator_name),
  description = VALUES(description),
  setup_notes = VALUES(setup_notes);

-- Reset only the example pack's entries so its fixed seed order can be
-- re-applied without colliding with an order changed during local testing.
DELETE FROM experience_pack_addons
WHERE experience_pack_id = 'b2f3a0ca-8b33-4e01-935e-1d0dd410e8d9';

INSERT INTO experience_pack_addons (id, experience_pack_id, addon_id, install_order)
VALUES
  ('ecef74ca-399a-40f1-8eef-9e03d3a5c0d1', 'b2f3a0ca-8b33-4e01-935e-1d0dd410e8d9', 'cocoas-monsters', 1),
  ('64e59d6b-fb5e-4c22-9a41-2d61c0d65291', 'b2f3a0ca-8b33-4e01-935e-1d0dd410e8d9', 'cocoas-structures', 2),
  ('be0ba8a5-07ac-4e87-8a6a-ac6a623b6161', 'b2f3a0ca-8b33-4e01-935e-1d0dd410e8d9', 'cocoas-birds', 3)
ON DUPLICATE KEY UPDATE install_order = VALUES(install_order);
