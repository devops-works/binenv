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

echo "Copying distribution to your config directory"
cp distributions/distributions.yaml ~/.config/binenv/

echo "Updating the cache (2 threads)"
binenv update -f -c2

echo "Importing resulting cache into code"
cat ~/.cache/binenv/cache.json | jq '.' > distributions/cache.json

echo "Please test the cache using `./validate.sh code`"
