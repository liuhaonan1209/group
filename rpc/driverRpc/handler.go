package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	driver0 "group/kitex_gen/car/driver"
	"log"
)

// DriverServiceImpl implements the last service interface defined in the IDL.
type DriverServiceImpl struct{}

// DriverDetail implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) DriverDetail(ctx context.Context, req *driver0.DriverDetailReq) (resp *driver0.DriverDetailResp, err error) {

	// 查询司机信息
	var driver model.Driver
	id, err := dao.GetDriverByID(&driver, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询司机信息失败")
	}

	// 构造响应
	resp = &driver0.DriverDetailResp{
		Name:         id.Name,
		Tel:          id.Tel,
		IDCard:       id.IDCard,
		License:      id.License,
		RegisterDate: id.RegisterDate.Format("2006-01-02"),
		Rating:       fmt.Sprintf("%.2f", id.Rating),
	}

	log.Printf("查询司机信息成功: %d", req.Id)
	return resp, nil
}
