package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	CONFIG_FILE_NAME = "ovpnller.json"
)

func loadConfig() error {
	rawBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", stateCLI.configPathFlag, CONFIG_FILE_NAME))
	if err != nil {
		return err
	}
	configs := &Configs{}
	err = json.Unmarshal(rawBytes, configs)
	if err != nil {
		return err
	}
	identities.ca.ident = configs.CA
	identities.server.ident = configs.Server
	identities.Actions = configs.Actions

	return nil
}
