package main

import (
	"fmt"
	"github.com/QuanAnh/rabbitmq-go/endpoints"
	"github.com/QuanAnh/rabbitmq-go/service"
	"github.com/QuanAnh/rabbitmq-go/transport"
	"github.com/ThomasVNN/go-base/storage/postgresql"
	"github.com/go-kit/kit/log"
	"os"
)
func main() {

	// MAKE DB CONNECTION
	conn, err := postgresql.GetConnection()
	if err != nil {
		fmt.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	repository := service.NewRepo(conn)
	addService := service.NewBasicService(repository, log.NewNopLogger())
	addEndpoints := endpoints.MakeEndpoints(addService)
	url := "amqp://guest:guest@localhost:5672/"
	err = transport.MakeAMQPHandle(addEndpoints, url)
	if err != nil {
		fmt.Errorf("Error", err)
	}
}