package metrics

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/denverquane/go-ec2-proxy/common"
	"log"
	"time"
)

func ListNetworkThroughputForInstance(creds *credentials.Credentials, region common.Region, instanceID string) (float64, float64) {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(string(region)),
		Credentials: creds,
	},
	)

	svc := cloudwatch.New(sess)
	now := time.Now()
	start := now.AddDate(0, 0, -1)
	var avg = "Sum"
	var period int64 = 60
	input := cloudwatch.GetMetricStatisticsInput{
		EndTime:    &now,
		StartTime:  &start,
		MetricName: aws.String("NetworkIn"),
		Namespace:  aws.String("AWS/EC2"),
		Statistics: []*string{
			&avg,
		},
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instanceID),
			},
		},
		Period: &period,
	}

	result, err := svc.GetMetricStatistics(&input)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("NetworkIn", result.Datapoints)
	totalIn := sumDatapointsByUnit(result.Datapoints, "Bytes")

	input.SetMetricName("NetworkOut")

	result, err = svc.GetMetricStatistics(&input)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("NetworkOut", result.Datapoints)

	totalOut := sumDatapointsByUnit(result.Datapoints, "Bytes")
	return totalIn, totalOut
}

func sumDatapointsByUnit(data []*cloudwatch.Datapoint, unit string) float64 {
	sum := 0.0
	for _, v := range data {
		if *v.Unit == unit {
			sum += *v.Sum
		}
	}
	return sum
}
