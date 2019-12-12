#!/bin/bash
clear

usage="build.sh TAG"

if [ -z "$1" ]; then
  echo $usage
  exit 1
fi;

git pull; docker build --rm -t figassis/wpoffload:$1 . && docker push figassis/wpoffload:$1