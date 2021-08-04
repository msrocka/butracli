package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func check(info string, err error) {
	if err != nil {
		fmt.Println("ERROR:", info)
		fmt.Println("  ->", err)
		os.Exit(1)
	}
}

func main() {

	// parse auth-data from program args
	args, err := parseArgs()
	check("failed to parse program arguments", err)

	// login
	authKey, err := login(args)
	check("failed to login", err)
	args.authKey = authKey
	fmt.Println("Connected to the BuildingTransparency API")
	fmt.Println("  with authentication token =", authKey)

}

func login(args *args) (string, error) {
	// post login data
	user := struct {
		User     string `json:"username"`
		Password string `json:"password"`
	}{
		User:     args.user,
		Password: args.password,
	}
	postData, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(args.endpoint+"​/rest-auth​/login",
		"application/json", bytes.NewReader(postData))
	if err != nil {
		return "", err
	}

	// read key from response
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var key struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(respData, &key); err != nil {
		return "", err
	}
	return key.Key, nil
}
