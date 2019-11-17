#!/bin/ash

DEFAULT_UID=1000
DEFAULT_GID=1000

usermod goautoyt -u ${PUID:-${DEFAULT_UID}}
groupmod goautoyt -g ${PGID:-${DEFAULT_GID}}

chown -R goautoyt:goautoyt /app
exec su -s /bin/ash -c "PATH=$PATH:/usr/local/bin;./main" goautoyt
