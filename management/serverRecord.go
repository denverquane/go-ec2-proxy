package management

import (
	"fmt"
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/denverquane/go-ec2-proxy/create"
	"github.com/denverquane/go-ec2-proxy/destroy"
	"github.com/denverquane/go-ec2-proxy/metrics"
	"log"
	"strconv"
	"time"
)

const KbInGb = 1000000000.0
const CloudwatchRefreshMinutesInterval = 5

type ServerRecord struct {
	InstanceId string

	Region     common.Region
	User       string
	Pass       string
	PublicIp   string
	PublicPort string

	//private elements; shouldn't modify these directly
	constraints ServerConstraints
	status      ServerStatus
}

func (sr *ServerRecord) FetchCurrentDataUsage() {
	creds, err := common.GetCredentialsFromEnvironment("CLOUDWATCH_ACCESS_KEY_ID", "CLOUDWATCH_SECRET_ACCESS_KEY")
	if err != nil {
		log.Print(err)
	}

	in, out := metrics.FetchNetworkThroughputForInstance(creds, sr.Region, sr.InstanceId)
	sr.status.InboundBytesUsed = in
	sr.status.OutboundBytesUsed = out

	fmt.Println(sr.InstanceId + " Gb in: " + strconv.FormatFloat(sr.status.InboundBytesUsed/KbInGb, 'f', -1, 64))
	fmt.Println(sr.InstanceId + " Gb out: " + strconv.FormatFloat(sr.status.OutboundBytesUsed/KbInGb, 'f', -1, 64))
}

func StartProxyAndReturnRecord(pc common.ProxyConfig, sc common.ServerConfig, duration time.Duration, byteCap float64) (*ServerRecord, error) {
	ec2Creds, err := common.GetCredentialsFromEnvironment("EC2_ACCESS_KEY_ID", "EC2_SECRET_ACCESS_KEY")
	if err != nil {
		log.Fatal(err)
	}

	sgId := create.FindOrCreateSecurityGroup(ec2Creds, common.USWest1, pc.Port)

	instance, err := create.CreateAndStartProxyServer(ec2Creds, pc, sc, sgId)
	if err != nil {
		log.Println(err)

		return nil, err
	}

	log.Println("Proxy instance " + *instance.InstanceId + " is running, IP: " + *instance.PublicIpAddress + ":" + pc.Port)

	constraints := ServerConstraints{time.Now().Add(duration), byteCap}
	status := ServerStatus{time.Now(), 0, 0, true, false}

	record := ServerRecord{*instance.InstanceId, sc.Region, "", "", *instance.PublicIpAddress, pc.Port, constraints, status}

	go monitorUsageAsync(&record)

	return &record, err
}

func monitorUsageAsync(record *ServerRecord) {
	kill := false

	for !kill {
		time.Sleep(time.Minute * CloudwatchRefreshMinutesInterval)
		record.FetchCurrentDataUsage()

		if record.constraints.DestructionTime.Before(time.Now()) {
			fmt.Println("Time ran out for the server, killing!")
			kill = true
		} else if record.constraints.TotalByteCap < (record.status.OutboundBytesUsed + record.status.InboundBytesUsed) {
			fmt.Println("Server has used too much data, killing!")
			kill = true
		}
	}

	destroy.TerminateInstance(record.Region, record.InstanceId)
	record.status.IsDestroyed = true
	record.status.IsRunning = false
}

type ServerConstraints struct {
	DestructionTime time.Time
	TotalByteCap    float64
}

type ServerStatus struct {
	CreationTime time.Time

	InboundBytesUsed  float64
	OutboundBytesUsed float64

	IsRunning   bool
	IsDestroyed bool
}
