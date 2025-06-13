include .env
export

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh trainer internal/trainer/ports ports

.PHONY: test
test:
	@./scripts/test.sh trainer .test.env