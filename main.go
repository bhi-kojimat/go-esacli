package main

import (
	"errors"
	"log"
	"os"
)

const (
	requireMainArgsSize = 1
)

type command interface {
	Run() error
}

func main() {
	cmd, err := parseCommand(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseCommand(args []string) (command, error) {
	if len(args) < requireMainArgsSize {
		return nil, errors.New("require sub command. read/write")
	}

	switch args[1] {
	case "read":
		return parseReadCommand(args[2:])

	case "write":
		return parseWriteCommand(args[2:])

	default:
		return parseReadCommand(args[1:])
	}
}
