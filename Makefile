# Project name 

PROJECT_NAME := SSR_snippet

GO_WEB_FILES := $(wildcard ./cmd/web/*.go)

BUILD_TARGET := $(PROJECT_NAME)

.PHONY: build
build: 
	go build -o $(PROJECT_NAME) $(GO_WEB_FILES)

.PHONY: run
run: build 
	./$(BUILD_TARGET)

## Docker 

DOCKER_CONTAINER_NAME := mysql-container
DDL_SQL_FILE := ./internal/sql/ddl.sql
DML_SQL_FILE := ./internal/sql/dml.sql

## Mysql

MYSQL_DATABASE := snippetbox
MYSQL_CONTAINER := mysql-container
MYSQL_USER_ROOT := root
MYSQL_USER_NAME := web
MYSQL_USER_PASSWORD := 123456

.PHONY: init_mysql
init_mysql:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_ROOT) -p$(MYSQL_USER_PASSWORD) < ./internal/sql/init_db.sql

.PHONY: check_mysql
check_mysql:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_NAME) -p$(MYSQL_USER_PASSWORD) < ./internal/sql/get_snippets.sql

.PHONY: open_db
open_db:
	docker exec -it $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_NAME) -p$(MYSQL_USER_PASSWORD) $(MYSQL_DATABASE)
