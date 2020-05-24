#!/bin/bash -xe
/usr/bin/git pull origin master

/usr/local/bin/docker-compose build

nohup /usr/local/bin/docker-compose up -d