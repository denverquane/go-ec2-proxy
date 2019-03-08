package create

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/denverquane/go-ec2-proxy/common"
	"log"
	"strconv"
)

func FindOrCreateSecurityGroup(creds *credentials.Credentials, region common.Region, port string) string {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(string(region)),
		Credentials: creds,
	},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	sgName := createSGNameFromPort(port)

	exists, result := doesSGExistAlready(svc, sgName)

	if exists {
		log.Println("Security group already exists, id: " + *result.SecurityGroups[0].GroupId)
		return *result.SecurityGroups[0].GroupId
	} else {
		log.Println("Security group does not exist for port: " + port + ", creating a new one...")
		input := &ec2.CreateSecurityGroupInput{
			Description: aws.String("Proxy SG for port " + port + " (Managed by go-ec2-proxy)"),
			GroupName:   aws.String(sgName),
		}
		createResult, err := svc.CreateSecurityGroup(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return ""
		} else {
			log.Println("Created security group: " + *createResult.GroupId)
			portInt, strerr := strconv.Atoi(port)
			if strerr != nil {
				log.Println(strerr)
				return ""
			}
			input := &ec2.AuthorizeSecurityGroupIngressInput{
				GroupId: aws.String(*createResult.GroupId),
				IpPermissions: []*ec2.IpPermission{
					{
						FromPort:   aws.Int64(22),
						IpProtocol: aws.String("tcp"),
						IpRanges: []*ec2.IpRange{
							{
								CidrIp:      aws.String("0.0.0.0/0"),
								Description: aws.String("SSH access"),
							},
						},
						ToPort: aws.Int64(22),
					},
					{
						FromPort:   aws.Int64(int64(portInt)),
						IpProtocol: aws.String("tcp"),
						IpRanges: []*ec2.IpRange{
							{
								CidrIp:      aws.String("0.0.0.0/0"),
								Description: aws.String("Proxy Port access"),
							},
						},
						ToPort: aws.Int64(int64(portInt)),
					},
				},
			}

			_, err := svc.AuthorizeSecurityGroupIngress(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					default:
						fmt.Println(aerr.Error())
					}
				} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
				}
				return ""
			} else {
				log.Println("Successfully added ingress rule for Security Group " + *createResult.GroupId)
			}

			return *createResult.GroupId
		}
	}
}

func createSGNameFromPort(port string) string {
	return "sg_proxies_port_" + port
}

func doesSGExistAlready(svc *ec2.EC2, sgName string) (bool, ec2.DescribeSecurityGroupsOutput) {
	input := &ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{
			aws.String(sgName),
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return false, *result
	}

	return true, *result
}
