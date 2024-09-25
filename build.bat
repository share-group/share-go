@echo off

echo updating... && (del /s /q go.sum || echo 1) && go get -u -t ./... && go mod tidy