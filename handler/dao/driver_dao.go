package dao

import (
	"group/global"
)

// GetDriverByID 根据ID查询司机详情
func GetDriverByID[T any](data *T, id int64) (driverDetail T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).Find(&driverDetail).Error
	return
}

// CreateDriver 创建司机
func CreateDriver[T any](data *T) error {
	return global.DB.Create(data).Error
}

// UpdateDriver 更新司机信息
func UpdateDriver[T any](data *T) error {
	return global.DB.Save(data).Error
}

// DeleteDriver 删除司机
func DeleteDriver[T any](data *T, id int64) error {
	return global.DB.Delete(data, id).Error
}
