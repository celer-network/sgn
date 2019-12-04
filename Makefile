include Makefile.ledger
.PHONY: all
all: lint install

.PHONY: install
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

.PHONY: update-test-data
update-test-data: remove-test-data copy-test-data

copy-test-config:
	cp test/data/.sgn/config/genesis.json ~/.sgn/config/genesis.json
	cp test/data/.sgncli/config/config.toml ~/.sgncli/config/config.toml

################################ Docker related ################################
.PHONY: build
build: go.sum
	go build -o build/sgn ./cmd/sgn
	go build -o build/sgncli ./cmd/sgncli

.PHONY: build-linux
build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

GETH_VER = geth-linux-amd64-1.9.1-b7b2f60f
.PHONY: get-geth
get-geth:
	curl -sL https://gethstore.blob.core.windows.net/builds/$(GETH_VER).tar.gz | tar -xz --strip 1 $(GETH_VER)/geth && mv geth ./build;

.PHONY: build-dockers
build-dockers:
	$(MAKE) -C networks/local

# Prepare docker environment for multinode testing
.PHONY: prepare-docker-env
prepare-docker-env: build-dockers get-geth build-linux

# Run geth
.PHONY: localnet-start-geth
localnet-start-geth:
	docker-compose stop geth
	docker-compose up -d geth

# Run a 3-node sgn testnet locally
.PHONY: localnet-start-nodes
localnet-start-nodes: localnet-stop-nodes
	docker-compose up -d sgnnode0 sgnnode1 sgnnode2

# Stop sgn testnet
.PHONY: localnet-stop-nodes
localnet-stop-nodes:
	docker-compose stop sgnnode0 sgnnode1 sgnnode2

# Stop both geth and sgn testnet
.PHONY: localnet-down
localnet-down:
	docker-compose down

# Prepare sgn nodes' data
.PHONY: prepare-sgn-data
prepare-sgn-data:
	rm -rf ./docker-volumes
	cp -r ./test/multi-node-data .
	mv ./multi-node-data ./docker-volumes
