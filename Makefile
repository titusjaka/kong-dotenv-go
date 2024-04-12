.PHONY: lint test

RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

export CGO_ENABLED=0

lint:
	@golangci-lint run

test:
	@go test ./...
