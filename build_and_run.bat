@echo off
chcp 65001 >nul
title 编译并运行 Go 程序

echo 🔨 开始编译程序...

REM 设置输出目录
set OUTPUT_DIR=bin
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

REM 编译各个服务
echo 📦 编译 routerApi...
cd api\routerApi
go build -o ..\..\%OUTPUT_DIR%\routerApi.exe .
if errorlevel 1 (
    echo ❌ routerApi 编译失败
    pause
    exit /b 1
)
cd ..\..

echo 📦 编译 financeApi...
cd api\financeApi
go build -o ..\..\%OUTPUT_DIR%\financeApi.exe .
if errorlevel 1 (
    echo ❌ financeApi 编译失败
    pause
    exit /b 1
)
cd ..\..

echo 📦 编译 routeRpc...
cd rpc\routeRpc
go build -o ..\..\%OUTPUT_DIR%\routeRpc.exe .
if errorlevel 1 (
    echo ❌ routeRpc 编译失败
    pause
    exit /b 1
)
cd ..\..

echo 📦 编译 financeRpc...
cd rpc\financeRpc
go build -o ..\..\%OUTPUT_DIR%\financeRpc.exe .
if errorlevel 1 (
    echo ❌ financeRpc 编译失败
    pause
    exit /b 1
)
cd ..\..

echo 📦 编译 seed_data...
go build -o %OUTPUT_DIR%\seed_data.exe seed_data.go
if errorlevel 1 (
    echo ❌ seed_data 编译失败
    pause
    exit /b 1
)

echo ✅ 所有程序编译完成！

echo.
echo 🚀 可执行文件位置:
echo   - %OUTPUT_DIR%\routerApi.exe
echo   - %OUTPUT_DIR%\financeApi.exe  
echo   - %OUTPUT_DIR%\routeRpc.exe
echo   - %OUTPUT_DIR%\financeRpc.exe
echo   - %OUTPUT_DIR%\seed_data.exe

echo.
echo 📋 运行说明:
echo   1. 先启动 RPC 服务: routeRpc.exe 和 financeRpc.exe
echo   2. 再启动 API 服务: routerApi.exe 和 financeApi.exe
echo   3. 运行数据插入: seed_data.exe

pause