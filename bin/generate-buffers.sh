#!/bin/sh
protoc -I=./api/pokemon --go_out=plugins=grpc:./api/pokemon ./api/pokemon/pokemon.proto