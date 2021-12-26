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

	"github.com/ergochat/ergo-ldap/ldap"
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

func run() (success bool, err error) {
	if len(os.Args) < 2 {
		return false, ErrNoConfigFile
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
	// sanity check; for certfp-based auth, these will be empty,
	// don't let the check pass even if the LDAP server is weirdly misconfigured
	if input.AccountName == "" || input.Passphrase == "" {
		return false, nil
	}

	ldapErr := ldap.CheckLDAPPassphrase(config, input.AccountName, input.Passphrase)
	switch ldapErr {
	case nil:
		// success
		return true, nil
	case ldap.ErrCouldNotFindUser, ldap.ErrInvalidCredentials, ldap.ErrUserNotInRequiredGroup:
		// auth failed, but not an error
		return false, nil
	default:
		return false, ldapErr
	}
}

func main() {
	var output AuthScriptOutput
	success, err := run()
	if success && err == nil {
		output.Success = true
	} else if err != nil {
		output.Error = err.Error()
	}

	out, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	out = append(out, '\n')
	os.Stdout.Write(out)
}
