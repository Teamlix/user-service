MIGRATE := ./migrate

UNAME_S := $(shell uname -s)

MIGRATE_VER?=v4.15.2
MIGRATE_PATH:=$(shell dirname $(MIGRATE))

MIGRATE_MONGO_PATH=./db/migrations/mongo

ifeq ($(UNAME_S),Linux)
	OS:=linux-amd64.tar.gz
endif

ifeq ($(UNAME_S),Darwin)
	OS:=darwin-arm64.tar.gz
endif

MONGODB_HOST := $(MONGODB_HOST)
MONGODB_DATABASE := $(MONGODB_DATABASE)
MONGODB_USER := $(MONGODB_USER)
MONGODB_PASSWORD := $(MONGODB_PASSWORD)



CONNECTION_STRING := mongodb://$(MONGODB_USER):$(MONGODB_PASSWORD)@$(MONGODB_HOST)/$(MONGODB_DATABASE)?sslmode=disable

install:
	@mkdir -p $(MIGRATE_PATH)
	curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VER)/migrate.$(OS) | tar -vxz -C $(MIGRATE_PATH) migrate

migrate-up:
	@echo $(CONNECTION_STRING)
	@${MIGRATE} -database $(CONNECTION_STRING) -path ${MIGRATE_MONGO_PATH} up


all: install migrate-up


