package go_ec2_proxy

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestCreateProxy(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	creds := credentials.NewStaticCredentials(os.Getenv("EC2_ACCESS_KEY_ID"), os.Getenv("EC2_SECRET_ACCESS_KEY"), "")

	proxyConfig := ProxyConfig{"http", "23455", "", ""}

	serverConfig := CreateServerConfig(USWest1, Micro)

	instance, err := CreateAndStartProxyServer(creds, proxyConfig, serverConfig)

	log.Println("Server created, IP: " + *instance.PublicIpAddress + ":" + proxyConfig.Port)

	TerminateInstance(creds, serverConfig, *instance.InstanceId)
}
