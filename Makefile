include Makefile.ledger
.PHONY: all
all: lint install

.PHONY: install
install: go.sum
		go install $(BUILD_FLAGS) ./cmd/sgn
		go install $(BUILD_FLAGS) ./cmd/sgncli

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

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	@if ! [ -f build/node0/sgn/config/genesis.json ]; then\
		docker run --rm -v $(CURDIR)/build:/sgn:Z celer-network/sgnnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ;\
	fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down


######### utils
.PHONY: prepare-docker-env
prepare-docker-env:
	rm -rf ./docker-test-mount
	cp -r ./test/multi-node-data .
	mv ./multi-node-data ./docker-test-mount
