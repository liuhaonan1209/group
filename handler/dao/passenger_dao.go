// Package dao 数据访问层
// 提供乘客相关的数据库操作方法
package dao

import (
	"group/global"
)

// GetPassengerByID 根据ID查询乘客详情
// 使用泛型实现，可以适配不同的数据模型
// 参数:
//   - data: 数据模型指针，用于指定查询的表结构
//   - id: 乘客ID
// 返回:
//   - passengerDetail: 查询到的乘客详情
//   - err: 错误信息
func GetPassengerByID[T any](data *T, id int64) (passengerDetail T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).Find(&passengerDetail).Error
	return
}

// GetPassengerByTelAndIDCard 根据手机号和身份证查询乘客
// 用于乘客身份验证，需要二要素完全匹配
// 参数:
//   - data: 数据模型指针
//   - tel: 手机号
//   - idCard: 身份证号
// 返回:
//   - passengerDetail: 查询到的乘客详情
//   - err: 错误信息（如果未找到匹配记录会返回错误）
func GetPassengerByTelAndIDCard[T any](data *T, tel, idCard string) (passengerDetail T, err error) {
	err = global.DB.Model(data).Where("tel = ? AND id_card = ?", tel, idCard).First(&passengerDetail).Error
	return
}

// GetPassengerByTel 根据手机号查询乘客
// 用于检查手机号是否已注册
// 参数:
//   - data: 数据模型指针
//   - tel: 手机号
// 返回:
//   - passengerDetail: 查询到的乘客详情
//   - err: 错误信息（如果未找到匹配记录会返回错误）
func GetPassengerByTel[T any](data *T, tel string) (passengerDetail T, err error) {
	err = global.DB.Model(data).Where("tel = ?", tel).First(&passengerDetail).Error
	return
}

// CreatePassenger 创建乘客
// 向数据库插入新的乘客记录
// 参数:
//   - data: 乘客数据模型指针，包含要插入的数据
// 返回:
//   - error: 错误信息（如果插入失败，如唯一索引冲突）
func CreatePassenger[T any](data *T) error {
	return global.DB.Create(data).Error
}

// UpdatePassenger 更新乘客信息
// 保存乘客的修改信息到数据库
// 参数:
//   - data: 乘客数据模型指针，包含要更新的数据
// 返回:
//   - error: 错误信息
func UpdatePassenger[T any](data *T) error {
	return global.DB.Save(data).Error
}

// DeletePassenger 删除乘客
// 根据ID删除乘客记录（软删除）
// 参数:
//   - data: 数据模型指针
//   - id: 要删除的乘客ID
// 返回:
//   - error: 错误信息
func DeletePassenger[T any](data *T, id int64) error {
	return global.DB.Delete(data, id).Error
}
