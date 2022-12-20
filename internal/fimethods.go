package internal

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func (c *Configuration) ReadConfig() *Configuration {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
