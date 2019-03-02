package create

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/denverquane/go-ec2-proxy/destroy"
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

	proxyConfig := common.ProxyConfig{"http", "23455", "", ""}

	serverConfig := common.CreateServerConfig(common.USWest1, common.Micro)

	instance, err := CreateAndStartProxyServer(creds, proxyConfig, serverConfig)

	log.Println("Server created, IP: " + *instance.PublicIpAddress + ":" + proxyConfig.Port)

	destroy.TerminateInstance(creds, serverConfig, *instance.InstanceId)
}
