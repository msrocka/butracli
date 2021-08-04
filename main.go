package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

	// logout
	err = logout(args)
	check("failed to logout", err)
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
	resp, err := http.Post(args.endpoint+"rest-auth/login",
		"application/json", bytes.NewReader(postData))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("login failed: " + resp.Status)
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

func logout(args *args) error {
	url := args.endpoint + "rest-auth/logout"
	fmt.Println("Logout:")
	fmt.Println("  POST", url)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+args.authKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("  response: " + resp.Status)
	return nil
}
