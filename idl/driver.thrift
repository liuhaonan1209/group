namespace go car.driver

struct DriverDetailReq{
    1: i64 id
}

struct DriverDetailResp {
    1: string name
    2: string tel
    3: string IDCard
    4: string license
    5: string registerDate
    6: string rating
}
service DriverService {
    DriverDetailResp DriverDetail(1:DriverDetailReq req)
}