NAME ?=
SQL_MIGRATION_DB_CONFIG = migrations/dbconfig.yml
SQL_MIGRATION_ENV = dynamic
SQL_WRITER_DATABASE_USER = backendtest
SQL_WRITER_DATABASE_HOST = 127.0.0.1
SQL_WRITER_DATABASE_PORT = 3306
SQL_WRITER_DATABASE_PASS = backendtest
SQL_WRITER_DATABASE = backendtest
SQL_SEED_PATH = migrations/testdata/seed.sql

.PHONY: migratenew migrateup migratedown

migratenew:
	sql-migrate new -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV) $(NAME)

migrateup:
	sql-migrate up -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV) --dryrun
	sql-migrate up -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV)

migratedown:
	sql-migrate down -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV) --dryrun
	sql-migrate down -config=$(SQL_MIGRATION_DB_CONFIG) -env=$(SQL_MIGRATION_ENV)

cleandatabase:
	mysql \
		--user $(SQL_WRITER_DATABASE_USER) \
		--host $(SQL_WRITER_DATABASE_HOST) \
		-P $(SQL_WRITER_DATABASE_PORT) \
		-p$(SQL_WRITER_DATABASE_PASS) \
		-e 'DROP DATABASE backendtest; CREATE DATABASE backendtest;'

seed:
	make migratedown
	make cleandatabase
	make migrateup
	mysql \
		--user $(SQL_WRITER_DATABASE_USER) \
		--host $(SQL_WRITER_DATABASE_HOST) \
		-P $(SQL_WRITER_DATABASE_PORT) \
		-p$(SQL_WRITER_DATABASE_PASS) \
		$(SQL_WRITER_DATABASE) \
		-e 'source $(SQL_SEED_PATH)'
