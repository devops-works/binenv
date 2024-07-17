#!/bin/bash

#
# Usage: ./updatecache.sh
#
# This script will rebuild the cache for all distributions in the distributions
# file.
#
# PLEASE MAKE SURE YOUR HAVE A GITHUB TOKEN !!!
#

set -eu

export BINENV_GLOBAL=false

if [ "$GITHUB_TOKEN" == "" ]; then
    echo "ERROR: a GITHUB_TOKEN is required but can not be found; exiting"
fi

make

echo "Updating the cache (5 threads)"
./bin/binenv update --cachedir ./distributions --confdir ./distributions -f -c5

echo "Importing resulting cache"
jq '.' < distributions/cache.json > distributions/cache.json.tmp
mv distributions/cache.json.tmp distributions/cache.json

echo "Please test the cache using './scripts/validate.sh code'"
