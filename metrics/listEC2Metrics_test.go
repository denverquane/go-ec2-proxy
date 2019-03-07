package metrics

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestListNetworkThroughputForInstance(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	creds := credentials.NewStaticCredentials(os.Getenv("CLOUDWATCH_ACCESS_KEY_ID"), os.Getenv("CLOUDWATCH_SECRET_ACCESS_KEY"), "")

	in, out := ListNetworkThroughputForInstance(creds, common.USWest1, "i-0cc634284080bf8d7")
	fmt.Println(strconv.FormatFloat(in/1000000000.0, 'f', -1, 64) + "Gb in")
	fmt.Println(strconv.FormatFloat(out/1000000000.0, 'f', -1, 64) + "Gb out")
}
