run-operator:
	@echo "Running operator locally..."
	@cd operator && air

build-wscli:
	docker build

run-operator-with-logs:
	@echo "Running operator locally with logs..."
	@go run operator/cmd/main.go --log-level=debug

build-operator:
	@echo "Building operator binary..."
	@go build -o bin/operator operator/cmd/main.go

run-web:
	@echo "Running web locally..."
	@cd web-app/dashboard && npm run dev

run-backend:
	@echo "Running backend locally..."
	@cd backend && air

remote-github:
	@echo "Remote to GitHub..."
	@git remote set-url origin  git@github.com:wafi11/workspaces.git

remote-gitlab:
	@echo "Remote to GitLab..."
	@git remote set-url origin http://192.168.1.31/root/workspaces.git

gen-proto:
	@echo "Generating Go code from proto..."
	protoc --proto_path=proto/workspace \
	       --go_out=proto/workspace \
	       --go_opt=paths=source_relative \
	       proto/workspace/workspace.proto