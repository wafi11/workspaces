run-operator:
	@echo "Running operator locally..."
	@cd operator && air

run-operator-with-logs:
	@echo "Running operator locally with logs..."
	@go run operator/cmd/main.go --log-level=debug

build-operator:
	@echo "Building operator binary..."
	@go build -o bin/operator operator/cmd/main.go

run-web:
	@echo "Running web locally..."
	@cd web && npm run dev

run-backend:
	@echo "Running backend locally..."
	@cd backend && air