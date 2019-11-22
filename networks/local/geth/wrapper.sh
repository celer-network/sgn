#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/geth/${BINARY:-geth}
LOG=${LOG:-geth.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'geth' E.g.: -e BINARY=geth_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export GETHHOME="/geth/geth-env"
export GETHDATA="/geth/geth-env/geth-data"

mkdir -p "${GETHDATA}"

"${BINARY}" --datadir "${GETHDATA}" init "${GETHHOME}"/mainchain_genesis.json

if [ -d "$(dirname "${GETHHOME}"/"${LOG}")" ]; then
  "${BINARY}" "$@" | tee "${GETHHOME}/${LOG}"
else
  "${BINARY}" "$@"
fi
