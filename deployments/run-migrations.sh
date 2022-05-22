#!/bin/bash

# database needs to be created first - this is done in db-init/01-db-init.sql
docker container run --rm -it --net=host -v "$PWD/../migrations":/ledger -e PASSWORD=password1234 sledger:latest --database postgres://postgres:password1234@localhost:5432/chronicler?sslmode=disable
