package main

import (
	"errors"
	"os"
	"strings"
)

type args struct {
	endpoint string
	user     string
	password string
	authKey  string
}

func parseArgs() (*args, error) {
	args := args{
		endpoint: "https://etl-api.cqd.io/api/",
	}

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
			args.user = arg
		case "-p", "-pw", "-password":
			args.password = arg
		}
		flag = ""
	}

	if args.user == "" || args.password == "" {
		return nil, errors.New("no user (-u USER) or password (-p PASSWORD) given")
	}

	return &args, nil
}
