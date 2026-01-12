package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	passenger "group/kitex_gen/car/passenger"
	"log"
)

// PassengerServiceImpl implements the last service interface defined in the IDL.
type PassengerServiceImpl struct{}

// PassengerDetail implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) PassengerDetail(ctx context.Context, req *passenger.PassengerDetailReq) (resp *passenger.PassengerDetailResp, err error) {
	log.Printf("查询乘客ID: %d", req.Id)

	var p model.Passenger
	id, err := dao.GetPassengerByID(&p, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询乘客失败")
	}

	// 构造响应
	resp = &passenger.PassengerDetailResp{
		Name:         id.Name,
		Tel:          id.Tel,
		IDCard:       id.IDCard,
		RegisterDate: id.RegisterDate.Format("2006-01-02"),
	}

	log.Printf("查询乘客信息成功: %d", req.Id)
	return resp, nil
}
