package create

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

func CreateAndStartProxyServer(creds *credentials.Credentials, proxyConfig common.ProxyConfig, serverConfig common.ServerConfig) (*ec2.Instance, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(string(serverConfig.Region)),
		Credentials: creds,
	},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	script := CreateGoProxyScriptString(proxyConfig)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(serverConfig.GetAmi()),
		InstanceType: aws.String(string(serverConfig.Tier)),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		UserData:     &script,
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		return &ec2.Instance{}, err
	}

	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)

	handle := common.GenerateUniqueInstanceHandle(proxyConfig, serverConfig)

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("Goproxy Node " + handle),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return &ec2.Instance{}, errtag
	}

	fmt.Println("Successfully tagged instance")

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			runResult.Instances[0].InstanceId,
		},
	}

	status := "pending"
	var result *ec2.DescribeInstancesOutput

	for status == "pending" {
		result, err = svc.DescribeInstances(params)
		if err != nil {
			log.Println(err)
		}

		// TODO Could cause an index panic?
		res := result.Reservations[0]
		inst := res.Instances[0]
		status = *inst.State.Name
		fmt.Println("Proxy isn't running yet, sleeping for 2 seconds...")
		time.Sleep(time.Second * 2)
	}

	return result.Reservations[0].Instances[0], err
}
