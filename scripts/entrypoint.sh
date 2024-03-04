#!/bin/bash

go build -o /app/service.exe /app/cmd/server/cmd/main.go

/app/service.exe
