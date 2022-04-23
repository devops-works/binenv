#!/bin/bash

export PATH=~/.binenv:$PATH
export BINENV_GLOBAL=false

set -u

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

if [ -z ${GITHUB_TOKEN+x} ]; then
    echo "GITHUB_TOKEN is not set; can not continue"
    exit 1
fi

check_links() {
  local errors=0
  local count=0
  echo -e "\n# Checking DISTRIBUTIONS.md URLs\n"
  while read l; do
    url=$(echo $l | cut -f2 -d'(' | cut -f 1 -d')' )
    name=$(echo $l | cut -f2 -d'[' | cut -f 1 -d']' )
    printf "%-70s %s" "[$name] $url"
    # printf "%s %-50s %s" $name $url
    ((count=count+1))
    curl -sfLIm5 $url -o /dev/null
    if [ $? -gt 0 ]; then
      echo -e "${RED}fail${NC}"
      ((errors=errors+1))
    else
      echo -e "${GREEN}ok${NC}"
    fi
  done <  ~/.config/binenv/DISTRIBUTIONS.md 

  echo -e "\n$count URLs tested"
  if [ $errors -gt 0 ]; then
    echo -e "${RED}$errors errors found${NC}"
  fi
}

check_update() {
  local errors=0
  local count=0
  echo -e "\n# Checking distributions release updates\n"
  for dist in $(grep "^\s\s[a-z]" ~/.config/binenv/distributions.yaml | cut -f3 -d' ' | tr -d ':'); do
    printf "testing %-28s %s" $dist
    ((count=count+1))
    if ! binenv update -f $dist > /dev/null 2&>1; then
      echo -e "${RED}fail${NC}"
      ((errors=errors+1))
    else
      echo -e "${GREEN}ok${NC}"
    fi
  done

  echo -e "\n$count releases updated"
  if [ $errors -gt 0 ]; then
    echo -e "${RED}$errors errors found${NC}"
  fi
}

check_install() {
  local errors=0
  local count=0
  echo -e "\n# Checking distributions install\n"
  for dist in $(grep "^\s\s[a-z]" ~/.config/binenv/distributions.yaml |cut -f3 -d' ' | tr -d ':'); do
    printf "testing %-28s %s" $dist
    ((count=count+1))
    if ! binenv install $dist > /dev/null 2>&1; then
      echo -e "${RED}fail${NC}"
      ((errors=errors+1))
    else
      echo -e "${GREEN}ok${NC}"
    fi
  done

  echo -e "\n$count distributions installed"
  if [ $errors -gt 0 ]; then
    echo -e "${RED}$errors errors found${NC}"
  fi
}

echo

binenv version

check_links
check_update
check_install