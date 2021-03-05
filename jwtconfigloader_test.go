package goboxer

import (
	"crypto/rsa"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestJwtConfigDefaultLoader_DecryptPrivateKey(t *testing.T) {
	validConfig, _ := os.Open("./testdata/dummykey/dummyconfig.json")
	jwtConfig, _ := JwtConfigDefaultLoader{}.Load(validConfig)

	_, _ = validConfig.Seek(0, io.SeekStart)
	pkeyempty, _ := JwtConfigDefaultLoader{}.Load(validConfig)
	pkeyempty.BoxAppSettings.AppAuth.PrivateKey = ""

	_, _ = validConfig.Seek(0, io.SeekStart)
	passWrong, _ := JwtConfigDefaultLoader{}.Load(validConfig)
	passWrong.BoxAppSettings.AppAuth.Passphrase = ""

	type args struct {
		jwtConfig *JwtConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"normal",
			args{jwtConfig: jwtConfig},
			false,
		},
		{
			"empty private key",
			args{jwtConfig: pkeyempty},
			true,
		},
		{
			"wrong passphrase",
			args{jwtConfig: passWrong},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jw := JwtConfigDefaultLoader{}
			got, err := jw.DecryptPrivateKey(tt.args.jwtConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if _, ok := got.(*rsa.PrivateKey); !ok {
					t.Errorf("DecryptPrivateKey() got unexpected key: %v", got)
				}
			}
		})
	}
}

func TestJwtConfigDefaultLoader_Load(t *testing.T) {
	type args struct {
		reader io.Reader
	}

	validConfig, _ := os.Open("./testdata/dummykey/dummyconfig.json")
	expectedJwtConfig := JwtConfig{
		BoxAppSettings: struct {
			ClientID     string `json:"clientID"`
			ClientSecret string `json:"clientSecret"`
			AppAuth      struct {
				PublicKeyID string `json:"publicKeyID"`
				PrivateKey  string `json:"privateKey"`
				Passphrase  string `json:"passphrase"`
			} `json:"appAuth"`
		}{
			ClientID:     "CLIENTID",
			ClientSecret: "CLIENTSECRET",
			AppAuth: struct {
				PublicKeyID string `json:"publicKeyID"`
				PrivateKey  string `json:"privateKey"`
				Passphrase  string `json:"passphrase"`
			}{
				PublicKeyID: "asdfzxcv",
				PrivateKey:  "-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFLTBXBgkqhkiG9w0BBQ0wSjApBgkqhkiG9w0BBQwwHAQIJQPKp9X369QCAggA\nMAwGCCqGSIb3DQIJBQAwHQYJYIZIAWUDBAEqBBB56mlhms0ZalLIZUwx9RLNBIIE\n0HGPL4YjxL38LF5L2aGQIad7NNIZLNGmqKswhNxqZylJ4AobU3MnKxQ6p2wyoXd6\nCFzymO1q9xCzHKFQcY9Q02ceeMZ+90sEK5GBEXwK0GvyhE/kBVwQeFLd6AbPP4a2\nvkhiNLgfgjtCW45Mj4KqBcxSD1WVnEOj50eaEIXeH/xieeEhxQ5r0S1Cw7MUebKZ\nLWtLsPklwikUgiB295qvE69U8JBsdY2w5E8imZYnlyrf6KHfmfxq/nAiP5WztCUQ\nXd6DgjmsDi3tPpiTKg+JnGGVXk5xsSj+X04AVobdbjbRNs1fDnOqYzOnZGRqvowk\n+kXF4IkNoiYaVy+sjYjcOXrk1UydOTF9nIRbEajpRUQ5nX6bedMDrotgHdVm9Aj0\n3nu2d48bERxd8ULHjuRlgo7urdvfQVpDq4kfY8I1bc1OuDU2K58j7BxqYsDFKxj9\ne53njyXaCdwav3Mjv0FlIQ1pJKAUruGNVVTkpq1/XOuPswxd8WGfRo/UAeDYNJLW\nxf6QcvqFQ0rkQz49h+RF4DZayEArQO/Q/p3++XrW2qe92xo2KRBqODulITzyTIF1\nfnmiXBcFxKuACAuNrk8S21pqLKUsQ/4wc3oPfbgGo2EImFRXCbUqN16h4AnNJyZ1\n4xtiU8wVij+hwLu/MTlNoCl3cWMt//l0XBqz/I0bJ9UrLkrmY/fnmnaZR+zmmzbG\n+igfbQoECUp0Gy4ukwBZ5gBwY8fs3KLuBCh7dYPQzphxaChUeeVibSo6q/0aGivn\nuDujy/hjZm2JZE0sC55cvLiUuyp00QDbpxu9iB/H/VeIZH4Zp1DHjh8Q5flJ634b\nvYprnujZRLU8j3pjiJ4b4o47+9nhKiDQJD4AfR/iwuRhp/DMfWR4QolemauJJGff\n11BZSaq3QwGahF6nt1kQ8yGomvwae18mQwRp8bDV7yQuUpWqOABLQi8IqqdsSoy0\nbLr2mpSBfTF6bU6lmKjLhcIt4qyn83Im3YZ+q52vA9kc+99TvouDTq8tqDe6aqJr\nCusXv6x/x/x+0PW2J9lg7PRHQUdglidbJ4UQRX0bl78xZxy/K2ToDI/eUhaVLbxn\nTvZNWCo1w07aUQoExgJRDnZqkxzBK8bhc4BV3jcnPwJxpW9q1Zm5ZZXKVsk+AkcE\nJhk3oH8n2FHv7zEHul5FQ2Z6wYJVYyQqCUUVC61N9NL0+VlCpX8YIsqU5ExsxqPm\nneyq1NSSb2oOBU/xPYdz1Jxx/i57jJneY0prNPi1+QOrWSqKZAmQQX+4dSjsBZmd\nypUAmjp2RPAnNy08SGx8eEEcDko2Fm2W9En9xVUOycwcrKKz9HU6HOo9W+KFeebE\nJuEO5IK/LxFslWLUuhCwLSB1QEaQg0O4D4+CcmSBdl8EBnjUmpffCWve5ClW0Vm2\n20WYvLxIorFFli2Kzl2M1qaDnW96q6VVAFCTAugW9wqjTLC5jCFO8PhXACjFsCas\nOIIaVjjnDxOO3IrcViyPlxfH1VI3Rzpyw2W+QfWf/i9PJlHgi/1CmtgqAUFuyGMw\n3wdNNAP3wUVcoxikkATjQkzrCNxWS5upi6PJXtsn2dqGbdG7TTN/MLo44sePNz2S\nxLjPWCrC9ngDkSCveiF6+N6r/rKSplsHv5cjWP63fuFC\n-----END ENCRYPTED PRIVATE KEY-----\n",
				Passphrase:  "12345",
			},
		},
		EnterpriseID: "9876",
	}

	empty := JwtConfig{}
	tests := []struct {
		name    string
		args    args
		want    *JwtConfig
		wantErr bool
	}{
		{
			"load valid config file",
			args{
				reader: validConfig,
			},
			&expectedJwtConfig,
			false,
		},
		{
			"reader is not configuration file(empty json).",
			args{
				reader: strings.NewReader("{}"),
			},
			&empty,
			false,
		},
		{
			"reader is not configuration file.",
			args{
				reader: strings.NewReader("plain text"),
			},
			nil,
			true,
		},
		{
			"reader is nil",
			args{
				reader: nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jw := JwtConfigDefaultLoader{}
			got, err := jw.Load(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
