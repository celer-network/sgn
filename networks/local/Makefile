# Makefile for the "sgnnode" docker image.

.PHONY: all
all: sgn geth

.PHONY: sgn
sgn:
	docker build --tag celer-network/sgnnode sgnnode

.PHONY: geth
geth:
	docker build --tag celer-network/geth geth