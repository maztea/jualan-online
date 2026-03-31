# Variables
DOCKER_COMPOSE = docker-compose
APP_NAME = jualan-online
SERVICES = auth-service

.PHONY: all up down restart build logs ps deploy migrate seed help

all: help

## 🚀 Development
up: ## Start all services in background
	$(DOCKER_COMPOSE) up -d

down: ## Stop and remove all containers
	$(DOCKER_COMPOSE) down

restart: down up ## Restart all services

build: ## Build or rebuild services
	$(DOCKER_COMPOSE) build

logs: ## Follow logs of all services
	$(DOCKER_COMPOSE) logs -f

ps: ## List all running containers
	$(DOCKER_COMPOSE) ps

## 🛠️ Database & Maintenance
migrate: ## Run database migrations (Note: Auth Service uses GORM AutoMigrate on startup)
	@echo "Running migrations via service startup..."
	$(DOCKER_COMPOSE) restart auth-service

seed: ## Seed initial data to database
	@echo "Seeding initial data for auth-service..."
	$(DOCKER_COMPOSE) exec auth-service ./main -seed

## 🚢 Deployment (Docker Swarm)
deploy: ## Deploy the stack to Docker Swarm
	@echo "Deploying stack to Docker Swarm..."
	docker stack deploy -c docker-compose.yml $(APP_NAME)

remove-stack: ## Remove the stack from Docker Swarm
	docker stack rm $(APP_NAME)

## ❓ Help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
