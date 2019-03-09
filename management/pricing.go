package management

import "github.com/denverquane/go-ec2-proxy/common"

const BaseCharge = 100.0 //1 dollar

func GetTotalPrice(region common.Region, tier common.Tier, durationInHrs int, quantity int) float64 {
	total := BaseCharge

	switch region {
	case common.USWest1:
		break
	default:
		total *= 1
		break
	}

	switch tier {
	case common.Micro:
		total *= 1.2

	case common.Nano:
	default:
		break
	}

	total *= float64(durationInHrs) / 10.0

	return total * float64(quantity)
}
