#!/bin/bash -xe
SCRIPT_DIR=$(cd $(dirname $0); pwd)

cd $SCRIPT_DIR

/usr/bin/git pull origin master

/usr/local/bin/docker-compose build

/usr/local/bin/docker-compose stop

nohup /usr/local/bin/docker-compose up -d
