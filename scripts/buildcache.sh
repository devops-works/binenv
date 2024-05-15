#!/bin/bash

#
# Usage: ./updatecache.sh
#
# This script will:
# - copy the source code distribution locally
# - rebuild the cache
# - copy the cache to source code
#
# PLEASE MAKE SURE YOUR HAVE A GITHUB TOKEN !!!

set -eu

export BINENV_GLOBAL=false 

if [ "$GITHUB_TOKEN" == "" ]; then
    echo "ERROR: a GITHUB_TOKEN is required but can not be found; exiting"
fi

make

echo "Copying distribution to your config directory"
cp distributions/distributions.yaml ~/.config/binenv/

echo "Updating the cache (5 threads)"
./bin/binenv update -f -c5

echo "Importing resulting cache into code"
cat ~/.cache/binenv/cache.json | jq '.' > distributions/cache.json

echo "Please test the cache using './scripts/validate.sh code'"
