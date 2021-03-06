package model

import (
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/j75689/easybot/pkg/logger"
)

// Iptable accessible ips
type Iptable struct {
	ID    string   `json:"id" bson:"_id"`
	Type  string   `json:"type" bson:"type"` // allow, deny
	IP    []string `json:"ip" bson:"ip"`
	Scope string   `json:"scope" bson:"scope"`
}

// Pass check clientIP
func (iptable *Iptable) Pass(clientIP string) bool {

	var preset bool
	switch iptable.Type {
	case "allow":
		preset = true
	case "deny":
		preset = false
	}

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
			return preset
		}
		if reflect.DeepEqual(client, ipfilter) {
			return preset
		}
	}
	return !preset
}
