package main

import (
	"os"
	"strconv"
)

type Configuration struct {
}

// GetConfigurationString returns the configuration item
// value. The second return determines if the item was
// actually provided. If the item was not provided the
// routine should return ("", false); otherwise it should
// provide (value, true)
func (Configuration) GetConfigurationString(name string) (string, bool) {
	return os.LookupEnv(name)
}

// GetConfigurationInt is similar to GetConfigurationString
// but validates that the item is both provided and an integer
// If the item is not provided or not an integer, this should
// return (0, false); otherwise it should return (value, true)
func (c Configuration) GetConfigurationInt(name string) (int, bool) {

	if value, isPresent := c.GetConfigurationString(name); isPresent {
		if retValue, err := strconv.Atoi(value); err == nil {
			return retValue, true
		}
	}
	return 0, false
}
