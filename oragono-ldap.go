// Copyright (c) Shivaram Lingamneni <slingamn@cs.stanford.edu>
// Released under the Apache 2.0 license

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/oragono/oragono-ldap/ldap"
)

var (
	ErrNoConfigFile = errors.New("no config file supplied")
)

// JSON-serializable input and output types for the script
type AuthScriptInput struct {
	AccountName string `json:"accountName,omitempty"`
	Passphrase  string `json:"passphrase,omitempty"`
	Certfp      string `json:"certfp,omitempty"`
	IP          string `json:"ip,omitempty"`
}

type AuthScriptOutput struct {
	AccountName string `json:"accountName"`
	Success     bool   `json:"success"`
	Error       string `json:"error"`
}

func LoadRawConfig(filename string) (config ldap.ServerConfig, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return
	}
	return
}

func run() (err error) {
	if len(os.Args) < 2 {
		return ErrNoConfigFile
	}
	config, err := LoadRawConfig(os.Args[1])
	if err != nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return
	}

	var input AuthScriptInput
	err = json.Unmarshal(line, &input)
	if err != nil {
		return
	}

	err = ldap.CheckLDAPPassphrase(config, input.AccountName, input.Passphrase)
	return
}

func main() {
	var output AuthScriptOutput
	err := run()
	if err == nil {
		output.Success = true
	} else {
		output.Success = false
		output.Error = err.Error()
	}

	out, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	out = append(out, '\n')
	os.Stdout.Write(out)
}
