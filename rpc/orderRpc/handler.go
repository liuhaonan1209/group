package main

import (
	"context"
	order "group/kitex_gen/car/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// OrderDetail implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) OrderDetail(ctx context.Context, req *order.OrderDetailReq) (resp *order.OrderDetailResp, err error) {
	// TODO: Your code here...
	return
}
