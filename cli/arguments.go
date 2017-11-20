package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

// RawArguments specifies the arguments in the program
// `name` is the cli argument long name. e.g. --type
// `sh` is the shorthand name e.g. --type == -t
type Arguments struct {
	// ModuleType is the module type name
	ModuleType string `name:"type" short:"t" description:"Module to select"`
	// UserPath is a newline delimited path of users
	UserPath string `name:"userpath" short:"up" description:"Path to a list of newline delimited users"`
	// User is a single user
	User string `name:"user" short:"u" description:"Selected user to brute force"`
	// PasswordPath is a list of newline delimited passwords
	PasswordPath string `name:"passwordpath" short:"pp" description:"Path to a list of newline delimited passwords"`
	// Password
	Password string `name:"password" short:"p" description:"Single password to try"`
	// Host
	Host string `name:"host" short:"ho" description:"A single host to brute force"`
	//HostPath
	HostPath string `name:"hostpath" short:"hp" description:"A list of hosts to brute force"`
	// List is a market to determine if the program should list all modules
	List bool `name:"list" short:"l" description:"List all available modules"`
	// Help prints help messgaes
	Help bool `name:"help" short:"h" description:"Get help from package"`
	// Represents the rest of the arguments
	Extra map[string]string
}

func createArgumentMap(args []string) (map[string]string, error) {
	outputMap := make(map[string]string)

	// args[i] should never reach a value
	for i := 0; i < len(args); i++ {
		argument := args[i]

		// Check the last argument independently
		if i == len(args)-1 {
			if argument[:2] == "--" {
				variableName := argument[2:len(argument)]
				outputMap[variableName] = "true"
			} else if argument[:1] == "-" {
				variableName := argument[1:len(argument)]
				outputMap[variableName] = "true"
			} else {
				fmt.Println("Error: Formatting error")
				return nil, errors.New("Argument Format Error")
			}
		} else {
			nextValue := args[i+1]
			// Assume that if the next value starts with a `--`,
			// then the current value is a boolean
			if nextValue[:2] == "--" || nextValue[:1] == "-" {
				variableName := argument[2:len(argument)]
				outputMap[variableName] = "true"
			} else if argument[:2] == "--" {
				variableName := argument[2:len(argument)]
				outputMap[variableName] = nextValue
				i++
			} else if argument[:1] == "-" {
				variableName := argument[1:len(argument)]
				outputMap[variableName] = nextValue
				i++
			}
		}
	}

	return outputMap, nil
}

func createArgumentModule(argMap map[string]string) Arguments {
	retArgs := reflect.New(reflect.TypeOf(Arguments{})).Elem()

	strType := reflect.TypeOf("")
	extraMap := reflect.MakeMap(reflect.MapOf(strType, strType))

	for variableName, value := range argMap {
		didSet := false
		for i := 0; i < retArgs.NumField(); i++ {
			tag := retArgs.Type().Field(i).Tag
			ln := tag.Get("name")
			sn := tag.Get("short")
			if ln == variableName || sn == variableName {
				fieldType := retArgs.Field(i).Kind()
				if fieldType == reflect.Bool {
					if value == "true" {
						retArgs.Field(i).SetBool(true)
					} else if value == "false" {
						retArgs.Field(i).SetBool(false)
					}
				} else if fieldType == reflect.String {
					retArgs.Field(i).SetString(value)
				}
				didSet = true
			}
		}
		if !didSet {
			vnValue := reflect.ValueOf(variableName)
			valValue := reflect.ValueOf(value)
			extraMap.SetMapIndex(vnValue, valValue)
		}
	}

	retArgs.FieldByName("Extra").Set(extraMap)

	return retArgs.Interface().(Arguments)
}

// createArguments when inputting `os.Args`, it will output an Arguments struct
func CreateArguments(args []string) (Arguments, error) {
	argMap, err := createArgumentMap(args)
	if err != nil {
		return Arguments{}, err
	}
	return createArgumentModule(argMap), nil
}

func fileToArray(fileName string) []string {
	ret := make([]string, 0)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}

	for _, str := range strings.Split(string(file), "\n") {
		if str != "" {
			ret = append(ret, str)
		}
	}
	return ret
}

func (args Arguments) ListString() (output string) {
	argsType := reflect.TypeOf(Arguments{})
	for i := 0; i < argsType.NumField(); i++ {
		fieldTag := argsType.Field(i).Tag
		name := fieldTag.Get("name")
		short := fieldTag.Get("short")
		description := fieldTag.Get("description")
		if name != "" && short != "" && description != "" {
			output += fmt.Sprintf("\t--%s/-%s - %s\n", name, short, description)
		}
	}
	return output
}

func (args Arguments) DidFillRequired() bool {
	if args.User == "" && args.UserPath == "" {
		return false
	}
	if args.Password == "" && args.PasswordPath == "" {
		return false
	}
	if args.Host == "" && args.HostPath == "" {
		return false
	}
	return true
}
