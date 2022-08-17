NAME ?=
SQL_MIGRATION_DB_CONFIG = migrations/dbconfig.yml
SQL_MIGRATION_ENV = dynamic

.PHONY: migratenew migrateup migratedown

migratenew:
	sql-migrate new -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV) $(NAME)

migrateup:
	sql-migrate up -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV) --dryrun
	sql-migrate up -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV)

migratedown:
	sql-migrate down -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV) --dryrun
	sql-migrate down -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV)
