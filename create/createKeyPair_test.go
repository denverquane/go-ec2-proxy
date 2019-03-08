package create

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestCreateKeyPair(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	creds := credentials.NewStaticCredentials(os.Getenv("EC2_ACCESS_KEY_ID"), os.Getenv("EC2_SECRET_ACCESS_KEY"), "")

	CreateKeyPair(creds, common.USWest1)
}
