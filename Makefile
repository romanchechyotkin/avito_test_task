.PHONY: build
build:
	go build -o build/bin cmd/app/main.go

.PHONY: run
run: build
	build/bin

.PHONY: mock
mock:
	go run go.uber.org/mock/mockgen@latest \
		-source internal/service/service.go \
		-destination internal/service/mocks/mocks.go \
		-package mocks

.PHONY: docs
docs:
	swag init -g ./cmd/app/main.go -o ./docs --parseDependency --parseInternal

.PHONY: gen
gen: mock docs

.PHONY: unit-test
unit-test:
	go test ./... -v -tags=unit

.PHONY: integration-test
integration-test:
	go test ./... -v -tags=integration

.PHONY: coverage
coverage:
	go test -coverprofile=cover.out -covermode=atomic -v -coverpkg=./... -tags=unit,integration ./...

.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html=cover.out

.PHONY: test
test: unit-test integration-test

.PHONY: compose-up
compose-up:
	docker compose up --build