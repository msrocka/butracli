package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Session struct {
	ID       string
	Endpoint string
}

func login(creds *Credentials) (*Session, error) {

	// post login data
	user := struct {
		User     string `json:"username"`
		Password string `json:"password"`
	}{
		User:     creds.User,
		Password: creds.Password,
	}
	postData, err := json.Marshal(&user)
	if err != nil {
		return nil, err
	}

	url := creds.Endpoint + "rest-auth/login"
	fmt.Println("Login: " + creds.User + "@" + creds.Endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewReader(postData))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("login failed: " + resp.Status)
	}

	// read key from response
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var key struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(respData, &key); err != nil {
		return nil, err
	}
	fmt.Println("OK: session=" + key.Key)

	session := &Session{
		ID:       key.Key,
		Endpoint: creds.Endpoint,
	}
	return session, nil
}

func (s *Session) logout() error {
	url := s.Endpoint + "rest-auth/logout"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.ID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("Logged out; " + resp.Status)
	return nil
}

func (s *Session) request(method, path string) error {
	url := s.Endpoint + path
	fmt.Println(" ", method, url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.ID)
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
