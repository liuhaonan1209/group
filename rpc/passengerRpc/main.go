package main

import (
	"group/core"
	passenger "group/kitex_gen/car/passenger/passengerservice"
	"log"
)

func main() {
	core.Nacos()
	core.Mysql()

	log.Println("Passenger service starting...")

	svr := passenger.NewServer(new(PassengerServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
