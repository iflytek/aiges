@echo off
go test . -cover -coverprofile cover.out
go tool cover -html=cover.out