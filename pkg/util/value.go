package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ReplaceVariables 取代變數
func ReplaceVariables(reply string, variables map[string]interface{}) string {
	r, _ := regexp.Compile("\\$\\{(.*?)\\}")

	for _, match := range r.FindAllStringSubmatch(reply, -1) {
		var (
			handleURLEncode = false
		)
		if strings.HasPrefix(match[1], "##URL##:") {
			match[1] = match[1][8:len(match[1])]
			handleURLEncode = true
		}

		if v := GetJSONValue(match[1], variables); v != nil {
			switch val := v.(type) {
			case string:
				if handleURLEncode {
					val = url.QueryEscape(val)
				}
				reply = strings.Replace(reply, match[0], val, -1)
			case bool:
				reply = strings.Replace(reply, match[0], fmt.Sprintf("%t", val), -1)
			case interface{}:
				reply = strings.Replace(reply, match[0], fmt.Sprintf("%v", val), -1)
			}
		}
	}

	return reply
}

// GetJSONValue by Layer
func GetJSONValue(path string, data map[string]interface{}) interface{} {
	layer := strings.Split(path, ".")
	return getValue(layer, data)
}

func getValue(layer []string, data interface{}) interface{} {
	if len(layer) <= 1 || data == nil {
		return data
	}

	switch data.(type) {
	case map[string]interface{}:
		data = (data.(map[string]interface{}))[layer[0]]
	default:
		return data
	}

	return getValue(layer[1:len(layer)], data)
}
