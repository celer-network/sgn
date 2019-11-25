include Makefile.ledger
.PHONY: all
all: lint install

.PHONY: install
install: go.sum
		go install -mod=readonly ./cmd/sgn
		go install -mod=readonly ./cmd/sgncli

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
	go build -mod=readonly -o build/sgn ./cmd/sgn
	go build -mod=readonly -o build/sgncli ./cmd/sgncli

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
.PHONY: prepare-geth-env
prepare-geth-env:
	rm -rf ./build/geth-env
	cp -r ./testing/env ./build/
	mv ./build/env ./build/geth-env

.PHONY: prepare-keys
prepare-keys:
	rm -rf ./build/keys
	cp -r ./test/keys ./build/

.PHONY: prepare-sgn-data
prepare-sgn-data:
	rm -rf ./build/node*
	cp -r ./test/data ./build/
	mv ./build/data/.sgn ./build/data/sgn
	mv ./build/data/.sgncli ./build/data/sgncli
	mv ./build/data ./build/node0
	cp -r ./build/node0 ./build/node1
	cp -r ./build/node0 ./build/node2
	cp -r ./build/node0 ./build/node3

.PHONY: prepare-config
prepare-config:
	cp config-docker.json ./build/config.json

.PHONY: prepare-docker-env
prepare-docker-env: prepare-geth-env prepare-keys prepare-sgn-data prepare-config
