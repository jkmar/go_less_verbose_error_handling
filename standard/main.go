package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

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
	// START GenericMonadGetCommandsFromFile OMIT
	rawConfiguration, err := readConfiguration(filename) // HL_generic_monad
	if err != nil {                                      // HL_generic_monad
		return nil, err // HL_generic_monad
	} // HL_generic_monad

	configuration, err := parseConfiguration(rawConfiguration) // HL_generic_monad
	if err != nil {                                            // HL_generic_monad
		return nil, err // HL_generic_monad
	} // HL_generic_monad

	commands, err := calculateCommands(configuration) // HL_generic_monad
	if err != nil {                                   // HL_generic_monad
		return nil, err // HL_generic_monad
	} // HL_generic_monad
	// END GenericMonadGetCommandsFromFile OMIT

	return commands, nil
}

// END getCommandsFromFile OMIT

// START readConfiguration OMIT
func readConfiguration(filename string) (*RawConfiguration, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f) // HL_error_in_struct
	// START ReadLineReadConfiguration OMIT
	header, _, err := reader.ReadLine() // HL_error_in_struct
	if err != nil {                     // HL_error_in_struct
		return nil, err // HL_error_in_struct
	} // HL_error_in_struct
	// HL_error_in_struct
	body, _, err := reader.ReadLine() // HL_error_in_struct
	if err != nil {                   // HL_error_in_struct
		return nil, err // HL_error_in_struct
	} // HL_error_in_struct
	// END ReadLineReadConfiguration OMIT

	return &RawConfiguration{
		header: header,
		body:   body,
	}, nil
}

// END readConfiguration OMIT

// START parseConfiguration OMIT
func parseConfiguration(configuration *RawConfiguration) (*Configuration, error) {
	// START ErrorCheckerParseConfiguration OMIT
	version, err := strconv.Atoi(string(configuration.header)) // HL_check
	if err != nil {                                            // HL_check
		return nil, err // HL_check
	} // HL_check

	var data map[string]string                      // HL_check
	err = json.Unmarshal(configuration.body, &data) // HL_check
	if err != nil {                                 // HL_check
		return nil, err // HL_check
	} // HL_check
	// END ErrorCheckerParseConfiguration OMIT

	return &Configuration{
		Version: version,
		Data:    data,
	}, nil
}

// END parseConfiguration OMIT

// START calculateCommands OMIT
func calculateCommands(configuration *Configuration) ([]string, error) {
	var commands []string
	// START MonadCalculateCommands OMIT
	downCommands, err := calculateDownCommands(configuration) // HL_monad
	if err != nil {                                           // HL_monad
		return nil, err // HL_monad
	} // HL_monad
	commands = append(commands, downCommands...)

	upCommands, err := calculateUpCommands(configuration) // HL_monad
	if err != nil {                                       // HL_monad
		return nil, err // HL_monad
	} // HL_monad
	commands = append(commands, upCommands...)
	// END MonadCalculateCommands OMIT

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
