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
MYSQL_CONTAINER := mysql-db
MYSQL_USER_ROOT := root
MYSQL_USER_NAME := web
MYSQL_USER_PASSWORD := 123456

.PHONY: init_mysql
init_mysql:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_ROOT) -p$(MYSQL_USER_PASSWORD) < ./internal/sql/init_db.sql

.PHONY: init_sessions_table
init_sessions_table:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_ROOT) -p$(MYSQL_USER_PASSWORD) < ./internal/sql/init_session_table.sql


.PHONY: init_users_table
init_users_table:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_ROOT) -p$(MYSQL_USER_PASSWORD) < ./internal/sql/init_users_table.sql

.PHONY: remove_db
remove_db:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_ROOT) -p$(MYSQL_USER_PASSWORD) < ./internal/sql/remove_db.sql

.PHONY: init_db
init_db: init_mysql init_sessions_table

.PHONY: clear_db
	init_db: remove_db

.PHONY: check_mysql
check_mysql:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_NAME) -p$(MYSQL_USER_PASSWORD) < ./internal/sql/get_snippets.sql


.PHONY: start_db
start_db:
	docker run --name $(MYSQL_CONTAINER) -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql:latest
	

.PHONY: open_db
open_db:
	docker exec -it $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER_NAME) -p$(MYSQL_USER_PASSWORD) $(MYSQL_DATABASE)

## App
APP_HOST := 192.168.1.12:8080
.PHONY: insert_test
insert_test:
	curl -iL -X POST $(APP_HOST)/snip/create

.PHONY: get_test
get_test:
	curl -iL -X GET $(APP_HOST)/snip/view?id=1
