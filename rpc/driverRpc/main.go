package main

import (
	"group/core"
	driver "group/kitex_gen/car/driver/driverservice"
	"log"
)

func main() {
	core.Nacos()
	core.Mysql()

	svr := driver.NewServer(new(DriverServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
