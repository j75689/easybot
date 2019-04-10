package auth

import (
	"net/http"
	"reflect"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		userID string
		aud    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestGenerateToken",
			args: args{
				userID: "Test",
				aud:    "User",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.userID, tt.args.aud)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if token, err := ParseToken(got.AccessToken); err != nil {
				t.Errorf("ParseToken [%v] error: %v", got, err)
			} else {
				if token.Subject != tt.args.userID {
					t.Errorf("UserID not match origin:[%s] parse:[%s]", tt.args.userID, token.Subject)
				}

				if token.Audience != tt.args.aud {
					t.Errorf("Audience not match origin:[%s] parse:[%s]", tt.args.aud, token.Audience)
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
		want    *jwt.StandardClaims
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
			want:    &jwt.StandardClaims{Subject: "Test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTokenFromRequest(tt.args.request)

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
