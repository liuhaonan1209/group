#!/bin/bash

# 班车管理系统批量数据插入脚本
# 每个表插入100万条数据

echo "=========================================="
echo "班车管理系统 - 批量数据插入"
echo "=========================================="

# 检查Python环境
if ! command -v python3 &> /dev/null; then
    echo "错误: 未找到 python3，请先安装 Python 3"
    exit 1
fi

# 检查并安装依赖
echo "检查Python依赖..."
pip3 install pymysql

# 设置数据库连接参数
echo "请确保MySQL数据库已启动，并且配置正确："
echo "- 主机: localhost"
echo "- 端口: 3306"
echo "- 用户: root"
echo "- 密码: password"
echo "- 数据库: bus_management"
echo ""

read -p "是否继续执行数据插入？(y/N): " confirm
if [[ $confirm != [yY] ]]; then
    echo "操作已取消"
    exit 0
fi

# 执行数据插入
echo "开始执行批量数据插入..."
python3 bulk_insert_sql.py

echo "批量数据插入完成！"