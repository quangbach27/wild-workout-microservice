include .env
export

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh trainer internal/trainer/ports ports
	@./scripts/openapi-http.sh training internal/training/ports ports

.PHONY: proto
proto:
	@./scripts/proto.sh trainer

.PHONY: test
test:
	@./scripts/test.sh trainer

.PHONY: migrate_create
migrate_create:
	migrate create -ext sql -dir sql/migrations/ -seq $(file)

.PHONY: migrate
migrate:
	@./scripts/migrate_sql.sh $(cmd)