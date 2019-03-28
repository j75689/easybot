package server

import (
	"github.com/j75689/easybot/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_replaceVariables(t *testing.T) {
	type args struct {
		reply     string
		variables map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Test_replaceVariables",
			args: args{
				reply: `{"name":"Test","message":"${test}"}`,
				variables: map[string]interface{}{
					"test": "Hello,Test replaceVariables Function",
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = `{"name":"Test","message":"Hello,Test replaceVariables Function"}`
			assert.Equal(t, result, util.ReplaceVariables(tt.args.reply, tt.args.variables), "Should be Same")
		})
	}
}
