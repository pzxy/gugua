#!/bin/bash
go mod tidy
go build -o gua gua.go
./gua &