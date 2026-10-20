package goboxer

import (
	"crypto/rsa"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loadTestJwtConfig(t *testing.T) (*JwtConfig, *rsa.PrivateKey) {
	t.Helper()
	f, err := os.Open("./testdata/dummykey/dummyconfig.json")
	if err != nil {
		t.Fatalf("failed to open test config: %v", err)
	}
	defer f.Close()
	loader := JwtConfigDefaultLoader{}
	conf, err := loader.Load(f)
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
	pkey, err := loader.DecryptPrivateKey(conf)
	if err != nil {
		t.Fatalf("failed to decrypt test private key: %v", err)
	}
	return conf, pkey.(*rsa.PrivateKey)
}

func TestNewJwtAuthClaim_InvalidArgs(t *testing.T) {
	conf, pkey := loadTestJwtConfig(t)

	if _, err := NewJwtAuthClaim(nil, pkey); err == nil {
		t.Errorf("NewJwtAuthClaim() with nil config should fail")
	}
	if _, err := NewJwtAuthClaim(conf, "not a key"); err == nil {
		t.Errorf("NewJwtAuthClaim() with non-rsa key should fail")
	}
	if _, err := NewJwtAuthClaimForUser(nil, pkey, "USER_ID"); err == nil {
		t.Errorf("NewJwtAuthClaimForUser() with nil config should fail")
	}
	if _, err := NewJwtAuthClaimForUser(conf, nil, "USER_ID"); err == nil {
		t.Errorf("NewJwtAuthClaimForUser() with nil key should fail")
	}
}

func TestJwtAuthClaim_Claim(t *testing.T) {
	conf, pkey := loadTestJwtConfig(t)
	const tokenUrl = "https://api.box.com/oauth2/token"

	tests := []struct {
		name        string
		newClaim    func() (*JwtAuthClaim, error)
		wantSubType string
		wantSubject string
	}{
		{
			"service account",
			func() (*JwtAuthClaim, error) { return NewJwtAuthClaim(conf, pkey) },
			"enterprise",
			conf.EnterpriseID,
		},
		{
			"user",
			func() (*JwtAuthClaim, error) { return NewJwtAuthClaimForUser(conf, pkey, "USER_ID1") },
			"user",
			"USER_ID1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claim, err := tt.newClaim()
			if err != nil {
				t.Fatalf("failed to create claim: %v", err)
			}
			before := time.Now()
			signed, err := claim.Claim(tokenUrl)
			if err != nil {
				t.Fatalf("Claim() error = %v", err)
			}

			parsed := &boxJwt{}
			token, err := jwt.ParseWithClaims(signed, parsed, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodRS256 {
					t.Errorf("unexpected signing method: %v", token.Method)
				}
				return &pkey.PublicKey, nil
			}, jwt.WithAudience(tokenUrl), jwt.WithIssuer(conf.BoxAppSettings.ClientID))
			if err != nil {
				t.Fatalf("failed to parse signed token: %v", err)
			}
			if !token.Valid {
				t.Fatalf("signed token is not valid")
			}
			if got := token.Header["kid"]; got != conf.BoxAppSettings.AppAuth.PublicKeyID {
				t.Errorf("kid = %v, want %v", got, conf.BoxAppSettings.AppAuth.PublicKeyID)
			}
			if parsed.BoxSubType != tt.wantSubType {
				t.Errorf("box_sub_type = %v, want %v", parsed.BoxSubType, tt.wantSubType)
			}
			if parsed.Subject != tt.wantSubject {
				t.Errorf("sub = %v, want %v", parsed.Subject, tt.wantSubject)
			}
			if parsed.ID == "" {
				t.Errorf("jti must not be empty")
			}
			if parsed.ExpiresAt == nil {
				t.Fatalf("exp must be set")
			}
			// Box requires exp to be within 60 seconds of issuance.
			remaining := parsed.ExpiresAt.Time.Sub(before)
			if remaining <= 0 || remaining > 60*time.Second {
				t.Errorf("exp is out of range: remaining = %v", remaining)
			}
		})
	}
}
