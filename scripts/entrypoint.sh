#!/bin/bash

export PATH=~/.binenv:$PATH

GREEN='\033[0;32m'
RED='\033[0;31m'
WHITE='\033[0;37m'
RESET='\033[0m'

if [ -z ${GITHUB_TOKEN+x} ]; then
    echo "GITHUB_TOKEN is not set; can not continue"
    exit 1
fi

function strip_ansi() {
  shopt -s extglob
  printf %s "${1//$'\e'\[*([0-9;])m/}"
}

echo "Updating local distributions cache"

binenv update -f

count=$(binenv search | wc -l)
line='..........................'
errorscount=0

echo "Installing ${count} distributions"
for i in $(binenv search | cut -f1 -d':'); do
    dist=$(strip_ansi $i)
    err=$(binenv install $dist 2> >(grep -i ERR))
    case "$?" in
    0)
        printf "%s %s ${GREEN}OK${RESET}\n" $dist "${line:${#dist}}"
        ;;
    *)
        errorscount=$((errorscount+1))
        echo $err >> /tmp/error.log
        err=$(echo "${err}" | cut -f2 -d'=')
        printf "%s %s ${RED}ERROR${RESET}: ${err}\n" $dist "${line:${#dist}}"
        ;;
    esac
done

if [ ${errorscount} -gt 0 ]; then
    echo -e "\n\n${errorscount} errors found for ${count} distributions\n\n"
    cat /tmp/error.log
    exit 1
fi

exit 0
