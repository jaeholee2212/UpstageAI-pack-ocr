#!/bin/bash

is_apple_silicon() {
  chip=$(/Volumes/Macintosh\ HD/usr/sbin/system_profiler SPHardwareDataType | grep "Apple M1" | sed -e 's/^ *//g' | cut -d " " -f 2,3)
  if [[ ! $? ]]; then return 1; fi
  if [[ ${chip} = "Apple M1" ]]; then
    return 0
  fi
  return 1
}

load_env() {
  # Load `.env` into the current run context since `docker stack` does not
  # respect dot-files unfortunately.
  export $(cat .env.base) > /dev/null 2>&1;
  export $(cat .env) > /dev/null 2>&1;
  if is_apple_silicon; then
    echo "[env]: load .env.arm64"
    export $(cat .env.arm64) > /dev/null 2>&1;
  else
    echo "[env]: load .env.amd64"
    export $(cat .env.amd64) > /dev/null 2>&1;
  fi
}
