package provider

import (
	"os"

	"github.com/ipinfo/go/v2/ipinfo"
)

var ipInfoAPIKey = os.Getenv("IP_INFO_API_KEY")

func NewIpInfoClient() *ipinfo.Client {
	client := ipinfo.NewClient(nil, nil, ipInfoAPIKey)

	return client
}
