#!/bin/bash -xe

/usr/local/bin/docker-compose build

nohup /usr/local/bin/docker-compose up -d