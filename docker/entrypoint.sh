#!/bin/bash

export PATH=~/.binenv:$PATH

for dist in "${BINENV_INSTALL}"; do
    binenv update -f $dist > /dev/null 2>&1
    binenv install $dist > /dev/null 2>&1
done

exec "$@"
