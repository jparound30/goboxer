package goboxer

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/youmark/pkcs8"
	"io"
	"log"
)

type JwtConfigLoader interface {
	Load(reader io.Reader) (*JwtConfig, error)
	DecryptPrivateKey(jwtConfig *JwtConfig) (any, error)
}

type JwtConfigDefaultLoader struct{}

func (JwtConfigDefaultLoader) DecryptPrivateKey(jwtConfig *JwtConfig) (any, error) {
	block, _ := pem.Decode([]byte(jwtConfig.BoxAppSettings.AppAuth.PrivateKey))
	if block == nil {
		return nil, errors.New("failed to decode a PEM")
	}

	pkey, _, err := pkcs8.ParsePrivateKey(
		block.Bytes,
		[]byte(jwtConfig.BoxAppSettings.AppAuth.Passphrase),
	)
	if err != nil {
		log.Printf("failed to parse private key.  %+v", err)
		return nil, fmt.Errorf("failed to parse private key. %w", err)
	}
	return pkey, nil
}

func (JwtConfigDefaultLoader) Load(reader io.Reader) (*JwtConfig, error) {
	if reader == nil {
		return nil, errors.New("reader must not be null")
	}
	configFile, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read Jwt Config File. %w", err)
	}
	jwtConfig := JwtConfig{}
	err = json.Unmarshal(configFile, &jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Jwt Config File. %w", err)
	}
	return &jwtConfig, nil
}
