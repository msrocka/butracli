package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Credentials struct {
	Endpoint string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Reads the credentials from the openLCA workspace if an .ec3 file exists
// there. Returns `nil` if no such file could be read there.
func readWorkspaceCredentials() *Credentials {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	file := filepath.Join(home, "openLCA-data-1.4", ".ec3")
	if _, err := os.Stat(file); err != nil {
		return nil
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	fileContent := Credentials{}
	if err := json.Unmarshal(data, &fileContent); err != nil {
		return nil
	}

	return &fileContent
}

// Reads the credentials from the command line arguments or from an .ec3 file
// in the openLCA workspace. Command line arguments overwrite possible arguments
// in the .ec3 file.
func readCredentials() (*Credentials, error) {

	creds := readWorkspaceCredentials()
	if creds == nil {
		creds = &Credentials{
			Endpoint: "https://buildingtransparency.org/api/",
		}
	}

	// set command line arguments
	flag := ""
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") {
			flag = arg
			continue
		}
		if flag == "" {
			continue
		}

		switch flag {
		case "-u", "-user":
			creds.User = arg
		case "-p", "-pw", "-password":
			creds.Password = arg
		case "-url", "-endpoint":
			creds.Endpoint = arg
		}
		flag = ""
	}

	if creds.User == "" || creds.Password == "" || creds.Endpoint == "" {
		return nil, errors.New("invalid credentials: " +
			"no user (-user), password (-password), or URL (-url) given")
	}

	if !strings.HasSuffix(creds.Endpoint, "/") {
		creds.Endpoint = creds.Endpoint + "/"
	}

	return creds, nil
}
