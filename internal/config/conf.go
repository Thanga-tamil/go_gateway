package config

import (
	"os"
	"log"
	"encoding/json"

	"github.com/Thanga-tamil/noway_service/internal/utils"
)

type Cfg struct {
	Host 	 	string 		`json:"host"`
	Port 		int	 		`json:"port"`
	Postgre 	postgre 	`json:"sql"`
	Rcache 		redisCache 	`json:"redis"`
	LogLevel	int 	 	`json:"loglevel"`	 
}

type postgre struct {
	Type 		 string `json:"type"`
	Host 		 string `json:"host"`
	Port 		 int	`json:"port"`
	User 		 string `json:"user"`
	Password 	 string `json:"password"`
	Dbns 		 string `json:"dbns"`
}
	
type redisCache struct {
	Host string `json:"host"`
	Port int 	`json:"port"`
}

func LoadConfig() Cfg {

	file, err := os.Open(utils.CONF_PATH)

	if err != nil {
		log.Fatalf("Error while opening %s file \n", err.Error())
	}

	var c Cfg 
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)

	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
	log.Printf("decoded config.json file: %#v\n", c)

	return c
}

