.PHONY: create status version up down

include ../.env

DATABASE_DRIVER := postgres

name := new_migration

# Usage: `make create name=my-cool-migration`
create:
	@goose create $(name) sql && goose fix && git add .

status:
	@goose $(DATABASE_DRIVER) $(DATABASE_CONNECTION_STRING) status

version:
	@goose $(DATABASE_DRIVER) $(DATABASE_CONNECTION_STRING) version

up:
	@goose $(DATABASE_DRIVER) $(DATABASE_CONNECTION_STRING) up

down:
	@goose $(DATABASE_DRIVER) $(DATABASE_CONNECTION_STRING) down
