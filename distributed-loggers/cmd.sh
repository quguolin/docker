#!/usr/bin/env bash

#node
docker-compose up --scale distributed-logger=3


docker kill (leader image id)