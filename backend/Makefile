MIGRATION_DIR="file://data/schema"

.Phony: all
all: install run

.Phony: install
install: gen_key gen_rpc gen_orm gen_mock

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

.Phony: @db_refresh
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
	@go test ./domain/...
	@go test ./app/...
	@go test ./util/...
	@echo "done unittest!!"

.Phony: integration
integration: depend
	@go test ./test/integration/...
	@echo "done integration test"

.Phony: all_test
all_test: install test integration
	@echo "done all tests!!"
