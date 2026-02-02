@echo off
chcp 65001 >nul
title 使用特殊标志运行 Go 程序

echo 🔧 使用特殊编译标志运行程序...

REM 设置 Go 环境变量，禁用一些可能触发安全策略的特性
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

echo 📋 当前目录: %cd%
echo 🎯 目标: api\routerApi

cd api\routerApi

echo 🔄 使用安全编译标志运行...
go run -ldflags="-s -w" -trimpath .

pause