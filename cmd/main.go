package main

import (
	"fmt"
	"github.com/varunturlapati/simpleWebSvc/web"
	
	"github.com/varunturlapati/simpleWebSvc/pkg/db/redisclient"
)

func main() {
	rClient, _ := redisclient.New()
	pong, err := rClient.PingPong()
	fmt.Printf("%v and %v\n", pong, err)
	fmt.Println("Rest API v2.0 - Mux Routers")
	web.HandleRequests()

}
