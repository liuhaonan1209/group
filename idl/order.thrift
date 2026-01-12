namespace go car.order

struct OrderDetailReq{
    1: i64 id
}

struct OrderDetailResp {
    1: string name
    2: string tel
    3: string IDCard
    4: string license
    5: string registerDate
    6: string rating
}
service OrderService {
    OrderDetailResp OrderDetail(1:OrderDetailReq req)
}