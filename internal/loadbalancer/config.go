package loadbalancer

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	asciiHeader = `
   ____                  _                                          
  / __ \                (_)                                         
 | |  | |_ __ ___  _ __  _                                          
 | |  | | '_ ' _ \| '_ \| |                                         
 | |__| | | | | | | | | | |                                         
  \____/|_| |_| |_|_| |_|_| ____        _                           
 | |                   | | |  _ \      | |                          
 | |     ___   __ _  __| | | |_) | __ _| | __ _ _ __   ___ ___ _ __ 
 | |    / _ \ / _' |/ _' | |  _ < / _' | |/ _' | '_ \ / __/ _ \ '__|
 | |___| (_) | (_| | (_| | | |_) | (_| | | (_| | | | | (_|  __/ |   
 |______\___/ \__,_|\__,_| |____/ \__,_|_|\__,_|_| |_|\___\___|_|   

`
)

type Config struct {
	Algorithm string `yaml:"algorithm"`
}

// ReadConfig read configuration from `fileName` file
func ReadConfig(fileName string) (*Config, error) {
	in, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Print() {
	println("%s\n", asciiHeader)
	fmt.Printf("Algorithm: %s\n", c.Algorithm)
}

func (c *Config) IsValid() error {
	if c.Algorithm != "round-robin" {
		return errors.New("the algorithm is unknown")
	}

	return nil
}
