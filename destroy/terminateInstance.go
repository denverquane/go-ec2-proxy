package destroy

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/denverquane/go-ec2-proxy/common"
	"log"
	"time"
)

func TerminateInstance(region common.Region, instanceId string) {
	ec2Creds, err := common.GetCredentialsFromEnvironment("EC2_ACCESS_KEY_ID", "EC2_SECRET_ACCESS_KEY")
	if err != nil {
		log.Println(err)
	}
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(string(region)),
		Credentials: ec2Creds,
	},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	input := ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			&instanceId,
		},
	}

	_, err = svc.TerminateInstances(&input)
	if err != nil {
		log.Println(err)
	} else {
		params := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				&instanceId,
			},
		}

		status := "shutting-down"

		for status == "shutting-down" {
			result, _ := svc.DescribeInstances(params)

			status = *result.Reservations[0].Instances[0].State.Name
			fmt.Println("Proxy isn't destroyed yet, sleeping for 5 seconds...")
			time.Sleep(5 * time.Second)
		}
		log.Println("Server " + instanceId + " terminated")
	}
}
