# Makefile for kmaster development

# prepare: Install pre-commit, commitlint, and tools from .tools-versions
prepare:
	@echo "Install pre-commit via Homebrew..."
	brew install pre-commit

	@echo "Initializing environment for git..."
	pre-commit install -t pre-commit
	pre-commit install -t pre-push
	pre-commit install -t commit-msg

	@echo "Enable go module..."
	@echo "Since you are actively in GOPATH"
	@export GO111MODULE=on

	@echo "Installing commit's message linter - commitlint..."
	brew install yarn --ignore-dependencies
	yarn global add commitlint

	@echo "Install tools from .tools-versions"
	asdf install

# rehooks for pre-commit hooks
rehooks:
	@echo "Clean git hooks cache..."
	pre-commit clean

	@echo "Initializing environment for git..."
	pre-commit install -t pre-commit
	pre-commit install -t pre-push
	pre-commit install -t commit-msg

# Run development environment with docker-compose-dev.yml (local development with Air)
dev.local.up:
	docker-compose -f docker-compose-dev.yml up --build

# Run development environment with standard docker-compose.yml
dev.up:
	docker-compose -f docker-compose.yml up --build

# Stop development containers for local setup
dev.local.down:
	docker-compose -f docker-compose-dev.yml down

# Stop development containers for standard setup
dev.down:
	docker-compose -f docker-compose.yml down