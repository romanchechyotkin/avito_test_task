.PHONY: build
build:
	go build -o build/bin cmd/app/main.go

.PHONY: run
run: build
	build/bin

.PHONY: mock
mock:
	go generate internal/service/service.go

.PHONY: docs
docs:
	echo docs

.PHONY: gen
gen: mock docs