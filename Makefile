DC := docker compose
PROJECT_NAME := ozon

.PHONY: up-postgres up-inmemory down 

up-postgres:
	$(DC) --project-name $(PROJECT_NAME) --profile postgres up -d

up-inmemory:
	$(DC) --project-name $(PROJECT_NAME) up -d

down:
	$(DC) --project-name $(PROJECT_NAME) down -v