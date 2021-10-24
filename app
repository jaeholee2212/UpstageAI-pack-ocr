#!/bin/bash

# -u: not allow unbound variables
# -o: pipeline fail if any of commands in the pipeline fail
set -uo pipefail;

cmd=${1}
shift

source ${PWD}/tools/env.sh
load_env

case ${cmd} in
  up)
    ${PWD}/tools/up.sh
    ;;
  dn)
    ${PWD}/tools/dn.sh
    ;;
  re)
    ${PWD}/app dn
    ${PWD}/app up
    ;;
  infra|prod)
    stack_type=${cmd}
    name_or_cmd=${1}
    case ${name_or_cmd} in
      up|dn|re)
        sub_dirs=($(ls -d -- ${PWD}/${stack_type}/*/ | rev | cut -d "/" -f 2 | rev))
        for stack_name in "${sub_dirs[@]}"; do
          stack_cfg=${PWD}/${stack_type}/${stack_name}/docker-stack.yaml
          if [[ -e ${stack_cfg} ]]; then
            ${PWD}/tools/stack.sh ${name_or_cmd} ${stack_name} ${stack_cfg}
          fi
        done
        ;;
      *)
        name=${name_or_cmd}
        stack_cmd=${2}
        if [[ -d "${PWD}/${stack_type}/${name}" ]]; then
          ${PWD}/tools/stack.sh "${stack_cmd}" "${name}" "${PWD}/${stack_type}/${name}/docker-stack.yaml"
        fi
        ;;
    esac
    ;;
  env)
    printenv | grep "^INFRA_"
    ;;
  *)
    echo "[app]: unknown command=${cmd}"
    ;;
esac