package main

import (
	"fmt"
	"log"

	"github.com/takama/daemon"
)

func main() {
	service, err := daemon.New("com.yulefox.name", "description")

	log.Printf("%+v", service)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	usage := "Usage: "
	status, err := service.Install()
	if err != nil {
		log.Fatal(status, "\nError: ", err)
	}
	fmt.Println(status)
}
