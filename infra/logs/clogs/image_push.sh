#! /bin/bash

arch=${1:-arm64}
tag="popopome/clogs:${arch}"

docker build -t ${tag} .  &&  \
  docker push ${tag}

