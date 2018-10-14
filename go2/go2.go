package main

/*

const (
	latest = 2

	down = "down"
	up   = "up"

	dhcp   = "dhcp"
	static = "static"
)

type RawConfiguration struct {
	header []byte
	body   []byte
}

type Configuration struct {
	Version int
	Data    map[string]string
}

// START getCommandsFromFile OMIT
func getCommandsFromFile(filename string) ([]string, error) {
	rawConfiguration := check readConfiguration(filename) // HL_check
	configuration, := check parseConfiguration(rawConfiguration) // HL_check
	commands := check calculateCommands(configuration) // HL_check
	return commands, nil
}

// END getCommandsFromFile OMIT

// START readConfiguration OMIT
func readConfiguration(filename string) (*RawConfiguration, error) {
	f := check os.Open(filename) // HL_check
	defer f.Close()

	reader := bufio.NewReader(f)
	header, _ := check reader.ReadLine() // HL_check
	body, _ := check reader.ReadLine() // HL_check

	return &RawConfiguration{
		header: header,
		body:   body,
	}, nil
}

// END readConfiguration OMIT

// START parseConfiguration OMIT
func parseConfiguration(configuration *RawConfiguration) (*Configuration, error) {
	version := check strconv.Atoi(string(configuration.header)) // HL_check

	var data map[string]string
	check json.Unmarshal(configuration.body, &data) // HL_check

	return &Configuration{
		Version: version,
		Data:    data,
	}, nil
}

// END parseConfiguration OMIT

// START calculateCommands OMIT
func calculateCommands(configuration *Configuration) ([]string, error) {
	var commands []string

	downCommands := check calculateDownCommands(configuration) // HL_check
	commands = append(commands, downCommands...)

	upCommands := check calculateUpCommands(configuration) // HL_check
	commands = append(commands, upCommands...)

	return commands, nil
}

// END calculateCommands OMIT

// START calculateDownCommands OMIT
func calculateDownCommands(configuration *Configuration) ([]string, error) {
	switch configuration.Data[down] {
	case dhcp:
		if configuration.Version < latest {
			return nil, errors.New("DHCP not supported")
		}
		return []string{
			"pkill dhclient",
			"ifdown eth0",
		}, nil
	case static:
		return []string{
			"ifdown eth0",
		}, nil
	default:
		return nil, errors.New("unsupported configuration mode")
	}
}

// END calculateDownCommands OMIT

// END OMIT

// START calculateUpCommands OMIT
func calculateUpCommands(configuration *Configuration) ([]string, error) {
	switch configuration.Data[up] {
	case dhcp:
		if configuration.Version < latest {
			return nil, errors.New("DHCP not supported")
		}
		return []string{
			"ifup eth0",
			"dhclient",
		}, nil
	case static:
		return []string{
			"ifdown eth0",
		}, nil
	default:
		return nil, errors.New("unsupported configuration mode")
	}
}

// END calculateUpCommands OMIT

// START main OMIT
func main() {
	fmt.Println(getCommandsFromFile("resources/not_enough_lines"))
	fmt.Println(getCommandsFromFile("resources/version_not_a_number"))
	fmt.Println(getCommandsFromFile("resources/invalid_json"))
	fmt.Println(getCommandsFromFile("resources/incorrect_version"))
	fmt.Println(getCommandsFromFile("resources/incorrect_mode"))
	fmt.Println(getCommandsFromFile("resources/valid"))
}

// END main OMIT

*/
