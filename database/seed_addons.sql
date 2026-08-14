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
