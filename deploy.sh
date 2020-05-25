#!/bin/bash -xe
SCRIPT_DIR=$(cd $(dirname $0); pwd)

if [[ "$IMAGE_VER" != "" ]]
  then
    sed "s/IMAGE_VER=.*/IMAGE_VER=$IMAGE_VER/" .env > .env.bk && mv .env.bk .env
fi

/usr/bin/git pull origin master

/usr/local/bin/docker-compose build

/usr/local/bin/docker-compose stop

nohup /usr/local/bin/docker-compose up -d

# https://stackoverflow.com/questions/33913020/docker-remove-none-tag-images
/usr/bin/docker image prune