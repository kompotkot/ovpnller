package main

import (
	"bufio"
	"bytes"
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
	name  string
	ident Config

	sshClient *ssh.Client
}

type Identities struct {
	server Identity
	ca     Identity

	Actions ActionMaps
}

func (i *Identity) runAction(actionMap ActionMap) error {
	switch actionMap.ActionType {
	case "command":
		err := i.remoteWorkflowRun(actionMap.Action, i.name)
		if err != nil {
			return err
		}
	case "download":
		err := i.downloadFile(actionMap.SourceFilePath, actionMap.DestinationFilePath, i.name)
		if err != nil {
			return err
		}
	case "copy":
		err := i.copyFile(actionMap.SourceFilePath, actionMap.DestinationFilePath, i.name)
		if err != nil {
			return err
		}
	case "upload":
		err := i.uploadFile(actionMap.SourceFilePath, actionMap.DestinationFilePath, i.name)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown action type: %s", actionMap.ActionType)
	}
	return nil
}

// Define an io.Reader to read from Stdin
type InputReader struct{}

// Read reads data from Stdin
func (irr *InputReader) Read(b []byte) (int, error) {
	fmt.Print("> ")
	return os.Stdin.Read(b)
}

// remoteWorkflowRun executes command and provide shell for dynamic prompt
func (i *Identity) remoteWorkflowRun(cmd, machineType string) error {
	session, err := i.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("Unable to create session, err: %v", err)
	}
	defer session.Close()

	wp, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Unable to create Stdin Pipe, err: %v", err)
	}
	rp, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to create Stdout Pipe, err: %v", err)
	}

	// Handle commands input at remote
	go func(stdin io.WriteCloser, stdout io.Reader) {
		// TODO(kompotkot): Re-write it to io.Copy or in some more clear way
		for {
			output := make([]byte, 4096)
			reader := bufio.NewReader(stdout)
			_, err := reader.Read(output)
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			}
			if err != nil {
				fmt.Println("Unable to read data", err)
				break
			}
			fmt.Printf("%s> %s", machineType, string(output))

			// Define writer to remote
			writer := bufio.NewWriter(stdin)

			var inputReader InputReader
			// Create buffer to hold input from terminal
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
		return fmt.Errorf("Command run failed, err: %v", err)
	}

	return nil
}

// downloadFile copy file from remote machine to ovpnller
func (i *Identity) downloadFile(fromFilePath, toFilePath, machineType string) error {
	// TODO(kompotkot): One sftp client for all session?
	sc, err := sftp.NewClient(i.sshClient)
	if err != nil {
		return err
	}
	defer sc.Close()

	srcFile, err := sc.OpenFile(fromFilePath, (os.O_RDONLY))
	if err != nil {
		return fmt.Errorf("Unable to open remote file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(toFilePath)
	if err != nil {
		return fmt.Errorf("Unable to open local file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("Unable to download remote file: %v", err)
	}
	fmt.Printf("File from %s machine copied to %s\n", machineType, toFilePath)

	return nil
}

func (i *Identity) uploadFile(fromFilePath, toFilePath, machineType string) error {
	// TODO(kompotkot): Continue https://pkg.go.dev/github.com/pkg/sftp#example-package
	// TODO: Remove sftp, too many dependencies
	sftpClient, err := sftp.NewClient(i.sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	file, err := sftpClient.Create(toFilePath)
	if err != nil {
		return err
	}
	if _, err := file.Write([]byte("Hello world!")); err != nil {
		return err
	}
	file.Close()

	// TODO: Fix all logs, hard to find normal language during music listening
	fmt.Printf("File to %s machine copied at %s\n", machineType, toFilePath)

	return nil
}

func (i *Identity) copyFile(fromFilePath, toFilePath, machineType string) error {
	cmd := fmt.Sprintf("cp %s %s", fromFilePath, toFilePath)
	_, err := i.remoteRun(cmd)
	if err != nil {
		return err
	}

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
