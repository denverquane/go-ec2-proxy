package create

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DoesSGExistAlready(creds *credentials.Credentials, region, sgName string) bool {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	},
	)

	// Create EC2 service client
	_ := ec2.New(sess)

	return false
}
