package create

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestCreateSecurityGroup(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	creds := credentials.NewStaticCredentials(os.Getenv("EC2_ACCESS_KEY_ID"), os.Getenv("EC2_SECRET_ACCESS_KEY"), "")

	port := "43534"

	groupID := FindOrCreateSecurityGroup(creds, common.USWest1, port)

	fmt.Println("Group id for port " + port + ": " + groupID)
}
