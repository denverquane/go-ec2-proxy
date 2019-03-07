package metrics

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestListNetworkThroughputForInstance(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	creds := credentials.NewStaticCredentials(os.Getenv("CLOUDWATCH_ACCESS_KEY_ID"), os.Getenv("CLOUDWATCH_SECRET_ACCESS_KEY"), "")

	ListNetworkThroughputForInstance(creds, common.USWest1, "i-0fe9e2bdf13375bf2")
}
