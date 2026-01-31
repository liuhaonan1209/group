// Package dao 数据访问层
// 提供司机相关的数据库操作方法
package dao

import (
	"group/global"
)

// GetDriverByID 根据ID查询司机详情
// 使用泛型实现，可以适配不同的数据模型
// 参数:
//   - data: 数据模型指针，用于指定查询的表结构
//   - id: 司机ID
// 返回:
//   - driverDetail: 查询到的司机详情
//   - err: 错误信息
func GetDriverByID[T any](data *T, id int64) (driverDetail T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).Find(&driverDetail).Error
	return
}

// GetDriverByTelAndIDCardAndLicense 根据手机号、身份证和驾驶证查询司机
// 用于司机身份验证，需要三要素完全匹配
// 参数:
//   - data: 数据模型指针
//   - tel: 手机号
//   - idCard: 身份证号
//   - license: 驾驶证号
// 返回:
//   - driverDetail: 查询到的司机详情
//   - err: 错误信息（如果未找到匹配记录会返回错误）
func GetDriverByTelAndIDCardAndLicense[T any](data *T, tel, idCard, license string) (driverDetail T, err error) {
	err = global.DB.Model(data).Where("tel = ? AND id_card = ? AND license = ?", tel, idCard, license).First(&driverDetail).Error
	return
}

// GetDriverByTel 根据手机号查询司机
// 用于检查手机号是否已注册
// 参数:
//   - data: 数据模型指针
//   - tel: 手机号
// 返回:
//   - driverDetail: 查询到的司机详情
//   - err: 错误信息（如果未找到匹配记录会返回错误）
func GetDriverByTel[T any](data *T, tel string) (driverDetail T, err error) {
	err = global.DB.Model(data).Where("tel = ?", tel).First(&driverDetail).Error
	return
}

// CreateDriver 创建司机
// 向数据库插入新的司机记录
// 参数:
//   - data: 司机数据模型指针，包含要插入的数据
// 返回:
//   - error: 错误信息（如果插入失败，如唯一索引冲突）
func CreateDriver[T any](data *T) error {
	return global.DB.Create(data).Error
}

// UpdateDriver 更新司机信息
// 保存司机的修改信息到数据库
// 参数:
//   - data: 司机数据模型指针，包含要更新的数据
// 返回:
//   - error: 错误信息
func UpdateDriver[T any](data *T) error {
	return global.DB.Save(data).Error
}

// DeleteDriver 删除司机
// 根据ID删除司机记录（软删除）
// 参数:
//   - data: 数据模型指针
//   - id: 要删除的司机ID
// 返回:
//   - error: 错误信息
func DeleteDriver[T any](data *T, id int64) error {
	return global.DB.Delete(data, id).Error
}

// ==================== DriverConfig 司机配置相关操作 ====================

// GetDriverConfigByDriverID 根据司机ID查询配置信息
// 参数:
//   - data: 配置数据模型指针
//   - driverId: 司机ID
// 返回:
//   - config: 配置信息对象
//   - err: 查询失败或记录不存在时返回错误
func GetDriverConfigByDriverID[T any](data *T, driverId int64) (config T, err error) {
	err = global.DB.Model(data).Where("driver_id = ?", driverId).First(&config).Error
	return
}

// CreateDriverConfig 创建司机配置
// 参数:
//   - data: 配置数据模型指针
// 返回:
//   - error: 创建失败时返回错误信息
func CreateDriverConfig[T any](data *T) error {
	return global.DB.Create(data).Error
}

// UpdateDriverConfig 更新司机配置
// 参数:
//   - data: 配置数据模型指针（包含要更新的数据）
// 返回:
//   - error: 更新失败时返回错误信息
func UpdateDriverConfig[T any](data *T) error {
	return global.DB.Save(data).Error
}

// UpdateDriverConfigFields 更新司机配置的指定字段
// 参数:
//   - driverId: 司机ID
//   - updates: 需要更新的字段map，key为字段名，value为新值
// 返回:
//   - error: 更新失败时返回错误信息
func UpdateDriverConfigFields[T any](data *T, driverId int64, updates map[string]interface{}) error {
	return global.DB.Model(data).Where("driver_id = ?", driverId).Updates(updates).Error
}
