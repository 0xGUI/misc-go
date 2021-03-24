#!/bin/bash
set -x 
 
export GO111MODULE=auto

# https://github.com/hashicorp/waypoint/issues/572 
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o misc-go -ldflags "-s -w" . && \

upx misc-go && \

docker build -t misc-go . && \

kind load docker-image misc-go && \

kubectl delete -f pod.yaml 

kubectl apply -f pod.yaml
