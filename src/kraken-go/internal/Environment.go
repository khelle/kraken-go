package internal

import (
	"io/ioutil"
	"../json"
)

/**
 * Environment class
 */
type Environment struct {
	Config 		*json.Json
}

/**
 * Environment constructor
 */
func CreateEnvironment() *Environment {
	env := &Environment{}

	// read configuration file
	configFile, err := ioutil.ReadFile("E:/Programy/WebServ2.1/httpd-users/Kraken-standalone/src/kraken/kraken-foundation/resource/config/kraken.json")
	if err != nil {
		return nil
	}

	// get json object
	configJson, err := json.NewJson(configFile)
	if err != nil {
		return nil
	}

	env.Config = configJson

	return env
}

/**
 * Environment.GetConfig() *json.Json
 */
func (env *Environment) GetConfig() *json.Json {
	return env.Config
}
