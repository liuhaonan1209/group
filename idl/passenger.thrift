namespace go car.passenger

struct PassengerDetailReq{
    1: i64 id
}

struct PassengerDetailResp {
    1: string name
    2: string tel
    3: string IDCard
    4: string registerDate
}
service PassengerService {
    PassengerDetailResp PassengerDetail(1:PassengerDetailReq req)
}