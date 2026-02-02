@echo off
chcp 65001 >nul
title 班车管理系统 - 批量数据插入

echo ==========================================
echo 班车管理系统 - 批量数据插入
echo ==========================================
echo.

REM 检查Python环境
python --version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 Python，请先安装 Python 3
    pause
    exit /b 1
)

REM 检查并安装依赖
echo 检查Python依赖...
pip install pymysql

echo.
echo 请确保MySQL数据库已启动，并且配置正确：
echo - 主机: localhost
echo - 端口: 3306
echo - 用户: root
echo - 密码: password
echo - 数据库: bus_management
echo.

set /p confirm="是否继续执行数据插入？(y/N): "
if /i not "%confirm%"=="y" (
    echo 操作已取消
    pause
    exit /b 0
)

echo.
echo 开始执行批量数据插入...
python bulk_insert_sql.py

echo.
echo 批量数据插入完成！
pause