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
	accumulateCertsCmd *flag.FlagSet
	caBuildCmd         *flag.FlagSet
	serverBuildCmd     *flag.FlagSet
	serverSignCmd      *flag.FlagSet

	// Common sets
	configureCmd *flag.FlagSet
	versionCmd   *flag.FlagSet

	// Common flags
	configPathFlag string
	helpFlag       bool

	startFlag int
}

func (s *StateCLI) usage() {
	fmt.Printf(`usage: ovpnller [-h] {%[1]s,%[2]s,%[3]s,%[4]s,%[5]s,%[6]s} ...

OpenVPN sign workflow automation CLI

optional arguments:
    -h, --help         show this help message and exit

subcommands:
    {%[1]s,%[2]s,%[3]s,%[4]s,%[5]s,%[6]s}
`, s.accumulateCertsCmd.Name(), s.configureCmd.Name(), s.caBuildCmd.Name(), s.serverBuildCmd.Name(), s.serverSignCmd.Name(), s.versionCmd.Name())
}

// Check if required flags are set
func (s *StateCLI) checkRequirements() {
	if s.helpFlag {
		switch {
		case s.accumulateCertsCmd.Parsed():
			fmt.Printf("Accumulate all required certs for ovpnller in config directory\n\n")
			s.accumulateCertsCmd.PrintDefaults()
			os.Exit(0)
		case s.configureCmd.Parsed():
			fmt.Printf("Configure commands\n\n")
			s.configureCmd.PrintDefaults()
			os.Exit(0)
		case s.caBuildCmd.Parsed():
			fmt.Printf("Build CA machine (runs 'easyrsa build-ca nopass' command)\n\n")
			s.caBuildCmd.PrintDefaults()
			os.Exit(0)
		case s.serverBuildCmd.Parsed():
			fmt.Printf("Build Server machine (runs 'easyrsa gen-req server nopass' command)\n\n")
			s.serverBuildCmd.PrintDefaults()
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
		s.configPathFlag = fmt.Sprintf("%s/%s/%s", homePath, ".ovpnller", "ovpnller.json")
	}
}

func (s *StateCLI) populateCLI() {
	// Subcommands setup
	s.accumulateCertsCmd = flag.NewFlagSet("accumulate-certs", flag.ExitOnError)
	s.configureCmd = flag.NewFlagSet("configure", flag.ExitOnError)
	s.caBuildCmd = flag.NewFlagSet("ca-build", flag.ExitOnError)
	s.serverBuildCmd = flag.NewFlagSet("server-build", flag.ExitOnError)
	s.serverSignCmd = flag.NewFlagSet("server-register", flag.ExitOnError)
	s.versionCmd = flag.NewFlagSet("version", flag.ExitOnError)

	// Common flag pointers
	for _, fs := range []*flag.FlagSet{s.accumulateCertsCmd, s.configureCmd, s.caBuildCmd, s.serverBuildCmd, s.serverSignCmd, s.versionCmd} {
		fs.StringVar(&s.configPathFlag, "config", "", "Config file path")
		fs.BoolVar(&s.helpFlag, "help", false, "Show help message")
	}

	for _, fs := range []*flag.FlagSet{s.accumulateCertsCmd, s.caBuildCmd, s.serverBuildCmd, s.serverSignCmd} {
		fs.IntVar(&s.startFlag, "start", 0, "Start with actions from")
	}
}
