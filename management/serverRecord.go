package management

import (
	"github.com/denverquane/go-ec2-proxy/common"
	"github.com/denverquane/go-ec2-proxy/metrics"
	"log"
	"time"
)

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

func FetchCurrentDataUsage(sr *ServerRecord) {
	creds, err := common.GetCredentialsFromEnvironment("CLOUDWATCH_ACCESS_KEY_ID", "CLOUDWATCH_SECRET_ACCESS_KEY")
	if err != nil {
		log.Print(err)
	}

	in, out := metrics.FetchNetworkThroughputForInstance(creds, sr.Region, sr.InstanceId)
	sr.status.InboundBytesUsed = in
	sr.status.OutboundBytesUsed = out
}

type ServerConstraints struct {
	DestructionTime time.Time
	TotalByteCap    float64
}

type ServerStatus struct {
	CreationTime time.Time

	InboundBytesUsed  float64
	OutboundBytesUsed float64

	IsRunning     bool
	ShouldDestroy bool
	IsDestroyed   bool
}
