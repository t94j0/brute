package cli

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	brute "github.com/t94j0/brute/brute-protocol"
)

func fillIn(args Arguments, moduleT brute.Brute) (brute.Brute, error) {
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
		argument, hasArgument := args.Extra[bruteName]
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

	// Make sure selected module exists
	moduleSkel, ok := bruteMap[args.ModuleType]
	if !ok {
		return errors.New("No such module " + args.ModuleType)
	}

	// Fill in module with data in arguments
	module, err := fillIn(args, moduleSkel)
	if err != nil {
		fmt.Println(err)
	}

	if args.Help {
		fmt.Println("Module-specific arguments")
		fmt.Print(brute.GetCliHelp(module))
		return nil
	}

	if !args.DidFillRequired() {
		return errors.New("Required arguments not filled in")
	}

	Brute(args, module)

	return nil
}
