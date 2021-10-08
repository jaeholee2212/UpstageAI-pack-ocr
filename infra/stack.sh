#!/bin/bash

stack_name=$1
stack_cfg=$2
command=$3

stop() {
  docker stack rm ${stack_name}
}

start() {
  # Sadly, a command 'docker stack' does not respect .env files
  # like `docker compose`. The following line is to load envs
  # from a `.env` file
  export $(cat .env) > /dev/null 2>&1; 
  docker stack deploy -c ${stack_cfg} ${stack_name}
}

loadenvs() {
  local envfile="${1:-.env}"
  set -a && . $envfile && set +a
}

ensure_no_network() {
  local name=$1
  shift
  local tries=0
  while [[ $tries -lt 10 ]]
  do
    local filtered=$(docker network ls --filter name=$name --format {{.Name}})
    if [[ -z $filtered ]]; then
      break
    fi
    echo "network: ${name} exists. trying again. (${tries})"
    sleep 1
    ((tries++))
  done
  return 0
}

case $command in
  up)
    start
    echo "${stack_name} is now up"
    ;;
  down)
    stop
    echo "${stack_name} is down"
    ;;
  re)
    stop
    ensure_no_network "${1}_default"
    start
    ;;
  *)
    echo "cmd [up|down|re]"
esac
