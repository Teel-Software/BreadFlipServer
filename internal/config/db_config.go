package config

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	User    string `yaml:"user"`
	Pswd    string `yaml:"password"`
	DBName  string `yaml:"db"`
	CrtPath string `yaml:"crt_path"`
}

var DbConf DBConfig

func LoadDBCfg(cfgPath string) error {
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&DbConf)
	if err != nil {
		return err
	}

	configBytes, err := json.MarshalIndent(&DbConf, "", "  ")
	if err != nil {
		return err
	}

	log.Println("Running config:", string(configBytes))
	return nil
}

func GetDBConfig() *DBConfig {
	LoadDBCfg("myconfig.yaml")
	return &DbConf
}
