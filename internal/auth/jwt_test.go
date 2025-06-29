package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwt(t *testing.T) {
	const secret string = "mykey"

	type testCase struct {
		secret string
		userIdToEncode uuid.UUID
		userIdToCheck uuid.UUID
		expiresIn time.Duration
		waitTime time.Duration
		matchExpected bool
		timeOutExpected bool
	}

	cases := []testCase{}

	id := uuid.New()
	cases = append(cases, testCase{
		secret: secret,
		userIdToEncode: id,
		userIdToCheck: id,
		expiresIn: 5 * time.Second,
		waitTime: 0 * time.Second,
		matchExpected: true,
		timeOutExpected: false,
	})

	cases = append(cases, testCase{
		secret: secret,
		userIdToEncode: id,
		userIdToCheck: uuid.New(),
		expiresIn: 5 * time.Second,
		waitTime: 0 * time.Second,
		matchExpected: false,
		timeOutExpected: false,
	})
	
	cases = append(cases, testCase{
		secret: secret,
		userIdToEncode: id,
		userIdToCheck: id,
		expiresIn: 0 * time.Second,
		waitTime: 1 * time.Second,
		matchExpected: true,
		timeOutExpected: true,
	})

	for _, c := range cases {
		token, err := MakeJWT(c.userIdToEncode, c.secret, c.expiresIn)
		if err != nil {
			t.Errorf("Error, creating JWT: %v", err)
			return
		}

		time.Sleep(c.waitTime)

		id, err := ValidateJWT(token, c.secret)
		if err != nil {
			if c.timeOutExpected {
				break
			}

			t.Errorf("Error, decoding JWT: %v", err)
			return
		}

		if c.timeOutExpected {
			t.Error("Expired Token did not error")
			return
		}

		if id != c.userIdToCheck {
			if c.matchExpected {
				t.Errorf("%v != %v", id, c.userIdToCheck)
				return
			}

			break
		}
		
		if !c.matchExpected {
			t.Errorf("%v == %v", id, c.userIdToCheck)
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
