package main

import (
	"fmt"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/denverquane/go-ec2-proxy/management"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	for i := 23000; i < 23050; i++ {
		port := strconv.Itoa(i)
		proxyConfig := common.ProxyConfig{"http", port, "", ""}

		serverConfig := common.CreateServerConfig(common.USWest1, common.Micro)

		go management.StartProxyAndReturnRecord(proxyConfig, serverConfig, time.Minute*10, 1000)
		fmt.Println("Sleeping for 5 seconds before starting next server...")
		time.Sleep(time.Second * 5)
	}

	for true {
		fmt.Println("Sleepy")
		time.Sleep(time.Minute * 5)
	}
}
