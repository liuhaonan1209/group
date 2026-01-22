namespace go car.route

// ========== 站点信息 ==========
struct StationInfo {
    1: i64 stationId       // 站点ID
    2: string name         // 站点名称
    3: i32 travelTime      // 从起点到此站点的行驶时间（分钟）
    4: i32 stopOrder       // 停靠顺序
}

// ========== 班线基础请求/响应 ==========
struct RouteReq {
    1: i64 RouteId
}

struct RouteResp {
    1: bool Success
    2: string Msg
    3: map<string, string> RouteInfo
}

// ========== 班线创建 ==========
struct RouteCreateReq {
    1: string routeNo           // 班线编号
    2: string fieet             // 车队
    3: string tage              // 标签
    4: i64 startStationId       // 起点站ID
    5: i64 endStationId         // 终点站ID
    6: list<StationInfo> upStations    // 上车站点列表
    7: list<StationInfo> downStations  // 下车站点列表
}

struct RouteCreateResp {
    1: bool Success
    2: string Msg
    3: i64 RouteId
}

// ========== 班线列表 ==========
struct RouteListReq {
    1: string routeNo
    2: string fieet
    3: string tage
    4: i32 page
    5: i32 size
}

struct RouteItem {
    1: i64 id
    2: string routeNo
    3: string fieet
    4: string tage
    5: string startStationName
    6: string endStationName
    7: bool isActive
}

struct RouteListResp {
    1: bool Success
    2: string Msg
    3: list<RouteItem> list
    4: i64 total
}

// ========== 班线详情 ==========
struct RouteInfoReq {
    1: i64 id
}

struct RouteInfoResp {
    1: bool Success
    2: string Msg
    3: i64 id
    4: string routeNo
    5: string fieet
    6: string tage
    7: i64 startStationId
    8: i64 endStationId
    9: list<StationInfo> upStations
    10: list<StationInfo> downStations
    11: bool isActive
}

// ========== 班线更新 ==========
struct RouteUpdateReq {
    1: i64 id
    2: string routeNo
    3: string fieet
    4: string tage
    5: i64 startStationId
    6: i64 endStationId
    7: list<StationInfo> upStations
    8: list<StationInfo> downStations
    9: bool isActive
}

struct RouteUpdateResp {
    1: bool Success
    2: string Msg
}

// ========== 班线删除 ==========
struct RouteDeleteReq {
    1: i64 id
}

struct RouteDeleteResp {
    1: bool Success
    2: string Msg
}

// ========== 添加站点到班线 ==========
struct AddStationReq {
    1: i64 routeId          // 班线ID
    2: i64 stationId        // 站点ID
    3: i32 travelTime       // 从起点的行驶时间（分钟）
    4: i32 stopType         // 站点类型：1-上车站点，2-下车站点
}

struct AddStationResp {
    1: bool Success
    2: string Msg
}

// ========== 移除站点 ==========
struct RemoveStationReq {
    1: i64 routeId
    2: i64 stationId
}

struct RemoveStationResp {
    1: bool Success
    2: string Msg
}

// ========== 班次管理 ==========
struct ScheduleInfo {
    1: i64 id
    2: i64 routeId
    3: string departureTime   // 发车时间 HH:mm
    4: string arrivalTime     // 预计到达时间 HH:mm
    5: i32 capacity           // 座位容量
    6: bool isActive
}

struct ScheduleCreateReq {
    1: i64 routeId
    2: string departureTime
    3: string arrivalTime
    4: i32 capacity
}

struct ScheduleCreateResp {
    1: bool Success
    2: string Msg
    3: i64 scheduleId
}

struct ScheduleListReq {
    1: i64 routeId
    2: i32 page
    3: i32 size
}

struct ScheduleListResp {
    1: bool Success
    2: string Msg
    3: list<ScheduleInfo> list
    4: i64 total
}

struct ScheduleUpdateReq {
    1: i64 id
    2: string departureTime
    3: string arrivalTime
    4: i32 capacity
    5: bool isActive
}

struct ScheduleUpdateResp {
    1: bool Success
    2: string Msg
}

struct ScheduleDeleteReq {
    1: i64 id
}

struct ScheduleDeleteResp {
    1: bool Success
    2: string Msg
}

// 微服务接口
service RouteService {
    // 班线管理
    RouteResp GetRoute(1: RouteReq req)
    RouteCreateResp RouteCreate(1: RouteCreateReq req)
    RouteListResp RouteList(1: RouteListReq req)
    RouteInfoResp RouteInfo(1: RouteInfoReq req)
    RouteUpdateResp RouteUpdate(1: RouteUpdateReq req)
    RouteDeleteResp RouteDelete(1: RouteDeleteReq req)
    
    // 站点管理
    AddStationResp AddStation(1: AddStationReq req)
    RemoveStationResp RemoveStation(1: RemoveStationReq req)
    
    // 班次管理
    ScheduleCreateResp ScheduleCreate(1: ScheduleCreateReq req)
    ScheduleListResp ScheduleList(1: ScheduleListReq req)
    ScheduleUpdateResp ScheduleUpdate(1: ScheduleUpdateReq req)
    ScheduleDeleteResp ScheduleDelete(1: ScheduleDeleteReq req)
}
