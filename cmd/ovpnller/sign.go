package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	sftp "github.com/pkg/sftp"
	ssh "golang.org/x/crypto/ssh"
)

var (
	identities Identities

	registerServer []ActionMap
)

type Identity struct {
	ident Config

	sshClient *ssh.Client
}

type Identities struct {
	server Identity
	ca     Identity

	Actions ActionMaps
}

func (is *Identities) registerServerAction() error {
	for _, actionMap := range is.Actions.RegisterServer {
		if actionMap.ActionType == "command" {
			serverResponse, err := identities.server.remoteRun(actionMap.Action)
			if err != nil {
				return err
			}
			fmt.Println(serverResponse)
		}

		if actionMap.ActionType == "copyToRemote" {
			err := identities.server.copyToRemote(actionMap.SourcePath, actionMap.TargetPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadConfig() error {
	rawBytes, err := ioutil.ReadFile(stateCLI.configPathFlag)
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

func (i *Identity) copyToRemote(fromFilePath, toFilePath string) error {

	// TODO(kompotkot): Continue https://pkg.go.dev/github.com/pkg/sftp#example-package
	sftpClient, err := sftp.NewClient(i.sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// leave your mark
	file, err := sftpClient.Create(toFilePath)
	if err != nil {
		return err
	}
	if _, err := file.Write([]byte("Hello world!")); err != nil {
		return err
	}
	file.Close()

	return nil
}

func (i *Identity) remoteRun(cmd string) (string, error) {
	session, err := i.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("Unable to create session, err: %v", err)
	}
	defer session.Close()

	var response bytes.Buffer
	session.Stdout = &response
	err = session.Run(cmd)
	return response.String(), nil
}

func (i *Identity) prepareConnection() error {
	pkRawBytes, err := ioutil.ReadFile(i.ident.PrivateKeyPath)
	if err != nil {
		return err
	}
	pkText := string(pkRawBytes)
	key, err := ssh.ParsePrivateKey([]byte(pkText))
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User:            i.ident.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	// Establish a connection
	client, err := ssh.Dial("tcp", net.JoinHostPort(i.ident.Address, "22"), config)
	if err != nil {
		return err
	}
	// TODO(kompotkot): Check if connection closed after script finished
	// defer client.Close()

	i.sshClient = client
	return nil
}

// func (is *Identities) loadPrivateKey(serverType string) error {
// 	var filePath string
// 	var identity Identity
// 	if serverType == "ca" {
// 		identity = is.server
// 		filePath = stateCLI.caSSHFilePathFlag
// 	} else if serverType == "server" {
// 		identity = is.ca
// 		filePath = stateCLI.serverSSHFilePathFlag
// 	} else {
// 		return fmt.Errorf("Unsupported server type")
// 	}

// 	rawBytes, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	text := string(rawBytes)
// 	identity.privateKey = text
// 	log.Printf("Loaded private key for %s\n", serverType)

// 	return nil
// }
