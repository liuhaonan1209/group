package main

import (
	order "group/kitex_gen/car/order/orderservice"
	"log"
)

func main() {
	svr := order.NewServer(new(OrderServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
