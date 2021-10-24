#!/bin/bash
if docker swarm leave --force; then
  echo "[dn]: done"
fi