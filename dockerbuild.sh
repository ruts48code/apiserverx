#!/usr/bin/bash

docker build --no-cache -t registry.rmutsv.app/authapon/apiserver:x .
docker push registry.rmutsv.app/authapon/apiserver:x