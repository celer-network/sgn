# Makefile for the "sgnnode" docker image.

.PHONY: all
all: sgn geth

.PHONY: sgn
sgn:
	DOCKER_BUILDKIT=1 docker build --tag celer-network/sgnnode sgnnode

.PHONY: geth
geth:
	DOCKER_BUILDKIT=1 docker build --tag celer-network/geth geth