@echo off
rem %1 是第一个参数，即 .proto 文件的路径
set PROTO_PATH=%~dp1
protoc --proto_path=%PROTO_PATH% --go_out=%PROTO_PATH% --go_opt=paths=source_relative --go-grpc_out=%PROTO_PATH% --go-grpc_opt=paths=source_relative %1