#!/bin/sh
if [ -z $AUTHORIZATION_DSN ] ; then
  redis-server /etc/redis.conf
fi

for file in /etc/nginx/http.d/template/*
do
    if test -f "$file"
    then
      envsubst "$(env | awk -F = '{printf " $%s", $1}')" < "$file" > /etc/nginx/http.d/"${file##*/}"
      echo "${file##*/}"
    fi
done

nginx -g "daemon on;"
# nohup /nginx & > /gateway-console.log 2>&1
exec "$@"
