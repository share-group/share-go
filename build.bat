@echo off

echo updating... && (del /s /q go.sum || echo 1) && go get -d -u -t ./... && go mod tidy