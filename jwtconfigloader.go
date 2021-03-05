package goboxer

import (
	"encoding/json"
	"encoding/pem"
	"github.com/youmark/pkcs8"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
)

type JwtConfigLoader interface {
	Load(reader io.Reader) (*JwtConfig, error)
	DecryptPrivateKey(jwtConfig *JwtConfig) (interface{}, error)
}

type JwtConfigDefaultLoader struct{}

func (JwtConfigDefaultLoader) DecryptPrivateKey(jwtConfig *JwtConfig) (interface{}, error) {
	block, _ := pem.Decode([]byte(jwtConfig.BoxAppSettings.AppAuth.PrivateKey))
	if block == nil {
		return nil, xerrors.New("failed to decode a PEM")
	}

	pkey, _, err := pkcs8.ParsePrivateKey(
		block.Bytes,
		[]byte(jwtConfig.BoxAppSettings.AppAuth.Passphrase),
	)
	if err != nil {
		log.Printf("failed to parse private key.  %+v", err)
		return nil, xerrors.Errorf("failed to parse private key. %w", err)
	}
	return pkey, nil
}

func (JwtConfigDefaultLoader) Load(reader io.Reader) (*JwtConfig, error) {
	if reader == nil {
		return nil, xerrors.New("reader must not be null")
	}
	configFile, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to read Jwt Config File. %w", err)
	}
	jwtConfig := JwtConfig{}
	err = json.Unmarshal(configFile, &jwtConfig)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse Jwt Config File. %w", err)
	}
	return &jwtConfig, nil
}
