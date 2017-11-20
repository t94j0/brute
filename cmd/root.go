package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/t94j0/progress"
	"github.com/t94j0/ssh-bruteforce/brute"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ssh-brute [host]",
	Short: "Brute force anything!",
	Long:  `As long as it's listed in '--list'`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bruteWelcome(args[0])
	},
}

var isList bool
var host string
var bruteType string
var userLoc string
var passLoc string
var user string
var pass string

func bruteWelcome(host string) {
	implementations := brute.GetBruteImplMap()
	switch {
	case isList:
		for name, _ := range implementations {
			fmt.Println(name)
		}
		break
	case host != "":
		if _, ok := implementations[bruteType]; ok {
			bruteForce(implementations[bruteType], host)
		} else {
			fmt.Println("Error: bruteforce type does not exist")
		}
		break
	default:
		fmt.Println("Please specify a flag")
	}
}

func constructLists() (userList, passList []string) {
	userList = make([]string, 0)
	passList = make([]string, 0)
	if user != "" {
		userList = append(userList, user)
	} else {
		file, err := ioutil.ReadFile(userLoc)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}

		for _, userStr := range strings.Split(string(file), "\n") {
			userList = append(userList, userStr)
		}
	}

	if pass != "" {
		passList = append(passList, pass)
	} else {
		file, err := ioutil.ReadFile(passLoc)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}

		for _, passStr := range strings.Split(string(file), "\n") {
			passList = append(passList, passStr)
		}
	}

	return
}

func bruteForce(bruteStrategy brute.Brute, host string) {
	userList, passList := constructLists()
	tracker := progress.CreateTrackerMax("Wordlist Completion", len(userList)*len(passList))
	for _, user := range userList {
		for _, pass := range passList {
			if bruteStrategy.Try(host, user, pass) {
				fmt.Println("[+] Found: \"" + user + ":" + pass + "\"")
			}
			tracker.Increment()
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&isList, "list", "l", false, "List all installed implementations")
	RootCmd.Flags().StringVar(&user, "user", "", "Username for user")
	RootCmd.Flags().StringVar(&pass, "password", "", "Password for user")
	RootCmd.Flags().StringVar(&userLoc, "usernamePath", "", "Location for username list")
	RootCmd.Flags().StringVar(&passLoc, "passwordPath", "", "Location of passwords")
	RootCmd.Flags().StringVarP(&bruteType, "type", "t", "", "Type of brute force to do")

}
