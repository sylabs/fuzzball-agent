#!/bin/sh

getent group fuzzball >/dev/null || groupadd -r fuzzball
getent passwd fuzzball >/dev/null || \
useradd -r -g fuzzball -d /var/run/fuzzball \
    -s /usr/sbin/nologin -c "fuzzball daemon" fuzzball
exit 0