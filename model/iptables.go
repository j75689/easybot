package model

import (
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/j75689/easybot/pkg/logger"
)

// Iptables accessible ips
type Iptables struct {
	Type  string   `json:"type" bson:"type"` // allow, deny
	IP    []string `json:"ip" bson:"ip"`
	Scope string   `json:"scope" bson:"scope"`
}

// Pass check clientIP
func (iptable *Iptables) Pass(clientIP string) bool {
	for _, ip := range iptable.IP {
		_, ipfilter, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Error("[iptables] ", err)
			continue
		}

		mask := strings.Split(ip, "/")[1]

		_, client, err := net.ParseCIDR(fmt.Sprintf("%s/%s", clientIP, mask))
		if err != nil {
			logger.Error("[iptables] ", err)
			return false
		}
		if reflect.DeepEqual(client, ipfilter) {
			return true
		}
	}
	return false
}
