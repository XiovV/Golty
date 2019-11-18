#!/bin/ash

DEFAULT_UID=1000
DEFAULT_GID=1000
DEFAULT_UMASK=022

usermod goautoyt -u ${PUID:-${DEFAULT_UID}}
groupmod goautoyt -g ${PGID:-${DEFAULT_GID}}

chown -R goautoyt:goautoyt /app
exec su -s /bin/ash -c "umask ${UMASK_SET:-${DEFAULT_UMASK}};PATH=$PATH:/usr/local/bin;./main" goautoyt
