package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	"group/kitex_gen/car/starstation"
)

// StarStationServiceImpl implements the last service interface defined in the IDL.
type StarStationServiceImpl struct{}

// GetStarStation implements the StarStationServiceImpl interface.
func (s *StarStationServiceImpl) GetStarStation(ctx context.Context, req *starstation.StarStationReq) (resp *starstation.StarStationResp, err error) {
	// TODO: Your code here...
	resp = &starstation.StarStationResp{}

	// Call the DAO method to get station information
	var stationInfo struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	// Use RouteInfo method from dao to get station details
	err = dao.RouteInfo(&stationInfo, int(req.StationId))
	if err != nil {
		resp.Success = false
		resp.Msg = "Failed to get station info: " + err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "Station retrieved successfully"
	resp.StationName = stationInfo.Name

	return
}

func (s *StarStationServiceImpl) CreateStarStation(ctx context.Context, req *starstation.CreateStarStationReq) (resp *starstation.CreateStarStationResp, err error) {

	resp = &starstation.CreateStarStationResp{}

	starstationModel := &model.StartStation{
		Name:    req.Name,
		Address: req.Address,
	}

	err = dao.StartStationCreate(starstationModel)
	if err != nil {
		return nil, fmt.Errorf("班点添加失败")
	}

	return resp, nil
}

func (s *StarStationServiceImpl) RemoveDuplicateStops(ctx context.Context, req *starstation.RemoveDuplicateStopsReq) (resp *starstation.RemoveDuplicateStopsResp, err error) {
	resp = &starstation.RemoveDuplicateStopsResp{}

	route := model.Route{}

	err = dao.ValidateBusLineStops(&route)
	if err != nil {
		return nil, err
	}

	time, err := dao.SortStopsByTime([]model.StartStation{})
	if err != nil {
		return nil, err
	}

	fmt.Println(time)
	return resp, nil
}

func (s *StarStationServiceImpl) StarStationList(ctx context.Context, req *starstation.StarStationListReq) (resp *starstation.StarStationListResp, err error) {

	resp = &starstation.StarStationListResp{}
	return resp, nil
}

//func (s *StarStationServiceImpl) StartStationInfo(ctx context.Context, req *starstation.StarStationReq) (resp *starstation.StarStationResp, err error) {
//	// TODO: Your code here...
//	resp = &starstation.StarStationResp{}
//	var Info struct {
//		ID      int64  `json:"ID"`
//		Name    string `json:"name"`
//		Address string `json:"Address"`
//	}
//	err = dao.RouteInfo(&Info, int(req.StationId))
//	if err != nil {
//		resp.Success = false
//		resp.Msg = "Failed to get station info: " + err.Error()
//		return resp, nil
//	}
//	resp.Success = true
//	resp.Msg = "站点查看详情成功"
//	//resp.StationName =
//	return
//}
//
//func (s *StarStationServiceImpl) CreateStarStation(ctx context.Context, req *starstation.StarStationReq) (resp *starstation.StarStationResp, err error) {
//	// TODO: Your code here...
//
//	resp = &starstation.StarStationResp{}
//
//	err = dao.StartStationCreate(resp)
//
//	if err != nil {
//		return nil, err
//	}
//	return
//}
