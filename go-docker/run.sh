#!/usr/bin/env bash

mkdir -p ~/logs/go-docker

docker build -t go-docker .

docker run -p 8080:8080 -v ~/logs/go-docker:/app/logs go-docker