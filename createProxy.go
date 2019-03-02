package go_ec2_proxy

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"time"
)

type ProxyConfig struct {
	Type string
	Port string
	User string
	Pass string
}

// idiomatic "enum" for acceptable tiers
type Tier string

const (
	Nano  Tier = "t2.nano"
	Micro Tier = "t2.micro"
)

// valid regions
type Region string

const (
	USWest1 Region = "us-west-1"
)

type ServerConfig struct {
	Region Region
	ami    string
	Tier   Tier
}

func CreateServerConfig(region Region, tier Tier) ServerConfig {
	//TODO this AMI is a hefty assumption for always being us-west-1!!!
	return ServerConfig{region, "ami-8d948ced", tier}
}

func generateUniqueInstanceHandle(config ProxyConfig, serverConfig ServerConfig) string {
	h := sha256.New()
	h.Write([]byte(time.Now().String()))

	h.Write([]byte(config.Type))
	h.Write([]byte(config.Port))
	h.Write([]byte(config.User))
	h.Write([]byte(config.Pass))

	h.Write([]byte(serverConfig.Region))
	h.Write([]byte(serverConfig.ami))
	h.Write([]byte(serverConfig.Tier))

	return hex.EncodeToString(h.Sum(nil))
}

func CreateAndStartProxyServer(creds *credentials.Credentials, proxyConfig ProxyConfig, serverConfig ServerConfig) (*ec2.Instance, error) {
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
		ImageId:      aws.String(serverConfig.ami),
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

	handle := generateUniqueInstanceHandle(proxyConfig, serverConfig)

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

func TerminateInstance(creds *credentials.Credentials, serverConfig ServerConfig, instance string) {
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
