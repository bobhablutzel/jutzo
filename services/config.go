package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func ensureConfigItem(name string) bool {
	_, isPresent := os.LookupEnv(name)
	if !isPresent {
		log.Printf("Configuration item %s is not defined", name)
	}
	return isPresent
}

func getConfig(name string, validator func(string) bool, defValue string) string {
	retrieved := os.Getenv(name)
	valid := validator(retrieved)
	if valid {
		return retrieved
	} else {
		return defValue
	}
}

func validString(input string) bool {
	return len(input) != 0
}

func getConfigString(name string, defValue string) string {
	return getConfig(name, validString, defValue)
}

func getRequiredConfigString(name string) string {
	if !ensureConfigItem(name) {
		panic(fmt.Sprintf("Required configuration item %s is missing", name))
	}
	return getConfigString(name, "")
}

func getConfigInt(name string, defValue int) int {
	stringVal := getConfigString(name, strconv.Itoa(defValue))
	retValue, err := strconv.Atoi(stringVal)
	if err != nil {
		return defValue
	} else {
		return retValue
	}
}

func getConfigItemPort() string {
	return getConfigString("JUTZO_SERVER_PORT", "8080")
}
