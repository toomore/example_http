#!/bin/bash

go build

docker build -t toomore/mailman .

rm -rf ./mailman
