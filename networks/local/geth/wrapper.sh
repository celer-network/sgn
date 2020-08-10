#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/geth/bin/${BINARY:-geth}
LOG=${LOG:-geth.log}

##
## Run binary with all parameters
##
export GETHHOME="/geth/env"
export GETHDATA="/geth/env/data"

mkdir -p "${GETHDATA}"

"${BINARY}" --datadir "${GETHDATA}" init "${GETHHOME}"/mainchain_genesis.json

if [ -d "$(dirname "${GETHHOME}"/"${LOG}")" ]; then
	"${BINARY}" --datadir "${GETHDATA}" "$@" 2>&1 | tee "${GETHHOME}/${LOG}"
else
	"${BINARY}" --datadir "${GETHDATA}" "$@"
fi
