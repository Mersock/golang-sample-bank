#!/bin/sh

go test $(go list ./... | grep -v /tools) -v -cover 
