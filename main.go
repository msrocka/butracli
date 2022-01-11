package main

import (
	"bufio"
	"errors"
	"fmt"
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

		fmt.Print("##> ")
		cmd, err := readCommand(reader)
		if err != nil {
			fmt.Println("  ERROR: failed to read command:", err)
			continue
		}
		if cmd.isExit() {
			break
		}
		err = cmd.exec(session)
		if err != nil {
			fmt.Println("  ERROR: ", err)
		}
	}

	// logout
	err = session.logout()
	check("failed to logout", err)
}

type command struct {
	method string
	path   string
	data   string
}

func (c *command) isExit() bool {
	switch c.method {
	case "q", "quit", "exit", "halt", "end":
		return true
	default:
		return false
	}
}

func readCommand(r *bufio.Reader) (*command, error) {

	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	parts := strings.Split(line, " ")

	cmd := &command{}
	k := 0
	for i := range parts {
		s := strings.TrimSpace(parts[i])
		if s == "" {
			continue
		}
		switch k {
		case 0:
			cmd.method = s
		case 1:
			cmd.path = strings.TrimLeft(s, "/")
		case 2:
			cmd.data += " " + s
		}
		k++
	}

	if cmd.method == "" {
		return nil, errors.New("command is empty")
	}

	return cmd, nil
}

func (c *command) exec(s *Session) error {
	switch strings.ToUpper(c.method) {
	case "GET":
		return s.request(http.MethodGet, c.path)
	case "DELETE":
		return s.request(http.MethodDelete, c.path)
	default:
		return errors.New("invalid/unsupported HTTP method: " + c.method)
	}

}
