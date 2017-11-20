package cli

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	brute "github.com/t94j0/brute/brute-protocol"
)

// fillIn takes a list of parameters and a module to fill in those parameters to
// This returns a filled in struct of type `moduleT`
func fillIn(fillData map[string]string, moduleT brute.Brute) (brute.Brute, error) {
	module := reflect.New(reflect.TypeOf(moduleT)).Elem()

	for i := 0; i < module.NumField(); i++ {
		fieldValue := module.Field(i)
		fieldType := module.Type().Field(i)

		// Get the CLI name
		bruteName := fieldType.Tag.Get("cli")
		// Check if it's required
		_, isRequired := fieldType.Tag.Lookup("required")
		// Get any default arguments
		defaultVal, hasDefault := fieldType.Tag.Lookup("default")
		argument, hasArgument := fillData[bruteName]
		if !fieldValue.CanSet() {
			return nil, errors.New("Module Error: Cannot set parameter")
		} else if hasArgument {
			fieldValue.SetString(argument)
		} else if isRequired && hasDefault {
			fieldValue.SetString(defaultVal)
		} else if isRequired && !hasDefault {
			return nil, errors.New("Required parameter: " + bruteName)
		} else if hasDefault {
			fieldValue.SetString(defaultVal)
		}
	}

	return module.Interface().(brute.Brute), nil
}

func Brute(args Arguments, module brute.Brute) {
	for _, h := range args.Hosts {
		for _, u := range args.Users {
			for _, p := range args.Passwords {
				if ok := module.Try(h, u, p); ok {
					fmt.Printf("[+] %s: \"%s\":\"%s\"\n", h, u, p)
				}
			}
		}
	}
}

func Run() error {
	// Get arguments and available brute modules
	args, err := createArguments(os.Args[1:])
	if err != nil {
		return err
	}

	// Print help if the `-h` flag is on
	if args.Help {
		fmt.Println("General Arguments:")
		fmt.Print(args.ListString())
		if args.ModuleType == "" {
			return nil
		}
	}

	bruteMap := brute.GetBruteMap()

	// If `--list` is selected, print all modules
	if args.List {
		for key, _ := range bruteMap {
			fmt.Println(key)
		}
		return nil
	}

	// Make sure selected module exists
	moduleSkel, ok := bruteMap[args.ModuleType]
	if !ok {
		return errors.New("No such module " + args.ModuleType)
	}

	// Fill in module with data in arguments
	module, err := fillIn(args.Extra, moduleSkel)
	if err != nil {
		fmt.Println(err)
	}

	// Print out module-specific help message if `-h` flag is set
	if args.Help {
		fmt.Println("Module-specific arguments")
		fmt.Print(brute.GetCliHelp(module))
		return nil
	}

	// Make sure all generic arguments are filled in
	if !args.DidFillRequired() {
		return errors.New("Required arguments not filled in")
	}

	Brute(args, module)

	return nil
}
