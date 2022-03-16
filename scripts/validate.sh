#!/bin/bash

# Validates distribution & cache consistency
#
# Usage: ./validate.sh [mode]
#
# where mode can be:
# - local (checks your local distribution & cache file)
# - code (checks source code distribution & cache file)

MODE=${1:-code}

# Set files to use
DISTRIBUTION=~/.config/binenv/distributions.yaml
CACHE=~/.cache/binenv/cache.json

if [ ${MODE} == "code" ]; then
    DISTRIBUTION=distributions/distributions.yaml
    CACHE=distributions/cache.json
fi

RED='\033[0;31m'
NC='\033[0m' # No Color

# Check which yq version we have
# Start with Go version
YQCMD='yq -M eval'

# Try python version
chk=$(yq -h | grep "See the manpage for more options" 2>&1 || true)

if [ "$chk" == "See the manpage for more options." ]; then
    YQCMD='yq . --yaml-output'
fi

count=0
errors=0

echo "Using distribution $DISTRIBUTION"
echo "Using cache        $CACHE"

for i in $($YQCMD '.sources | keys' "${DISTRIBUTION}" | sed -e 's/^- //' | grep -v '#'); do
    # dist=$(echo $i | cut -f2 -d' ')
    if [[ $i == \#* ]]; then
      continue
    fi
    count=$((count + 1))

    echo -n "checking versions for $i"
    vcount=$(cat "${CACHE}" | jq ".\"${i}\"[]" 2> /dev/null| wc -l)

    if [ $vcount -eq 0 ]; then
        echo -e "...${RED}found none${NC}"
        errors=$((errors + 1))
    else
        echo "...found $vcount versions"
    fi
done

echo "found $count distributions"

if [ $errors -ne 0 ]; then
    echo "got $errors errors"
    exit 1
fi