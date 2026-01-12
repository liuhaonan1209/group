package dao

import (
	"group/global"
)

// GetDriverByID 根据ID查询乘客详情
func GetPassengerByID[T any](data *T, id int64) (passengerDetail T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).Find(&passengerDetail).Error
	return
}

// CreateDriver 创建乘客
func CreatePassenger[T any](data *T) error {
	return global.DB.Create(data).Error
}

// UpdateDriver 更新乘客信息
func UpdatePassenger[T any](data *T) error {
	return global.DB.Save(data).Error
}

// DeleteDriver 删除乘客
func DeletePassenger[T any](data *T, id int64) error {
	return global.DB.Delete(data, id).Error
}
