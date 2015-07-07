#!/bin/bash

go build

docker build -t toomore/mailmanweb .

rm -rf ./webserver
