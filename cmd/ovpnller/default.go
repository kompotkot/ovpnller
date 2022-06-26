// Default configuration
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ActionMap struct {
	Action     string `json:"action"`
	ActionType string `json:"action_type"`

	// File paths
	SourcePath string `json:"source_file_path"`
	TargetPath string `json:"target_file_path"`
}

type ActionMaps struct {
	RegisterServer []ActionMap `json:"register_server"`
	RegisterClient []ActionMap `json:"register_client"`
}

// Connection configurations to CA, Server
type Config struct {
	Address        string `json:"address"`
	Port           string `json:"port"`
	Username       string `json:"username"`
	PrivateKeyPath string `json:"private_key_path"`
}

// Top level configuration file keys
type Configs struct {
	CA     Config `json:"ca"`
	Server Config `json:"server"`

	Actions ActionMaps `json:"actions"`
}

// Generate default connection configuration
func genConnectionConfig() (Config, error) {
	var data Config
	homePath, err := os.UserHomeDir()
	if err != nil {
		return data, fmt.Errorf("Unable to parse home directory, err: %v\n", err)
	}

	data = Config{
		Address:        "127.0.0.1",
		Port:           "22",
		Username:       "ubuntu",
		PrivateKeyPath: fmt.Sprintf("%s/.ssh/id_rsa", homePath),
	}

	return data, nil
}

// Generate server registration action map
func genRegisterServerConfig() []ActionMap {
	data := []ActionMap{
		{Action: "ls -la ~/", ActionType: "command"},
		{Action: "", ActionType: "copyToRemote", SourcePath: "test.txt", TargetPath: "/home/ubuntu/test.txt"},
	}

	return data
}

// Generate client registration map
func genRegisterClientConfig() []ActionMap {
	data := []ActionMap{}

	return data
}

// Check, generate and write configuration file
func genConfig() error {
	fileExists := true
	_, err := os.Stat(stateCLI.configPathFlag)
	if err != nil {
		fileExists = false
	}

	if fileExists {
		// TODO(kompotkot): Ask to rewrite the file
		return fmt.Errorf("Config file already exists %s\n", stateCLI.configPathFlag)
	}

	conConfig, err := genConnectionConfig()
	if err != nil {
		return err
	}
	regServerConfig := genRegisterServerConfig()
	regClientConfig := genRegisterClientConfig()

	// Generate full configuration
	configs := Configs{
		CA:     conConfig,
		Server: conConfig,

		Actions: ActionMaps{
			RegisterServer: regServerConfig,
			RegisterClient: regClientConfig,
		},
	}

	configsJson, err := json.Marshal(configs)
	if err != nil {
		return fmt.Errorf("Unable to marshal configuration data, err: %v\n", err)
	}

	err = ioutil.WriteFile(stateCLI.configPathFlag, configsJson, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to write default config to file, err: %v\n", err)
	}

	return nil
}
