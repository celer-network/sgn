include Makefile.ledger
all: lint install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/sgn
		go install $(BUILD_FLAGS) ./cmd/sgncli

install-test: go.sum
	go install $(BUILD_FLAGS) ./cmd/sgntest

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

e2e-test:
	go test -failfast -v -timeout 15m github.com/celer-network/sgn/test/e2e

copy-test-data:
	cp -r test/data/.sgn ~/.sgn
	cp -r test/data/.sgncli ~/.sgncli

remove-test-data:
	rm -rf ~/.sgn ~/.sgncli

update-test-data: remove-test-data copy-test-data

copy-test-config:
	cp test/data/.sgn/config/genesis.json ~/.sgn/config/genesis.json
	cp test/data/.sgncli/config/config.toml ~/.sgncli/config/config.toml