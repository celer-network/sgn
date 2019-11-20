include Makefile.ledger
all: lint install

install: go.sum
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/sgn
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/sgncli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

copy-test-data:
	cp -r test/data/.sgn ~/.sgn
	cp -r test/data/.sgncli ~/.sgncli

remove-test-data:
	rm -rf ~/.sgn ~/.sgncli

update-test-data: remove-test-data copy-test-data

copy-test-config:
	cp test/data/.sgn/config/genesis.json ~/.sgn/config/genesis.json
	cp test/data/.sgncli/config/config.toml ~/.sgncli/config/config.toml