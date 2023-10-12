@echo off

echo updating... && del /s /q go.sum && go get -d -u -t ./... && go mod tidy