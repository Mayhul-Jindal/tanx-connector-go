package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignPayload(t *testing.T) {
	testCases := []struct {
		name    string
		Params  []string
		wantErr bool
	}{
		{
			name:    "0x present in private key",
			Params:  []string{"eww", "0x1"},
			wantErr: true,
		},
		{
			name:    "invalid private key",
			Params:  []string{"eww", "1111111111223232423423424"},
			wantErr: true,
		},
		{
			name:    "successful sign",
			Params:  []string{"eww", "7d6384d6877be027aa25bd458f2058e3f7ff68347dc583a9baf96f5f97b413a8"},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := SignPayload(tc.Params[0], tc.Params[1])
			if tc.wantErr {
				assert.Error(t, err)
				t.Logf("Error: %v", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNonce(t *testing.T) {
	testCases := []struct {
		name         string
		roundTripper http.RoundTripper
		timeout      time.Duration
		want         NonceResponse
		wantErr      bool
	}{
		{
			name: "Successful roundtrip",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
                        "status": "success",
                        "message": "Working fine!",
                        "payload": ""
                    }`)),
				},
			},
			want: NonceResponse{
				Status:  "success",
				Message: "Working fine!",
				Payload: "",
			},
			wantErr: false,
		},
		{
			name: "HTTP protocol error",
			roundTripper: &MockRoundTripper{
				Err: errors.New("mock error"),
			},
			wantErr: true,
		},
		{
			name: "JSON decoding error",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"status: 1}`)),
				},
			},
			wantErr: true,
		},
		{
			name: "Context timeout error",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"status": "success"}`)),
				},
				Delay: time.Second * 2,
			},
			timeout: time.Nanosecond,
			wantErr: true,
		},

		// status code 404
		{
			name: "404 status code",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"status": "error", "message": "Not found"}`)),
				},
			},
			wantErr: true,
		},

		// 500 status code
		{
			name: "500 status code",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`{"status": "error", "message": "Internal server error"}`)),
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := New(TESTNET)
			assert.NoError(t, err)

			client.httpClient.Transport = tc.roundTripper

			ctx := context.Background()
			if tc.timeout != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tc.timeout)
				defer cancel()
			}

			got, err := client.Nonce(ctx, "0x1234567890")
			if tc.wantErr {
				assert.Error(t, err)
				t.Logf("Error: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestRefreshTokens(t *testing.T) {
	testCases := []struct {
		name         string
		roundTripper http.RoundTripper
		timeout      time.Duration
		want         RefreshTokenResponse
		wantErr      bool
	}{
		{
			name: "Successful roundtrip",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"status": "success",
						"message": "",
						"payload": {
							"access": "eybdgyr0iOi....gVbdghy7g_Pn3QCU",
							"refresh": "eyJ0eXAiOi....mVeL25ytndhltJL4"
						}
					}`)),
				},
			},
			want: RefreshTokenResponse{
				Status:  "success",
				Message: "",
				Payload: RefreshTokenPayload{
					Access:  "eybdgyr0iOi....gVbdghy7g_Pn3QCU",
					Refresh: "eyJ0eXAiOi....mVeL25ytndhltJL4",
				},
			},
			wantErr: false,
		},

		// 404 status code
		{
			name: "404 status code",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"status": "error", "message": "Not found"}`)),
				},
			},
			wantErr: true,
		},

		// 500 error
		{
			name: "500 status code",
			roundTripper: &MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`{"status": "error", "message": "Internal server error"}`)),
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := New(TESTNET)
			assert.NoError(t, err)

			client.httpClient.Transport = tc.roundTripper

			ctx := context.Background()
			if tc.timeout != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tc.timeout)
				defer cancel()
			}

			got, err := client.RefreshTokens(ctx)
			if tc.wantErr {
				assert.Error(t, err)
				t.Logf("Error: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name         string
		roundTripper roundTripFunc
		params       [2]string
		timeout      time.Duration
		want         JWTResponse
		wantErr      bool
	}{
		{
			name: "Successful roundtrip",
			roundTripper: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == NONCE_ENDPOINT {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"status": "success"}`)),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"status": "success",
						"message": "",
						"payload": {
							"uid": "IDDAF4F5E3C7",
							"signature": ""
						},
						"token": {
							"refresh": "",
							"access": ""
						}
					}`)),
				}, nil
			}),
			params: [2]string{"0x713Cf80b7c71440E7a09Dede1ee23dCBf862fB66", "7d6384d6877be027aa25bd458f2058e3f7ff68347dc583a9baf96f5f97b413a8"},
			want: JWTResponse{
				Status:  "success",
				Message: "",
				Payload: JWTPayload{
					UID:       "IDDAF4F5E3C7",
					Signature: "",
				},
				Token: Token{
					Refresh: "",
					Access:  "",
				},
			},
			wantErr: false,
		},
		{
			name: "Nonce error",
			roundTripper: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == NONCE_ENDPOINT {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"status": "error"}`)),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"status": "success"}`)),
				}, nil
			}),
			wantErr: true,
		},

		// status code 404
		{
			name: "404 status code",
			roundTripper: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == NONCE_ENDPOINT {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"status": "success", "message": "Not found"}`)),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"status": "error", "message": "Not found"}`)),
				}, nil
			}),
			params:  [2]string{"0x713Cf80b7c71440E7a09Dede1ee23dCBf862fB66", "7d6384d6877be027aa25bd458f2058e3f7ff68347dc583a9baf96f5f97b413a8"},
			wantErr: true,
		},

		// 500 status error
		{
			name: "500 status code",
			roundTripper: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == NONCE_ENDPOINT {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"status": "success"}`)),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`{"status": "error", "message": "Internal server error"}`)),
				}, nil
			}),
			params:  [2]string{"0x713Cf80b7c71440E7a09Dede1ee23dCBf862fB66", "7d6384d6877be027aa25bd458f2058e3f7ff68347dc583a9baf96f5f97b413a8"},
			wantErr: true,
		},

		// invalid params
		{
			name: "invalid ethAddress and privateKey",
			roundTripper: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == NONCE_ENDPOINT {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"status": "success"}`)),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"status": "success"}`)),
				}, nil
			}),
			params:  [2]string{"", ""},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := New(TESTNET)
			assert.NoError(t, err)

			client.httpClient.Transport = tc.roundTripper

			ctx := context.Background()
			if tc.timeout != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tc.timeout)
				defer cancel()
			}

			got, err := client.Login(ctx, tc.params[0], tc.params[1])
			if tc.wantErr {
				assert.Error(t, err)
				t.Logf("Error: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
