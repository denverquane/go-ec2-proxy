package destroy

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/denverquane/go-ec2-proxy/common"
	"log"
	"time"
)

func TerminateInstance(creds *credentials.Credentials, serverConfig common.ServerConfig, instance string) {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(string(serverConfig.Region)),
		Credentials: creds,
	},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	input := ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			&instance,
		},
	}

	_, err := svc.TerminateInstances(&input)
	if err != nil {
		log.Println(err)
	} else {
		params := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				&instance,
			},
		}

		status := "shutting-down"

		for status == "shutting-down" {
			result, _ := svc.DescribeInstances(params)

			status = *result.Reservations[0].Instances[0].State.Name
			fmt.Println("Proxy isn't destroyed yet, sleeping for 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}
