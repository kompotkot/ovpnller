package main

import (
	"fmt"
	"os"
)

func main() {
	stateCLI.populateCLI()
	if len(os.Args) < 2 {
		stateCLI.usage()
		os.Exit(1)
	}

	// Parse subcommands and appropriate FlagSet
	switch os.Args[1] {
	case "configure":
		stateCLI.configureCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		err := genConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Saved config file at: %s\n", stateCLI.configPathFlag)
	case "sign-server":
		stateCLI.signServerCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		err := loadConfig()
		if err != nil {
			fmt.Printf("Unable to load config, err: %v\n", err)
			os.Exit(1)
		}

		// identities.ca.prepareConnection()
		err = identities.server.prepareConnection()
		if err != nil {
			fmt.Printf("Unable initialize connection, err: %v\n", err)
			os.Exit(1)
		}

		err = identities.registerServerAction()
		if err != nil {
			fmt.Printf("Action failed, err: %v\n", err)
			os.Exit(1)
		}

	case "version":
		stateCLI.versionCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		fmt.Printf("v%s\n", OVPNLLER_VERSION)
	default:
		stateCLI.usage()
		os.Exit(1)
	}
}
