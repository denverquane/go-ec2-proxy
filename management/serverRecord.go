package management

import (
	"github.com/denverquane/go-ec2-proxy/common"
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
