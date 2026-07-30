package service

import (
	"encoding/json"
	"github.com/Thanga-tamil/noway_service/internal/utils"
	"os"
	"time"

	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

type SecretJwk struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	K   string `json:"k"`
	Alg string `json:"alg"`
}

var JwtSignK []byte

func LoadJwtSignKeyInCache() error {

	file, err := os.Open(utils.SECRET_K)

	if err != nil {
		logger.Infof("Error while opening JWT secret key file: %s", err.Error())
		return err
	}

	var secretJwk SecretJwk

	dec := json.NewDecoder(file)
	err = dec.Decode(&secretJwk)

	if err != nil {
		logger.Infof("Error while decoding JWT secret sign key: %s", err.Error())
		return err
	}

	JwtSignK = []byte(secretJwk.K)
	logger.Info("JWT secret key has been added in inmemory")

	return nil
}

func ServeJwt(username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "username": username, 
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

    jwtToken, err := token.SignedString(JwtSignK)

	if err != nil {
		return "", err
	}

	logger.Info("Generated Jwt token: ", jwtToken)
	
	return jwtToken, nil

}

