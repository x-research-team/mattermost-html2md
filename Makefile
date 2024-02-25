infrastructure.local.up:
	@echo "Running infrastructure local..."
	@docker compose -f deploy/docker/local/docker-compose.yml up --build

infrastructure.local.down:
	@echo "Shutting down infrastructure local..."
	@docker compose -f deploy/docker/local/docker-compose.yml down

tests.run:
	@echo "Running tests..."
	@cd tests
	@go test -v -coverprofile=coverage.txt -covermode atomic -timeout 30s -run ^TestMain$ mattermost-html2md/tests
