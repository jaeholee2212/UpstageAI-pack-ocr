#!/bin/bash

# -u: not allow unbound variables
# -o: pipeline fail if any of commands in the pipeline fail
set -uo pipefail;

cmd=${1}
shift

source ${PWD}/tools/env.sh
load_env

case ${cmd} in
  env)
    printenv | grep "^INFRA_"
    ;;
  default)
    echo "[app]: unknown command=${cmd}"
    ;;
esac