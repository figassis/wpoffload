#!/bin/bash
clear

tag=`./version.sh`
git pull; docker build --rm -t figassis/wpoffload:$tag . && docker push figassis/wpoffload:$tag