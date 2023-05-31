#!/bin/bash

go build
go run main.go -in dirty_tmpl -out dirty_out -import github.com/hanxi/dirty-go/dirty_tmpl

