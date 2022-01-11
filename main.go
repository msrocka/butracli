package main

import (
	"bufio"
	"bytes"
	"encoding/json"
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
	creds, err := readCredentials()
	check("failed to read credentials", err)

	// login
	session, err := login(creds)
	check("failed to login", err)

	// repl
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("##> GET ", session.Endpoint, "  ")
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
		err = get(session, path)
		if err != nil {
			fmt.Println("  ERROR: request failed:", err)
		}
	}

	// logout
	err = session.logout()
	check("failed to logout", err)
}

func get(session *Session, path string) error {
	url := session.Endpoint + path
	fmt.Println("  GET", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+session.ID)
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
	var formatted bytes.Buffer
	err = json.Indent(&formatted, data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(formatted.String())
	return nil
}
