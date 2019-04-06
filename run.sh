#!/usr/bin/env bash
export CGO_LDFLAGS="-lmecab -lstdc++"
go run main.go