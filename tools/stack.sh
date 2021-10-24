#!/bin/bash

command=$1
stack_name=$2
stack_cfg=$3

ensure_docker_image_downloaded() {
  local name=$1
  local match=$(echo $name | grep -E "\\\$\{.+\}$")
  if [[ ${match} = ${name} ]]; then
    name=$(eval "echo \"${name}\"")
  fi

  local filtered=$(docker image ls --format "{{.Repository}}:{{.Tag}}" | grep ${name})
  if [[ $name = $filtered ]]; then
    return 0
  fi
  echo "[image]: downloading - ${name}"
  docker image pull ${name}
  if [[ $? ]]; then
    echo "[image]: downloaded - ${name}"
    return 0
  fi
  echo "[image]: failed to download (${name})"
  exit 1
}

download_images_in_yaml() {
  yaml=$1
  image_names=$(cat ${yaml} | grep image: | sed -e 's/^ *//g' | cut -d " " -f 2)
  for img in ${image_names[@]}; do
    ensure_docker_image_downloaded ${img}
  done
}

stop() {
  docker stack rm ${stack_name}
}

start() {
  download_images_in_yaml ${stack_cfg}
  docker stack deploy -c ${stack_cfg} ${stack_name}
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
    echo "[stack]: starting ${stack_name}... config=${stack_cfg}"
    start
    echo "[stack]: ${stack_name} is now up"
    ;;
  dn|down)
    stop
    echo "[stack]: ${stack_name} is down"
    ;;
  re)
    echo "[stack]: ${stack_name} is being stopped"
    stop
    ensure_no_network "${1}_default"
    echo "[stack]: ${stack_name} is starting"
    start
    ;;
  *)
    echo "cmd [up|down|re]"
esac
