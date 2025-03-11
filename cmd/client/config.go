package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Queueue0/qpass/internal/dbman"
)

type Config struct {
	configPath    string
	ServerAddress string
	ServerPort    string
}

func ConfigInit() (*Config, error) {
	qpasshome, err := dbman.GetQpassHome()
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	conf.configPath = fmt.Sprintf("%s/%s", qpasshome, "config.toml")
	if _, err := os.Stat(conf.configPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, err = os.Create(conf.configPath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	_, err = toml.DecodeFile(conf.configPath, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) Save() error {
	file, err := os.Create(c.configPath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(c)
}

func (a *Application) ServerAddress() string {
	return fmt.Sprintf("%s:%s", a.Config.ServerAddress, a.Config.ServerPort)
}
