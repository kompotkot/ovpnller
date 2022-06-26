package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// Storing CLI definitions
	stateCLI StateCLI
)

// Command Line Interface state
type StateCLI struct {
	configureCmd  *flag.FlagSet
	signServerCmd *flag.FlagSet
	versionCmd    *flag.FlagSet

	// Common flags
	configPathFlag string
	helpFlag       bool
}

func (s *StateCLI) usage() {
	fmt.Printf(`usage: ovpnller [-h] {%[1]s,%[2]s,%[3]s} ...

OpenVPN sign workflow automation CLI

optional arguments:
    -h, --help         show this help message and exit

subcommands:
    {%[1]s,%[2]s,%[3]s}
`, s.configureCmd.Name(), s.signServerCmd.Name(), s.versionCmd.Name())
}

// Check if required flags are set
func (s *StateCLI) checkRequirements() {
	if s.helpFlag {
		switch {
		case s.configureCmd.Parsed():
			fmt.Printf("Configure commands\n\n")
			s.configureCmd.PrintDefaults()
			os.Exit(0)
		case s.signServerCmd.Parsed():
			fmt.Printf("Sign Server machine commands\n\n")
			s.signServerCmd.PrintDefaults()
			os.Exit(0)
		case s.versionCmd.Parsed():
			fmt.Printf("Show version\n\n")
			s.versionCmd.PrintDefaults()
			os.Exit(0)
		default:
			s.usage()
			os.Exit(0)
		}
	}

	if s.configPathFlag == "" {
		homePath, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Unable to parse home directory, err: %v", err)
			os.Exit(1)
		}
		s.configPathFlag = fmt.Sprintf("%s/%s", homePath, "ovpnller.json")
	}

	// if _, err := os.Stat(s.caSSHFilePathFlag); err != nil {
	// 	fmt.Sprintf("Not able to get CA machine ssh file, err: %s", err)
	// 	os.Exit(0)
	// }
	// if _, err := os.Stat(s.serverSSHFilePathFlag); err != nil {
	// 	fmt.Sprintf("Not able to get Server machine ssh file, err: %s", err)
	// 	os.Exit(0)
	// }
}

func (s *StateCLI) populateCLI() {
	// Subcommands setup
	s.configureCmd = flag.NewFlagSet("configure", flag.ExitOnError)
	s.signServerCmd = flag.NewFlagSet("sign-server", flag.ExitOnError)
	s.versionCmd = flag.NewFlagSet("version", flag.ExitOnError)

	// Common flag pointers
	for _, fs := range []*flag.FlagSet{s.configureCmd, s.signServerCmd, s.versionCmd} {
		fs.StringVar(&s.configPathFlag, "config", "", "Config file path")
		fs.BoolVar(&s.helpFlag, "help", false, "Show help message")
	}

	// for _, fs := range []*flag.FlagSet{s.signServerCmd} {
	// 	fs.StringVar(&s.caSSHFilePathFlag, "ca-ssh", "", "Path to CA machine ssh file path")
	// 	fs.StringVar(&s.serverSSHFilePathFlag, "server-ssh", "", "Path to Server machine ssh file path")
	// }

	// Configure subcommand flag pointers
	// s.configureCmd.StringVar(&s.addressFlag, "address", "127.0.0.1", "Machine IP address")
	// s.configureCmd.StringVar(&s.portFlag, "port", "22", "Machine SSH port")
	// s.configureCmd.StringVar(&s.usernameFlag, "username", "ubuntu", "Machine username identity")
	// s.configureCmd.StringVar(&s.privateKeyPath, "pk", "", "Machine SSH private key")

	// sign-server subcommand flag pointers
}
