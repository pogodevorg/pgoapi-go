#!/bin/sh
protoc -I=./api/pokemon --go_out=./api/pokemon ./api/pokemon/pokemon.proto