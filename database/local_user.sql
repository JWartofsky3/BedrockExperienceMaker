-- Local development API account. Do not use this password outside a local copy.
CREATE USER IF NOT EXISTS 'bedrock_app'@'127.0.0.1' IDENTIFIED BY 'change-me';
GRANT SELECT, INSERT, UPDATE, DELETE ON bedrock_experience_maker.* TO 'bedrock_app'@'127.0.0.1';
FLUSH PRIVILEGES;
