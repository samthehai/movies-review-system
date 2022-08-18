
-- +migrate Up
CREATE FULLTEXT INDEX `fulltext_movies` ON `movies` (`original_title`, `overview`, `original_language`);

-- +migrate Down
ALTER TABLE `movies` DROP INDEX `fulltext_movies`;
