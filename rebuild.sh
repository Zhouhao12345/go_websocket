#! /bin/bash
docker-compose build
docker-compose stop && docker-compose rm -f && docker-compose up -d
