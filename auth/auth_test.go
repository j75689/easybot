package auth

import (
	"net/http"
	"reflect"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/j75689/easybot/auth/claim"
	"github.com/j75689/easybot/model"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		info *model.ServiceAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestGenerateToken",
			args: args{
				info: &model.ServiceAccount{
					Name: "Test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if token, err := ParseToken(got.AccessToken); err != nil {
				t.Errorf("ParseToken [%v] error: %v", got, err)
			} else {
				if token.Name != tt.args.info.Name {
					t.Errorf("Name not match origin:[%s] parse:[%s]", tt.args.info.Name, token.Name)
				}

			}

		})
	}
}

func TestGetTokenFromRequest(t *testing.T) {
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *claim.ServiceAccountClaims
		wantErr bool
	}{
		{
			name: "TestGetTokenFromRequest",
			args: args{
				request: &http.Request{
					Header: http.Header{
						"Authorization": []string{"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJUZXN0In0.PpXF9tI9qIkPH-ZOCKJUR0I2ynXhtlsFIcl6f3DE3WLNd_So4_sHwwP0bQXVBEbOg5AbqfgLoMVopwSUc8kHnw"},
					},
				},
			},
			want:    &claim.ServiceAccountClaims{StandardClaims: jwt.StandardClaims{Subject: "Test"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := GetTokenFromRequest(tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenFromRequest() error = %v, wantErr %v, %v", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokenFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
