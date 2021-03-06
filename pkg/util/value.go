package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/lineemoji"
)

// ExcuteGoTemplate parse and excute template
func ExcuteGoTemplate(reply string, variables map[string]interface{}) string {
	tmpl := template.New("temp")
	tmpl.Parse(reply)
	var (
		data       bytes.Buffer
		dataString string
	)
	err := tmpl.Execute(&data, variables)
	if err != nil {
		logger.Error("[pkg] ", "Excute Template Error: ", err)
	}
	dataString, err = strconv.Unquote(data.String())
	if err != nil {
		logger.Error("[pkg] ", "Strings Unquote Error: ", err)
		dataString = data.String()
	}
	dataString = html.UnescapeString(dataString)
	return dataString
}

// ReplaceVariables from string
func ReplaceVariables(reply string, variables map[string]interface{}) string {
	r, _ := regexp.Compile("\\$\\{(.*?)\\}")

	for _, match := range r.FindAllStringSubmatch(reply, -1) {
		var (
			handleURLEncode = false
			value           string
		)
		if strings.HasPrefix(match[1], "##URL##:") {
			match[1] = match[1][8:len(match[1])]
			handleURLEncode = true
		}

		// find variable
		if v := GetJSONValue(match[1], variables); v != nil {
			value = fmt.Sprintf("%v", v)
		} else {
			value = match[1]
		}
		// URLEncode
		if handleURLEncode {
			value = url.QueryEscape(value)
		}
		// replace
		reply = strings.Replace(reply, match[0], value, -1)
	}

	return reply
}

// ReplaceLineEmoji replace line custom emoji character
func ReplaceLineEmoji(reply string) string {
	r, _ := regexp.Compile("(\\(.*?\\))")
	for _, match := range r.FindAllStringSubmatch(reply, -1) {
		if emoji, err := lineemoji.GetEmoji(match[0]); err == nil {
			reply = strings.Replace(reply, match[0], string(emoji), -1)
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
	if len(layer) <= 0 || data == nil {
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

// GetJSONBytes Unmarshal json to []byte
func GetJSONBytes(input interface{}) []byte {
	param, _ := json.Marshal(input)
	return param
}

// Itob int to []byte
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// ReflectFieldValue reflect struct field
func ReflectFieldValue(data interface{}, fieldName string) reflect.Value {
	var ptr reflect.Value
	var value reflect.Value
	value = reflect.ValueOf(data)

	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(data))
		temp := ptr.Elem()
		temp.Set(value)
		value = temp
	}

	return value.FieldByName(fieldName)
}
