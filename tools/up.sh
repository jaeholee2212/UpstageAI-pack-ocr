#!/bin/bash

is_in_swarm_mode() {
  local status=$(docker node ls --filter "role=manager" --format {{.Status}})
  if [[ ${status} = "Ready" ]]; then
    return 0
  fi

  return 1
}

ensure_network() {
  local name=$1
  shift
  local filtered=$(docker network ls --filter name=$name --format {{.Name}})
  if [[ $name = $filtered ]]; then
    echo "[up]: previous network destroyed - $(docker network rm $name)"
  fi
  echo "[up]: new network - $(docker network create $@ $name)"
}

node_add_label() {
  nodeid=$(docker info -f '{{.Swarm.NodeID}}')
  docker node update --label-add $@ ${nodeid}
  echo "[up]: node labels: $(docker node inspect --format '{{.Spec.Labels}}' ${nodeid})"
}


if is_in_swarm_mode; then 
  echo "[up]: We're already in the swarm mode now. no further actions."
  exit 0
fi

echo "[up]: swarm advertise address - ${INFRA_SWARM_ADVERTISE_ADDR}"
docker swarm init --advertise-addr=${INFRA_SWARM_ADVERTISE_ADDR}

echo "[up]: creating mesh networks"
ensure_network 'infra-public' --driver=overlay --attachable
ensure_network 'infra-private' --driver=overlay --attachable --internal

echo "[up]: node configurations"
node_add_label db=true
node_add_label ocr.pack=true

echo "[up]: starting infra"
${PWD}/app infra up
echo "[up]: starting prod"
${PWD}/app prod up

echo "[up]: ok"

