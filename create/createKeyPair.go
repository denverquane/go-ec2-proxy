package create

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/denverquane/go-ec2-proxy/common"
)

func CreateKeyPair(creds *credentials.Credentials, region common.Region, name string) {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(string(region)),
		Credentials: creds,
	},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	input := ec2.CreateKeyPairInput{
		KeyName: aws.String(name),
	}

	out, _ := svc.CreateKeyPair(&input)
	fmt.Println("Made keypair")
	fmt.Println(out)
}
