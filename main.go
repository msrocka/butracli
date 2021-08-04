package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

	// repl
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("##> GET", args.endpoint, "")
		path, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("  ERROR: failed to read path:", err)
			continue
		}
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		if path == "q" || path == "quit" || path == "exit" || path == "halt" {
			break
		}
		err = get(args, path)
		if err != nil {
			fmt.Println("  ERROR: request failed:", err)
		}
	}

	// logout
	err = logout(args)
	check("failed to logout", err)
}

func get(args *args, path string) error {
	url := args.endpoint + path
	fmt.Println("  GET", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+args.authKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("  status:", resp.Status)

	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
