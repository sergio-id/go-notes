#!/bin/bash

PREFIX=simple025/go-notes

# building Docker images
for f in $(find . -name 'Dockerfile')
do
  echo "BUILD $f"
  BASE=$(basename "$(dirname "${f}")")
  docker build -f "${f}" ../ --tag ${PREFIX}-"${BASE}":1.0.0
done

docker login

# push Docker images
for f in $(find . -name 'Dockerfile')
do
  echo "PUSH $f"
  BASE=$(basename "$(dirname "${f}")")
  docker push ${PREFIX}-"${BASE}":1.0.0
done
