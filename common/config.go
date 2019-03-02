package common

import (
	"crypto/sha256"
	"encoding/hex"
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

func (sc ServerConfig) GetAmi() string {
	return sc.ami
}

func CreateServerConfig(region Region, tier Tier) ServerConfig {
	//TODO this AMI is a hefty assumption for always being us-west-1!!!
	return ServerConfig{region, "ami-8d948ced", tier}
}

func GenerateUniqueInstanceHandle(config ProxyConfig, serverConfig ServerConfig) string {
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
