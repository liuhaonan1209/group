#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
高性能批量数据插入脚本
支持MySQL数据库，每个表插入100万条数据
"""

import pymysql
import random
import time
from datetime import datetime, timedelta
from concurrent.futures import ThreadPoolExecutor
import threading

# 数据库配置
DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': 'password',
    'database': 'bus_management',
    'charset': 'utf8mb4'
}

# 配置参数
TOTAL_RECORDS = 1000000  # 每个表插入100万条数据
BATCH_SIZE = 5000        # 每批插入5000条
THREAD_COUNT = 4         # 并发线程数

# 基础数据
CITIES = ['北京', '上海', '广州', '深圳', '杭州', '南京', '武汉', '成都', '重庆', '西安', '天津', '苏州', '长沙', '郑州', '青岛']
FLEETS = ['车队A', '车队B', '车队C', '车队D', '车队E', '车队F', '车队G', '车队H']
TAGS = ['快线', '普通', '夜班', '高峰', '节假日', '直达', '环线', '支线']
SURNAMES = ['张', '李', '王', '赵', '钱', '孙', '周', '吴', '郑', '冯', '陈', '褚', '卫', '蒋', '沈', '韩', '杨']
GIVEN_NAMES = ['伟', '芳', '娜', '敏', '静', '丽', '强', '磊', '军', '洋', '勇', '艳', '杰', '娟', '涛', '明', '超', '秀英', '华', '建国']

# 线程锁
print_lock = threading.Lock()

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
    area_codes = ['110101', '310101', '440101', '440301', '330101', '320101', '420101', '510101', '500101', '610101']
    area = random.choice(area_codes)
    year = random.randint(1970, 2000)
    month = random.randint(1, 12)
    day = random.randint(1, 28)
    birthday = f"{year:04d}{month:02d}{day:02d}"
    random_code = f"{random.randint(0, 9999):04d}"
    return area + birthday + random_code

def generate_phone():
    """生成随机手机号"""
    prefixes = ['130', '131', '132', '133', '134', '135', '136', '137', '138', '139',
                '150', '151', '152', '153', '155', '156', '157', '158', '159',
                '180', '181', '182', '183', '184', '185', '186', '187', '188', '189']
    return random.choice(prefixes) + f"{random.randint(0, 99999999):08d}"

def safe_print(message):
    """线程安全的打印函数"""
    with print_lock:
        print(f"[{datetime.now().strftime('%H:%M:%S')}] {message}")

def create_tables():
    """创建数据库表"""
    conn = get_connection()
    cursor = conn.cursor()
    
    # 创建表的SQL语句
    tables = {
        'start_stations': '''
            CREATE TABLE IF NOT EXISTS start_stations (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                name VARCHAR(255) NOT NULL UNIQUE,
                latitude FLOAT,
                longitude FLOAT,
                time_from_start DATETIME,
                address VARCHAR(255) NOT NULL,
                is_active TINYINT(1) DEFAULT 1,
                INDEX idx_name (name),
                INDEX idx_is_active (is_active)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'drivers': '''
            CREATE TABLE IF NOT EXISTS drivers (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                name VARCHAR(100) NOT NULL,
                tel VARCHAR(20) NOT NULL UNIQUE,
                id_card VARCHAR(18) NOT NULL UNIQUE,
                license VARCHAR(50) NOT NULL,
                register_date DATETIME NOT NULL,
                rating DECIMAL(3,2) DEFAULT 5.00,
                INDEX idx_tel (tel),
                INDEX idx_id_card (id_card)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'passengers': '''
            CREATE TABLE IF NOT EXISTS passengers (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                name VARCHAR(100) NOT NULL,
                tel VARCHAR(20) NOT NULL UNIQUE,
                id_card VARCHAR(18) NOT NULL UNIQUE,
                register_date DATETIME NOT NULL,
                INDEX idx_tel (tel),
                INDEX idx_id_card (id_card)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'routes': '''
            CREATE TABLE IF NOT EXISTS routes (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                route_no VARCHAR(20) NOT NULL,
                fieet VARCHAR(20),
                tage VARCHAR(50),
                tag_show_status VARCHAR(50),
                start_station_id INT,
                end_station_id INT,
                pass_stations VARCHAR(500),
                route_path VARCHAR(500),
                is_active VARCHAR(10),
                INDEX idx_route_no (route_no),
                INDEX idx_fieet (fieet)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'bus_schedules': '''
            CREATE TABLE IF NOT EXISTS bus_schedules (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                route_id INT NOT NULL,
                departure_time INT,
                arrival_time INT,
                capacity INT,
                is_active TINYINT(1) DEFAULT 1,
                INDEX idx_route_id (route_id)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'route_stops': '''
            CREATE TABLE IF NOT EXISTS route_stops (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                route_id INT NOT NULL,
                station_id INT NOT NULL,
                travel_time INT DEFAULT 0,
                stop_order INT,
                stop_type TINYINT DEFAULT 1,
                is_active TINYINT(1) DEFAULT 1,
                INDEX idx_route_station (route_id, station_id)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'transactions': '''
            CREATE TABLE IF NOT EXISTS transactions (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                trans_date DATE,
                trans_time DATETIME,
                order_no VARCHAR(50),
                driver_id INT,
                driver_name VARCHAR(50),
                passenger_id INT,
                passenger_name VARCHAR(50),
                start_station VARCHAR(100),
                end_station VARCHAR(100),
                amount DECIMAL(10,2),
                payment_method TINYINT,
                trans_type TINYINT,
                remark VARCHAR(255),
                INDEX idx_trans_date (trans_date),
                INDEX idx_order_no (order_no)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'balance_sheets': '''
            CREATE TABLE IF NOT EXISTS balance_sheets (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                settle_date DATE,
                order_income DECIMAL(12,2),
                order_refund DECIMAL(12,2),
                order_settle DECIMAL(12,2),
                status TINYINT DEFAULT 1,
                INDEX idx_settle_date (settle_date)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        ''',
        
        'driver_settles': '''
            CREATE TABLE IF NOT EXISTS driver_settles (
                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                deleted_at DATETIME NULL,
                settle_date DATE,
                driver_id INT,
                driver_name VARCHAR(50),
                phone VARCHAR(20),
                trip_count INT,
                total_income DECIMAL(12,2),
                commission DECIMAL(12,2),
                net_income DECIMAL(12,2),
                status TINYINT DEFAULT 1,
                INDEX idx_settle_date (settle_date),
                INDEX idx_driver_id (driver_id)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        '''
    }
    
    safe_print("开始创建数据库表...")
    for table_name, sql in tables.items():
        cursor.execute(sql)
        safe_print(f"表 {table_name} 创建完成")
    
    conn.commit()
    cursor.close()
    conn.close()
    safe_print("所有表创建完成")

def insert_stations_batch(start_id, batch_size):
    """批量插入站点数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        name = f"{random.choice(CITIES)}站{record_id}"
        latitude = 39.9 + random.uniform(-5, 5)
        longitude = 116.4 + random.uniform(-5, 5)
        time_from_start = datetime.now() - timedelta(days=random.randint(0, 365))
        address = f"{random.choice(CITIES)}市{random.choice(CITIES)}区{random.choice(CITIES)}路{random.randint(1, 999)}号"
        is_active = 1 if random.randint(1, 10) > 1 else 0
        
        values.append((name, latitude, longitude, time_from_start, address, is_active))
    
    sql = """
        INSERT INTO start_stations (name, latitude, longitude, time_from_start, address, is_active)
        VALUES (%s, %s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_drivers_batch(start_id, batch_size):
    """批量插入司机数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        name = generate_random_name()
        tel = generate_phone()
        id_card = generate_id_card()
        license = f"C1{random.randint(10000000, 99999999)}"
        register_date = datetime.now() - timedelta(days=random.randint(0, 1095))
        rating = round(4.0 + random.random(), 2)
        
        values.append((name, tel, id_card, license, register_date, rating))
    
    sql = """
        INSERT INTO drivers (name, tel, id_card, license, register_date, rating)
        VALUES (%s, %s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_passengers_batch(start_id, batch_size):
    """批量插入乘客数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        name = generate_random_name()
        tel = generate_phone()
        id_card = generate_id_card()
        register_date = datetime.now() - timedelta(days=random.randint(0, 730))
        
        values.append((name, tel, id_card, register_date))
    
    sql = """
        INSERT INTO passengers (name, tel, id_card, register_date)
        VALUES (%s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_routes_batch(start_id, batch_size):
    """批量插入路线数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        route_no = f"R{record_id:06d}"
        fieet = random.choice(FLEETS)
        tage = random.choice(TAGS)
        tag_show_status = random.choice(['show', 'hide'])
        start_station_id = random.randint(1, 1000)
        end_station_id = random.randint(1, 1000)
        pass_stations = f"{random.randint(1, 1000)},{random.randint(1, 1000)},{random.randint(1, 1000)}"
        route_path = f"path_{record_id}"
        is_active = str(random.randint(0, 1))
        
        values.append((route_no, fieet, tage, tag_show_status, start_station_id, 
                      end_station_id, pass_stations, route_path, is_active))
    
    sql = """
        INSERT INTO routes (route_no, fieet, tage, tag_show_status, start_station_id,
                           end_station_id, pass_stations, route_path, is_active)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_schedules_batch(start_id, batch_size):
    """批量插入班次数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        route_id = random.randint(1, 100000)
        departure_time = random.randint(0, 1440)  # 0-1440分钟
        arrival_time = departure_time + random.randint(30, 210)
        capacity = random.randint(20, 70)
        is_active = 1 if random.randint(1, 10) > 1 else 0
        
        values.append((route_id, departure_time, arrival_time, capacity, is_active))
    
    sql = """
        INSERT INTO bus_schedules (route_id, departure_time, arrival_time, capacity, is_active)
        VALUES (%s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_route_stops_batch(start_id, batch_size):
    """批量插入路线站点关联数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        route_id = random.randint(1, 100000)
        station_id = random.randint(1, 100000)
        travel_time = random.randint(0, 120)
        stop_order = random.randint(1, 20)
        stop_type = random.randint(1, 2)
        is_active = 1 if random.randint(1, 10) > 1 else 0
        
        values.append((route_id, station_id, travel_time, stop_order, stop_type, is_active))
    
    sql = """
        INSERT INTO route_stops (route_id, station_id, travel_time, stop_order, stop_type, is_active)
        VALUES (%s, %s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_transactions_batch(start_id, batch_size):
    """批量插入交易数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        trans_time = datetime.now() - timedelta(days=random.randint(0, 365))
        trans_date = trans_time.date()
        order_no = f"ORD{record_id:013d}"
        driver_id = random.randint(1, 100000)
        driver_name = generate_random_name()
        passenger_id = random.randint(1, 100000)
        passenger_name = generate_random_name()
        start_station = f"{random.choice(CITIES)}站"
        end_station = f"{random.choice(CITIES)}站"
        amount = round(random.uniform(5.0, 105.0), 2)
        payment_method = random.randint(1, 3)
        trans_type = random.randint(1, 3)
        remark = f"交易备注{record_id}"
        
        values.append((trans_date, trans_time, order_no, driver_id, driver_name,
                      passenger_id, passenger_name, start_station, end_station,
                      amount, payment_method, trans_type, remark))
    
    sql = """
        INSERT INTO transactions (trans_date, trans_time, order_no, driver_id, driver_name,
                                passenger_id, passenger_name, start_station, end_station,
                                amount, payment_method, trans_type, remark)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_balance_sheets_batch(start_id, batch_size):
    """批量插入收支对账数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        settle_date = (datetime.now() - timedelta(days=random.randint(0, 365))).date()
        order_income = round(random.uniform(100.0, 1000.0), 2)
        order_refund = round(random.uniform(0.0, 50.0), 2)
        order_settle = order_income - order_refund
        status = random.randint(1, 2)
        
        values.append((settle_date, order_income, order_refund, order_settle, status))
    
    sql = """
        INSERT INTO balance_sheets (settle_date, order_income, order_refund, order_settle, status)
        VALUES (%s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_driver_settles_batch(start_id, batch_size):
    """批量插入司机结算数据"""
    conn = get_connection()
    cursor = conn.cursor()
    
    values = []
    for i in range(batch_size):
        record_id = start_id + i
        settle_date = (datetime.now() - timedelta(days=random.randint(0, 365))).date()
        driver_id = random.randint(1, 100000)
        driver_name = generate_random_name()
        phone = generate_phone()
        trip_count = random.randint(10, 100)
        total_income = round(random.uniform(100.0, 500.0), 2)
        commission = round(total_income * 0.1, 2)
        net_income = total_income - commission
        status = random.randint(1, 2)
        
        values.append((settle_date, driver_id, driver_name, phone, trip_count,
                      total_income, commission, net_income, status))
    
    sql = """
        INSERT INTO driver_settles (settle_date, driver_id, driver_name, phone, trip_count,
                                  total_income, commission, net_income, status)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
    """
    
    cursor.executemany(sql, values)
    conn.commit()
    cursor.close()
    conn.close()

def insert_table_data(table_name, insert_func):
    """多线程插入表数据"""
    safe_print(f"开始插入 {table_name} 数据...")
    start_time = time.time()
    
    # 计算每个线程处理的批次数
    batches_per_thread = (TOTAL_RECORDS // BATCH_SIZE) // THREAD_COUNT
    
    def worker(thread_id):
        start_batch = thread_id * batches_per_thread
        end_batch = (thread_id + 1) * batches_per_thread
        if thread_id == THREAD_COUNT - 1:  # 最后一个线程处理剩余的批次
            end_batch = TOTAL_RECORDS // BATCH_SIZE
        
        for batch_idx in range(start_batch, end_batch):
            start_id = batch_idx * BATCH_SIZE + 1
            try:
                insert_func(start_id, BATCH_SIZE)
                if (batch_idx + 1) % 50 == 0:
                    safe_print(f"{table_name} 线程{thread_id}: 已完成 {batch_idx + 1} 批次")
            except Exception as e:
                safe_print(f"{table_name} 线程{thread_id} 批次{batch_idx + 1} 插入失败: {e}")
    
    # 启动多线程
    with ThreadPoolExecutor(max_workers=THREAD_COUNT) as executor:
        futures = [executor.submit(worker, i) for i in range(THREAD_COUNT)]
        for future in futures:
            future.result()
    
    elapsed_time = time.time() - start_time
    safe_print(f"{table_name} 数据插入完成，耗时: {elapsed_time:.2f}秒")

def main():
    """主函数"""
    print("=" * 60)
    print("班车管理系统 - 批量数据插入脚本")
    print(f"目标: 每个表插入 {TOTAL_RECORDS:,} 条数据")
    print(f"批次大小: {BATCH_SIZE}")
    print(f"并发线程: {THREAD_COUNT}")
    print("=" * 60)
    
    # 创建表
    create_tables()
    
    # 插入数据的表和对应的插入函数
    tables_to_insert = [
        ('start_stations', insert_stations_batch),
        ('drivers', insert_drivers_batch),
        ('passengers', insert_passengers_batch),
        ('routes', insert_routes_batch),
        ('bus_schedules', insert_schedules_batch),
        ('route_stops', insert_route_stops_batch),
        ('transactions', insert_transactions_batch),
        ('balance_sheets', insert_balance_sheets_batch),
        ('driver_settles', insert_driver_settles_batch),
    ]
    
    total_start_time = time.time()
    
    # 依次插入各表数据
    for table_name, insert_func in tables_to_insert:
        insert_table_data(table_name, insert_func)
        print("-" * 60)
    
    total_elapsed_time = time.time() - total_start_time
    total_records = len(tables_to_insert) * TOTAL_RECORDS
    
    print("=" * 60)
    print("数据插入完成！")
    print(f"总记录数: {total_records:,}")
    print(f"总耗时: {total_elapsed_time:.2f}秒")
    print(f"平均速度: {total_records/total_elapsed_time:.0f} 记录/秒")
    print("=" * 60)

if __name__ == "__main__":
    main()