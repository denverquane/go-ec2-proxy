package main

import (
	"fmt"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/denverquane/go-ec2-proxy/create"
	"github.com/denverquane/go-ec2-proxy/metrics"
	"github.com/joho/godotenv"
	"log"
	"strconv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cloudwatchCreds, err := common.GetCredentialsFromEnvironment("CLOUDWATCH_ACCESS_KEY_ID", "CLOUDWATCH_SECRET_ACCESS_KEY")
	ec2Creds, err := common.GetCredentialsFromEnvironment("EC2_ACCESS_KEY_ID", "EC2_SECRET_ACCESS_KEY")

	proxyConfig := common.ProxyConfig{"http", "23455", "", ""}

	serverConfig := common.CreateServerConfig(common.USWest1, common.Micro)

	sgId := create.FindOrCreateSecurityGroup(ec2Creds, common.USWest1, proxyConfig.Port)

	instance, err := create.CreateAndStartProxyServer(ec2Creds, proxyConfig, serverConfig, sgId)

	log.Println("Instance ID: " + *instance.InstanceId)

	log.Println("Server created, IP: " + *instance.PublicIpAddress + ":" + proxyConfig.Port)

	in, out := metrics.FetchNetworkThroughputForInstance(cloudwatchCreds, common.USWest1, *instance.InstanceId)
	fmt.Println("In: " + strconv.FormatFloat(in, 'f', -1, 64) + " , Out: " + strconv.FormatFloat(out, 'f', -1, 64))
}
