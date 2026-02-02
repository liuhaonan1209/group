#!/bin/bash

echo "🚀 开始运行数据插入脚本..."
echo "📋 默认参数: 每个表100万条数据，每批2000条"
echo "⚙️ 可以使用参数修改: -n 总数量 -b 批次大小"
echo ""

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 未找到Go环境，请先安装Go"
    exit 1
fi

echo "✅ Go环境检查通过"
echo ""

# 运行数据插入脚本
echo "🔄 正在执行数据插入..."
go run seed_data.go "$@"

echo ""
echo "🎉 数据插入脚本执行完成！"