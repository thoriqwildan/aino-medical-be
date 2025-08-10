.PHONY: scan@library lint gen@app gen@users gen@gorm compose@up compose@down \
        migrate@create migrate@up migrate@down migrate@version migrate@clean help

# ---------- Configuration ----------
APP_PATH := ./internal/app
USER_PATH := ./internal/features/users
MIGRATE_PATH := db/migrations
DB_CONNECTION := mysql://username:password@tcp(host:port)/dbname?param1=value1&param2=value2

# ---------- Security & Lint ----------
scan@library: ## Run OSV scanner on dependencies
	osv-scanner scan source go.mod

lint: ## Run Go linter
	golangci-lint run ./...

# --------------- Running Golang Command ---------------

run@app: ## Running Golang App
	go run cmd/web/main.go

run@seed: ## Running Seeders
	go run cmd/seed/main.go

run@build.prod: ## Running Build Application to Binary File ( Production )
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/main cmd/web/main.go

run@build.dev: ## Running Build Application to Binary File ( Development )
	go build -o build/main cmd/web/main.go

# ---------- Generation Any Tools or Documentation ----------
gen@doc: ## Generate wire for app
	swag init -g cmd/web/main.go

# ----------- Docker -----------------

ct@build.prod:
	@bash -c ' \
	read -p "Enter image name: " name; \
	read -p "Enter unique id app: " uid; \
	read -p "Enter username app: " username; \
	docker build . -t $$name:latest --build-arg username=$$username --build-arg uid=$$uid \
	'
ct@build.dev:
	@bash -c ' \
	read -p "Enter image name: " name; \
	read -p "Enter unique id app: " uid; \
	read -p "Enter username app: " username; \
	docker build . -f dev.Dockerfile -t $$name:dev --build-arg username=$$username --build-arg uid=$$uid \
	'

# ---------- Docker Compose ----------
compose@up: ## Start Docker Compose in detached mode
	docker compose up -d

compose@down: ## Stop and remove containers/volumes (extendable)
	docker compose down $(ARGS)

compose@all: ## See all container by scoped docker-compose
	docker compose ps -a

# ---------- Database Migration ----------
migrate@create: ## Create a new migration file
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./$(MIGRATE_PATH) -seq $$name

migrate@up: ## Apply all up migrations
	migrate -source file://$(MIGRATE_PATH) -database "$(DB_CONNECTION)" up $(ARGS)

migrate@down: ## Revert all applied migrations
	migrate -source file://$(MIGRATE_PATH) -database "$(DB_CONNECTION)" down $(ARGS)

migrate@version: ## Show current migration version
	@echo "Checking current version..."
	@migrate -source file://$(MIGRATE_PATH) -database "$(DB_CONNECTION)" version

migrate@clean: ## Clean the dirty migrate
	@read -p "Enter migration version: " version; \
	migrate -source file://$(MIGRATE_PATH) -database "$(DB_CONNECTION)" force $$version


# ---------- Help ----------
help: ## Show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z0-9_@-]+:.*##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
