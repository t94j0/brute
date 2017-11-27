package cli

// RawArguments specifies the arguments in the program
// `long` is the cli argument long long. e.g. --type
// `sh` is the shorthand long e.g. --type == -t
type Arguments struct {
	// ModuleType is the module type name
	ModuleType string `long:"type" short:"t" description:"Module to select"`
	// UserPath is a newline delimited path of users
	UserPath string `long:"userpath" short:"up" description:"Path to a list of newline delimited users"`
	// User is a single user
	User string `long:"user" short:"u" description:"Selected user to brute force"`
	// PasswordPath is a list of newline delimited passwords
	PasswordPath string `long:"passwordpath" short:"pp" description:"Path to a list of newline delimited passwords"`
	// Password
	Password string `long:"password" short:"p" description:"Single password to try"`
	// Host
	Host string `long:"host" short:"ho" description:"A single host to brute force"`
	//HostPath
	HostPath string `long:"hostpath" short:"hp" description:"A list of hosts to brute force"`
	// List is a market to determine if the program should list all modules
	List []bool `long:"list" short:"l" description:"List all available modules"`
	// Help prints help messgaes
	Help []bool `long:"help" short:"h" description:"Get help from package"`
	// Represents the rest of the arguments
	Extra map[string]string
}
