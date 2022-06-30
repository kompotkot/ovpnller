package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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

// Define an io.Reader to read from Stdin
type InputReader struct{}

// Read reads data from Stdin
func (irr *InputReader) Read(b []byte) (int, error) {
	fmt.Print("> ")
	return os.Stdin.Read(b)
}

func (i *Identity) remoteWorkflowRun(cmd, machine_type string) error {
	session, err := i.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("Unable to create session, err: %v", err)
	}
	defer session.Close()

	wp, err := session.StdinPipe()
	if err != nil {
		return err
	}
	rp, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	// Handle commands input at remote
	go func(stdin io.WriteCloser, stdout io.Reader) {
		for {
			reader := bufio.NewReader(stdout)
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			}
			if err != nil {
				fmt.Println("Unable to read data", err)
				break
			}
			fmt.Printf("server> %s", line)

			// Define writer to remote
			writer := bufio.NewWriter(stdin)

			var inputReader InputReader
			// Create buffer to hold input from terminal
			// TODO(kompotkot): Re-write it to io.Copy or in some more clear way
			input := make([]byte, 4096)
			_, err = inputReader.Read(input)
			if err != nil {
				fmt.Println("Unable to read input data", err)
				break
			}
			// Pass input to remote writer
			_, err = writer.Write(input)
			if err != nil {
				fmt.Println("Unable to write data from input", err)
				break
			}
			writer.Flush()
		}
	}(wp, rp)

	// Execute startup command
	var response bytes.Buffer
	session.Stdout = &response
	err = session.Run(cmd)
	if err != nil {
		return err
	}
	// Last message from remote
	fmt.Println(response.String())

	return nil
}

func (is *Identities) runAction(activeAction string) error {
	switch {
	case "ca-init" == activeAction:
		for _, actionMap := range is.Actions.CaInit {
			fmt.Println(actionMap.Action)
			err := identities.ca.remoteWorkflowRun(actionMap.Action, "ca")
			if err != nil {
				return err
			}
		}
	case "server-register" == activeAction:
		for _, actionMap := range is.Actions.ServerRegister {
			err := identities.server.remoteWorkflowRun(actionMap.Action, "server")
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Unknown action: %s\n", activeAction)
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
