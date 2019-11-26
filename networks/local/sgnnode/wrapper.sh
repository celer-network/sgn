#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/sgn/bin/${BINARY:-sgn}
LOG=${LOG:-sgn.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'sgn' E.g.: -e BINARY=sgn_my_test_version"
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
export SGNHOME="/sgn/env/sgn"

if [ -d "$(dirname "${SGNHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${SGNHOME}" "$@" | tee "${SGNHOME}/${LOG}"
else
  "${BINARY}" --home "${SGNHOME}" "$@"
fi
