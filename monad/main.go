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
	rawConfiguration, err := readConfiguration(filename)
	if err != nil {
		return nil, err
	}

	configuration, err := parseConfiguration(rawConfiguration)
	if err != nil {
		return nil, err
	}

	commands, err := calculateCommands(configuration)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

// END getCommandsFromFile OMIT

// START ErrorReader  OMIT
type ErrorReader struct {
	err    error
	reader *bufio.Reader
}

func NewErrorReader(reader *bufio.Reader) *ErrorReader {
	return &ErrorReader{
		reader: reader,
	}
}

func (r *ErrorReader) Err() error {
	return r.err
}

func (r *ErrorReader) ReadLine() []byte {
	if r.err != nil {
		return nil
	}

	var result []byte
	result, _, r.err = r.reader.ReadLine()
	return result
}

// END ErrorReader  OMIT

// START readConfiguration OMIT
func readConfiguration(filename string) (*RawConfiguration, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := NewErrorReader(bufio.NewReader(f)) // HL_error_in_struct
	header := reader.ReadLine()                  // HL_error_in_struct
	body := reader.ReadLine()                    // HL_error_in_struct
	if err = reader.Err(); err != nil {          // HL_error_in_struct
		return nil, err // HL_error_in_struct
	} // HL_error_in_struct

	return &RawConfiguration{
		header: header,
		body:   body,
	}, nil
}

// END readConfiguration OMIT

// START ErrorParser  OMIT
type ErrorParser struct {
	err error
}

func NewErrorParser() *ErrorParser {
	return &ErrorParser{}
}

func (p *ErrorParser) Err() error {
	return p.err
}

func (p *ErrorParser) parseVersion(configuration *RawConfiguration) int {
	if p.err != nil {
		return 0
	}

	var result int
	result, p.err = strconv.Atoi(string(configuration.header))
	return result
}

func (p *ErrorParser) parseData(configuration *RawConfiguration) map[string]string {
	if p.err != nil {
		return nil
	}

	var result map[string]string
	p.err = json.Unmarshal(configuration.body, &result)
	return result
}

// END ErrorParser  OMIT

// START ErrorChecker  OMIT
type ErrorChecker struct {
	err error
}

func NewErrorChecker() *ErrorChecker {
	return &ErrorChecker{}
}

func (c *ErrorChecker) Err() error {
	return c.err
}

func (c *ErrorChecker) StrconvAtoi(s string) int {
	if c.err != nil {
		return 0
	}

	var result int
	result, c.err = strconv.Atoi(s)
	return result
}

func (c *ErrorChecker) JsonUnmarshal(data []byte, v interface{}) {
	if c.err != nil {
		return
	}

	c.err = json.Unmarshal(data, v)
}

// END ErrorChecker  OMIT

// START parseConfiguration OMIT
func parseConfiguration(configuration *RawConfiguration) (*Configuration, error) {
	checker := NewErrorChecker()
	version := checker.StrconvAtoi(string(configuration.header))

	var data map[string]string
	checker.JsonUnmarshal(configuration.body, &data)
	if err := checker.Err(); err != nil {
		return nil, err
	}

	return &Configuration{
		Version: version,
		Data:    data,
	}, nil
}

// END parseConfiguration OMIT

// START ConfigurationCalculator OMIT
type ConfigurationCalculator struct {
	configuration *Configuration
	commands      []string
}

// END ConfigurationCalculator OMIT

// START NewConfigurationCalculator OMIT
func NewConfigurationCalculator(configuration *Configuration) *ConfigurationCalculator {
	return &ConfigurationCalculator{
		configuration: configuration,
	}
}

// END NewConfigurationCalculator OMIT

// START ConfigurationCalculator GetCommands OMIT
func (c *ConfigurationCalculator) GetCommands() []string {
	return c.commands
}

// END ConfigurationCalculator GetCommands OMIT

// START ConfigurationCalculator calculateDownCommands OMIT
func (c *ConfigurationCalculator) calculateDownCommands() error {
	downCommands, err := calculateDownCommands(c.configuration)
	c.commands = append(c.commands, downCommands...)
	return err
}

// END ConfigurationCalculator calculateDownCommands OMIT

// START ConfigurationCalculator calculateUpCommands OMIT
func (c *ConfigurationCalculator) calculateUpCommands() error {
	upCommands, err := calculateUpCommands(c.configuration)
	c.commands = append(c.commands, upCommands...)
	return err
}

// END ConfigurationCalculator calculateUpCommands OMIT

// START Do OMIT
func Do(fs ...func() error) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// END Do OMIT

// START calculateCommands OMIT
func calculateCommands(configuration *Configuration) ([]string, error) {
	calculator := NewConfigurationCalculator(configuration) // HL_monad
	if err := Do(                                           // HL_monad
		calculator.calculateDownCommands, // HL_monad
		calculator.calculateUpCommands,   // HL_monad
	); err != nil { // HL_monad
		return nil, err // HL_monad
	} // HL_monad

	return calculator.GetCommands(), nil
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
