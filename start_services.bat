@echo off
chcp 65001 >nul
title 启动班车管理系统服务

echo 🚀 启动班车管理系统服务...

REM 检查可执行文件是否存在
if not exist "bin\routeRpc.exe" (
    echo ❌ 可执行文件不存在，请先运行 build_and_run.bat 编译
    pause
    exit /b 1
)

echo 📋 启动顺序:
echo   1. routeRpc (端口: 8892)
echo   2. financeRpc (端口: 8894)  
echo   3. routerApi (端口: 8890)
echo   4. financeApi (端口: 8895)

echo.
echo ⚠️  请按顺序在新窗口中启动服务...

echo.
echo 🔄 启动 routeRpc 服务...
start "Route RPC Service" cmd /k "cd /d %cd% && bin\routeRpc.exe"
timeout /t 3 /nobreak >nul

echo 🔄 启动 financeRpc 服务...
start "Finance RPC Service" cmd /k "cd /d %cd% && bin\financeRpc.exe"
timeout /t 3 /nobreak >nul

echo 🔄 启动 routerApi 服务...
start "Router API Service" cmd /k "cd /d %cd% && bin\routerApi.exe"
timeout /t 3 /nobreak >nul

echo 🔄 启动 financeApi 服务...
start "Finance API Service" cmd /k "cd /d %cd% && bin\financeApi.exe"
timeout /t 3 /nobreak >nul

echo.
echo ✅ 所有服务启动完成！
echo.
echo 📊 服务地址:
echo   - Route API: http://localhost:8890
echo   - Finance API: http://localhost:8895
echo   - Frontend: http://localhost:3000 (需要单独启动)
echo.
echo 💡 提示:
echo   - 关闭服务请直接关闭对应的命令行窗口
echo   - 查看日志请查看各服务窗口的输出
echo   - 如需重启服务，请先关闭再重新运行此脚本

pause