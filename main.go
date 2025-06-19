package main

import (
	controller "bimage/controller/http"
	codeprocessor "bimage/infrastructure/consumer"
	"bimage/usecase/service"
	"fmt"
)

func main() {
	consumer := codeprocessor.New()
	usecase := service.New(consumer)
	server := controller.New(usecase)

	r := server.WithObjectHandler()

	fmt.Println("gol")

	if err := controller.CreateAndRunServer(":8080", r); err != nil {
		panic(err)
	} else {
		fmt.Println("start")
	}

}
