package main

import (
	"log"
	"os"
)

const (
	requireMainArgsSize = 2
)

type command interface {
	Parse([]string) error
	Run() error
}

func main() {
	var cmd command

	if len(os.Args) < requireMainArgsSize {
		log.Fatal("require sub command. read/write")
	}

	switch os.Args[1] {
	case "read":
		cmd = newReadCommand()

	case "write":
		cmd = newWriteCommand()

	default:
		log.Fatalf("unsupported command: %s", os.Args[1])
	}

	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
