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
	case "accumulate-certs":
		stateCLI.accumulateCertsCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		err := loadConfig()
		if err != nil {
			fmt.Printf("Unable to load config, err: %v\n", err)
			os.Exit(1)
		}
		err = identities.ca.prepareConnection()
		if err != nil {
			fmt.Printf("Unable initialize connection with CA, err: %v\n", err)
			os.Exit(1)
		}

		for _, actionMap := range identities.Actions.AccumulateCerts[stateCLI.startFlag:] {
			err := identities.server.runAction(actionMap)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	case "configure":
		// TODO(kompotkot): Fix, outdated (need changes to work with config directory and new commands)
		stateCLI.configureCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		err := genConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Saved config file at: %s\n", stateCLI.configPathFlag)
	case "ca-build":
		stateCLI.caBuildCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		err := loadConfig()
		if err != nil {
			fmt.Printf("Unable to load config, err: %v\n", err)
			os.Exit(1)
		}

		err = identities.ca.prepareConnection()
		if err != nil {
			fmt.Printf("Unable initialize connection with CA, err: %v\n", err)
			os.Exit(1)
		}

		for _, actionMap := range identities.Actions.CaBuild[stateCLI.startFlag:] {
			err := identities.server.runAction(actionMap)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	case "server-build":
		stateCLI.serverBuildCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		err := loadConfig()
		if err != nil {
			fmt.Printf("Unable to load config, err: %v\n", err)
			os.Exit(1)
		}

		err = identities.server.prepareConnection()
		if err != nil {
			fmt.Printf("Unable initialize connection with Server, err: %v\n", err)
			os.Exit(1)
		}

		for _, actionMap := range identities.Actions.ServerBuild[stateCLI.startFlag:] {
			err := identities.server.runAction(actionMap)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	case "server-sign":
		stateCLI.serverSignCmd.Parse(os.Args[2:])
		stateCLI.checkRequirements()

		err := loadConfig()
		if err != nil {
			fmt.Printf("Unable to load config, err: %v\n", err)
			os.Exit(1)
		}

		// err = identities.ca.prepareConnection()
		err = identities.server.prepareConnection()
		if err != nil {
			fmt.Printf("Unable initialize connection with Server, err: %v\n", err)
			os.Exit(1)
		}

		for _, actionMap := range identities.Actions.ServerRegister[stateCLI.startFlag:] {
			err := identities.server.runAction(actionMap)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
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
