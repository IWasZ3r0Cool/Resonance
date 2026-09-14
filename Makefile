.PHONY: bootstrap dev-api dev-web generate-api test build check fmt

bootstrap:
	pnpm --dir web install --frozen-lockfile
	pnpm --dir web generate:api

dev-api:
	go run ./cmd/api

dev-web:
	pnpm --dir web dev

generate-api:
	pnpm --dir web generate:api

test:
	go test ./...
	go test -tags=integration ./tests/integration
	pnpm --dir web typecheck
	pnpm --dir web lint
	pnpm --dir web test

build:
	mkdir -p bin
	go build -o bin/resonance-api ./cmd/api
	pnpm --dir web build

fmt:
	gofmt -w $$(find api cmd internal -name '*.go')

check: fmt generate-api test build
	git diff --check
