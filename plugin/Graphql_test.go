package plugin

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestGraphql(t *testing.T) {
	log, _ := zap.NewProduction()
	logger := log.Sugar()
	type args struct {
		input     interface{}
		variables map[string]interface{}
		logger    *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Test Graphql Plugin",
			args: args{
				input: map[string]interface{}{
					"apiURL": "http://127.0.0.1:8080/graphql",
					"query": `query checkLine($lineID:String){
						checkLineIDUser(lineID:$lineID){
						success
						message
						debuglog
					  }
					}`,
					"variables": map[string]string{
						"lineID": "test",
					},
					"output": map[string]string{
						"isLineUser": "checkLineIDUser.success",
					},
				},
				variables: make(map[string]interface{}),
				logger:    logger,
			},
			want:    map[string]interface{}{"isLineUser": true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Graphql(tt.args.input, tt.args.variables, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("Graphql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Graphql() = %v, want %v", got, tt.want)
			}
		})
	}
}
