#!/bin/bash

env GOOS=linux GOARCH=386 go build -o homebugh-api ./cmd/api &&
  rsync -a homebugh-api deploy@homebugh.info:/home/deploy/apps/homebugh-api/
