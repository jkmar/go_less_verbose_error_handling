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

// START Error  OMIT
type Error struct {
	err error
}

// END Error  OMIT

// START NewError OMIT
func NewError(err error) *Error {
	return &Error{err: err}
}

// END NewError OMIT

// START handle OMIT
func handle(err *error) {
	if r := recover(); r != nil {
		if recoveredError, ok := r.(*Error); ok {
			*err = recoveredError.err
		} else {
			panic(r)
		}
	}
}

// END handle OMIT

// START check OMIT
func check(x interface{}, err error) interface{} {
	if err != nil {
		panic(&Error{err: err})
	}
	return x
}

// END check OMIT

// START getCommandsFromFile OMIT
func getCommandsFromFile(filename string) (commands []string, err error) {
	defer handle(&err)

	rawConfiguration := check(readConfiguration(filename)).(*RawConfiguration)    // HL_check
	configuration := check(parseConfiguration(rawConfiguration)).(*Configuration) // HL_check
	commands = check(calculateCommands(configuration)).([]string)                 // HL_check
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

	reader := bufio.NewReader(f)        // HL_error_in_struct
	header, _, err := reader.ReadLine() // HL_error_in_struct
	if err != nil {                     // HL_error_in_struct
		return nil, err // HL_error_in_struct
	} // HL_error_in_struct
	// HL_error_in_struct
	body, _, err := reader.ReadLine() // HL_error_in_struct
	if err != nil {                   // HL_error_in_struct
		return nil, err // HL_error_in_struct
	} // HL_error_in_struct

	return &RawConfiguration{
		header: header,
		body:   body,
	}, nil
}

// END readConfiguration OMIT

// START parseConfiguration OMIT
func parseConfiguration(configuration *RawConfiguration) (*Configuration, error) {
	version, err := strconv.Atoi(string(configuration.header))
	if err != nil {
		return nil, err
	}

	var data map[string]string
	err = json.Unmarshal(configuration.body, &data)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Version: version,
		Data:    data,
	}, nil
}

// END parseConfiguration OMIT

// START calculateCommands OMIT
func calculateCommands(configuration *Configuration) ([]string, error) {
	var commands []string
	downCommands, err := calculateDownCommands(configuration)
	if err != nil {
		return nil, err
	}
	commands = append(commands, downCommands...)

	upCommands, err := calculateUpCommands(configuration)
	if err != nil {
		return nil, err
	}
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
