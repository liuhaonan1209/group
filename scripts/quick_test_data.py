#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
快速测试数据插入脚本
每个表插入1万条数据，用于快速测试
"""

import pymysql
import random
import time
from datetime import datetime, timedelta

# 数据库配置
DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': 'password',
    'database': 'bus_management',
    'charset': 'utf8mb4'
}

# 配置参数 - 快速测试版本
TOTAL_RECORDS = 10000    # 每个表插入1万条数据
BATCH_SIZE = 1000        # 每批插入1000条

# 基础数据
CITIES = ['北京', '上海', '广州', '深圳', '杭州', '南京', '武汉', '成都', '重庆', '西安']
FLEETS = ['车队A', '车队B', '车队C', '车队D', '车队E']
TAGS = ['快线', '普通', '夜班', '高峰', '节假日']
SURNAMES = ['张', '李', '王', '赵', '钱', '孙', '周', '吴', '郑', '冯', '陈', '褚']
GIVEN_NAMES = ['伟', '芳', '娜', '敏', '静', '丽', '强', '磊', '军', '洋', '勇', '艳']

def get_connection():
    """获取数据库连接"""
    return pymysql.connect(**DB_CONFIG)

def generate_random_name():
    """生成随机姓名"""
    surname = random.choice(SURNAMES)
    given_name = random.choice(GIVEN_NAMES)
    if random.randint(0, 1):
        given_name += random.choice(GIVEN_NAMES)
    return surname + given_name

def generate_id_card():
    """生成随机身份证号"""
    area_codes = ['110101', '310101', '440101', '440301', '330101']
    area = random.choice(area_codes)
    year = random.randint(1970, 2000)
    month = random.randint(1, 12)
    day = random.randint(1, 28)
    birthday = f"{year:04d}{month:02d}{day:02d}"
    random_code = f"{random.randint(0, 9999):04d}"
    return area + birthday + random_code

def generate_phone():
    """生成随机手机号"""
    prefixes = ['130', '131', '132', '133', '134', '135', '136', '137', '138', '139']
    return random.choice(prefixes) + f"{random.randint(0, 99999999):08d}"

def create_tables():
    """创建数据库表"""
    conn = get_connection()
    cursor = conn.cursor()
    
    # 简化的表结构
    tables = {
        'test_stations': '''
            CREATE TABLE IF NOT EXISTS test_stations (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                latitude FLOAT,
                longitude FLOAT,
                address VARCHAR(255),
                is_active TINYINT(1) DEFAULT 1,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'test_drivers': '''
            CREATE TABLE IF NOT EXISTS test_drivers (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                name VARCHAR(100) NOT NULL,
                tel VARCHAR(20) NOT NULL,
                id_card VARCHAR(18) NOT NULL,
                rating DECIMAL(3,2) DEFAULT 5.00,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'test_routes': '''
            CREATE TABLE IF NOT EXISTS test_routes (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                route_no VARCHAR(20) NOT NULL,
                fleet VARCHAR(20),
                tag VARCHAR(50),
                start_station_id INT,
                end_station_id INT,
                is_active TINYINT(1) DEFAULT 1,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'test_transactions': '''
            CREATE TABLE IF NOT EXISTS test_transactions (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                order_no VARCHAR(50),
                driver_name VARCHAR(50),
                passenger_name VARCHAR(50),
                amount DECIMAL(10,2),
                trans_date DATE,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        '''
    }
    
    print("创建测试表...")
    for table_name, sql in tables.items():
        cursor.execute(sql)
        print(f"表 {table_name} 创建完成")
    
    conn.commit()
    cursor.close()
    conn.close()

def insert_test_data():
    """插入测试数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    print(f"开始插入测试数据，每个表 {TOTAL_RECORDS:,} 条记录...")
    start_time = time.time()
    
    # 插入站点数据
    print("插入站点数据...")
    for batch in range(0, TOTAL_RECORDS, BATCH_SIZE):
        values = []
        for i in range(BATCH_SIZE):
            if batch + i >= TOTAL_RECORDS:
                break
            record_id = batch + i + 1
            values.append((
                f"{random.choice(CITIES)}站{record_id}",
                39.9 + random.uniform(-2, 2),
                116.4 + random.uniform(-2, 2),
                f"{random.choice(CITIES)}市{random.choice(CITIES)}区{record_id}号"
            ))
        
        if values:
            cursor.executemany(
                "INSERT INTO test_stations (name, latitude, longitude, address) VALUES (%s, %s, %s, %s)",
                values
            )
            conn.commit()
    
    # 插入司机数据
    print("插入司机数据...")
    for batch in range(0, TOTAL_RECORDS, BATCH_SIZE):
        values = []
        for i in range(BATCH_SIZE):
            if batch + i >= TOTAL_RECORDS:
                break
            values.append((
                generate_random_name(),
                generate_phone(),
                generate_id_card(),
                round(4.0 + random.random(), 2)
            ))
        
        if values:
            cursor.executemany(
                "INSERT INTO test_drivers (name, tel, id_card, rating) VALUES (%s, %s, %s, %s)",
                values
            )
            conn.commit()
    
    # 插入路线数据
    print("插入路线数据...")
    for batch in range(0, TOTAL_RECORDS, BATCH_SIZE):
        values = []
        for i in range(BATCH_SIZE):
            if batch + i >= TOTAL_RECORDS:
                break
            record_id = batch + i + 1
            values.append((
                f"R{record_id:06d}",
                random.choice(FLEETS),
                random.choice(TAGS),
                random.randint(1, 1000),
                random.randint(1, 1000)
            ))
        
        if values:
            cursor.executemany(
                "INSERT INTO test_routes (route_no, fleet, tag, start_station_id, end_station_id) VALUES (%s, %s, %s, %s, %s)",
                values
            )
            conn.commit()
    
    # 插入交易数据
    print("插入交易数据...")
    for batch in range(0, TOTAL_RECORDS, BATCH_SIZE):
        values = []
        for i in range(BATCH_SIZE):
            if batch + i >= TOTAL_RECORDS:
                break
            record_id = batch + i + 1
            trans_date = (datetime.now() - timedelta(days=random.randint(0, 365))).date()
            values.append((
                f"ORD{record_id:010d}",
                generate_random_name(),
                generate_random_name(),
                round(random.uniform(10.0, 100.0), 2),
                trans_date
            ))
        
        if values:
            cursor.executemany(
                "INSERT INTO test_transactions (order_no, driver_name, passenger_name, amount, trans_date) VALUES (%s, %s, %s, %s, %s)",
                values
            )
            conn.commit()
    
    cursor.close()
    conn.close()
    
    elapsed_time = time.time() - start_time
    total_records = 4 * TOTAL_RECORDS
    
    print(f"\n测试数据插入完成！")
    print(f"总记录数: {total_records:,}")
    print(f"耗时: {elapsed_time:.2f}秒")
    print(f"速度: {total_records/elapsed_time:.0f} 记录/秒")

def verify_data():
    """验证数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    print("\n验证数据...")
    tables = ['test_stations', 'test_drivers', 'test_routes', 'test_transactions']
    
    for table in tables:
        cursor.execute(f"SELECT COUNT(*) FROM {table}")
        count = cursor.fetchone()[0]
        print(f"{table}: {count:,} 条记录")
        
        # 显示示例数据
        cursor.execute(f"SELECT * FROM {table} LIMIT 3")
        rows = cursor.fetchall()
        print(f"  示例数据: {len(rows)} 条")
        for row in rows:
            print(f"    {row}")
        print()
    
    cursor.close()
    conn.close()

def main():
    """主函数"""
    print("=" * 50)
    print("班车管理系统 - 快速测试数据插入")
    print(f"每个表插入: {TOTAL_RECORDS:,} 条数据")
    print("=" * 50)
    
    try:
        # 创建表
        create_tables()
        
        # 插入数据
        insert_test_data()
        
        # 验证数据
        verify_data()
        
        print("快速测试数据生成完成！")
        
    except Exception as e:
        print(f"错误: {e}")
        print("请检查数据库连接配置和权限")

if __name__ == "__main__":
    main()