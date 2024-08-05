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
	echo docs

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
	go test -coverprofile=coverage.out ./... -tags=unit,integration inte

.PHONY: test
test: unit-test integration-test
