MIGRATION_DIR="file://data/schema"

all: depend db_refresh gen_rpc gen_orm gen_key run

.Phony: gen_key
gen_key:
	@echo "PRIVATE_KEY_PATH=$(PRIVATE_KEY_PATH)"
	@mkdir -p $(KEY_DIR)
	@openssl genrsa -out $(PRIVATE_KEY_PATH) 2048
	@chmod 700 $(KEY_DIR)
	@chmod 600 $(PRIVATE_KEY_PATH)
	@echo "done generate key!!"

.Phony: fmt
fmt:
	@go fmt ./...

.Phony: depend
depend:
	@go mod tidy

.Phony: run
run: depend gen_key
	@air -c .air.toml
	@echo "done"

.Phony: gen_rpc
gen_rpc: 
	@buf lint
	@buf format -w
	@buf generate
	@go mod tidy
	@echo "done gen_rpc!!"

.Phony: db_refresh
db_refresh:
	@migrate -database=$(DATABASE_URL) -source=$(MIGRATION_DIR) drop -f
	@migrate -database=$(DATABASE_URL) -source=$(MIGRATION_DIR) up
	@echo "done refresh database!!"

.Phony: gen_orm
gen_orm: db_refresh
	@sqlc generate
	@go mod tidy
	@echo "done gen_orm!!"

.Phony: gen_mock
gen_mock:
	@mkdir -p ./test/mocks
	@mockery
	@go mod tidy
	@echo "done generate mocks!!"

.Phony: test
test: depend
	@go test $(REPO_NAME)/domain/object/entity
	@go test $(REPO_NAME)/domain/service
	@go test $(REPO_NAME)/app/usecase
	@go test $(REPO_NAME)/app/util/auth
	@go test $(REPO_NAME)/app/util/contextkey
	@go test $(REPO_NAME)/app/handler
	@echo "done unittest!!"

.Phony: test_all
test_all: depend gen_mock
	@go test -v ./...
	@echo "done all tests!!"
