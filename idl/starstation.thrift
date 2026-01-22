namespace go car.starstation


// 1. 定义请求结构体
struct StarStationReq {
    1: i64 StationId  // 站点ID
}

// 2. 定义响应结构体
struct StarStationResp {
    1: bool Success
    2: string Msg
    3: string StationName  // 站点名称
}


// 1. 定义请求结构体
struct CreateStarStationReq {
    1: string Name  // 站点名称
    2: string Address //详细地址
}

// 2. 定义响应结构体
struct CreateStarStationResp {
    1: bool Success
    2: string Msg
    3: string StationName  // 站点名称
}

// 1. 定义请求结构体
struct RemoveDuplicateStopsReq {
    1: string Name  // 站点名称
    2: string Address //详细地址
}

// 2. 定义响应结构体
struct RemoveDuplicateStopsResp {
    1: bool Success
    2: string Msg
    3: string StationName  // 站点名称
}

struct StarStationListResp {
    1: i64 page
    2: i64 size
}

// 2. 定义响应结构体
struct StarStationListReq {
    1: bool Success
    2: string Msg
    3: map<string,string> StationName
}

// 3. 必须的service定义（Kitex核心依赖）
service StarStationService {
    StarStationResp GetStarStation(1: StarStationReq req)  // 接口方法
    CreateStarStationResp CreateStarStation(1: CreateStarStationReq req)
    RemoveDuplicateStopsResp RemoveDuplicateStops(1: RemoveDuplicateStopsReq req)//站点排重
    StarStationListResp StarStationList(1: StarStationListReq req)
}