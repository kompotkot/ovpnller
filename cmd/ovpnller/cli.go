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
	caInitCmd     *flag.FlagSet
	serverSignCmd *flag.FlagSet
	versionCmd    *flag.FlagSet

	// Common flags
	configPathFlag string
	helpFlag       bool
}

func (s *StateCLI) usage() {
	fmt.Printf(`usage: ovpnller [-h] {%[1]s,%[2]s,%[3]s,%[4]s} ...

OpenVPN sign workflow automation CLI

optional arguments:
    -h, --help         show this help message and exit

subcommands:
    {%[1]s,%[2]s,%[3]s,%[4]s}
`, s.configureCmd.Name(), s.caInitCmd.Name(), s.serverSignCmd.Name(), s.versionCmd.Name())
}

// Check if required flags are set
func (s *StateCLI) checkRequirements() {
	if s.helpFlag {
		switch {
		case s.configureCmd.Parsed():
			fmt.Printf("Configure commands\n\n")
			s.configureCmd.PrintDefaults()
			os.Exit(0)
		case s.caInitCmd.Parsed():
			fmt.Printf("Initialize CA machine setup\n\n")
			s.caInitCmd.PrintDefaults()
			os.Exit(0)
		case s.serverSignCmd.Parsed():
			fmt.Printf("Sign Server machine commands\n\n")
			s.serverSignCmd.PrintDefaults()
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
}

func (s *StateCLI) populateCLI() {
	// Subcommands setup
	s.configureCmd = flag.NewFlagSet("configure", flag.ExitOnError)
	s.caInitCmd = flag.NewFlagSet("ca-init", flag.ExitOnError)
	s.serverSignCmd = flag.NewFlagSet("server-register", flag.ExitOnError)
	s.versionCmd = flag.NewFlagSet("version", flag.ExitOnError)

	// Common flag pointers
	for _, fs := range []*flag.FlagSet{s.configureCmd, s.caInitCmd, s.serverSignCmd, s.versionCmd} {
		fs.StringVar(&s.configPathFlag, "config", "", "Config file path")
		fs.BoolVar(&s.helpFlag, "help", false, "Show help message")
	}
}
