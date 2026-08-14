-- Local development API account. Do not use this password outside a local copy.
CREATE USER IF NOT EXISTS 'bedrock_app'@'127.0.0.1' IDENTIFIED BY 'change-me';
GRANT SELECT, INSERT, UPDATE, DELETE ON bedrock_experience_maker.* TO 'bedrock_app'@'127.0.0.1';
FLUSH PRIVILEGES;

-- Local browser login. The password is represented only by its bcrypt hash.
INSERT INTO bedrock_experience_maker.users (id, username, password_hash)
VALUES ('c37c9db0-4a8d-44a1-9ad7-2472c428257b', 'cocoawarrior', '$2a$10$dNDhppL6Qo/hBLd3UKS/U.RMTgo/JtHGojZDVh9TzFizF9udfxM36')
ON DUPLICATE KEY UPDATE password_hash = VALUES(password_hash);

UPDATE bedrock_experience_maker.experience_packs
SET creator_user_id = 'c37c9db0-4a8d-44a1-9ad7-2472c428257b'
WHERE creator_name = 'CocoaWarrior';
