#!/bin/bash

cp -a bin/pigeon docker/pigeon && cp -a external/website docker/website && cp configs/pigeon.yaml docker/pigeon.yaml

docker build -t $1 docker
