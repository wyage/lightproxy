package main

import (
	"fmt"
	"log"
	"myproxy"
	"net/http"
)

func main() {
	startUp()
}

func startUp() {
	configfile := "config.json"
	config := myproxy.NewConfig(configfile)

	log.Println(config)

	handler := myproxy.NewHandler(config)
	fmt.Println("will listen at", config.ListenHost)

	err := http.ListenAndServe(config.ListenHost, handler)
	if err != nil {
		fmt.Println("startup failed:", err)
	}
}
