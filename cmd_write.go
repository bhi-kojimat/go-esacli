package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
)

type writeCommand struct {
	*flag.FlagSet

	Team  string
	Token string

	Wip  bool
	Tags string
}

func newWriteCommand() writeCommand {
	cmd := writeCommand{
		FlagSet: flag.NewFlagSet("write", flag.ExitOnError),
	}

	cmd.StringVar(&cmd.Team, "team", os.Getenv("ESA_TEAM"), "")
	cmd.StringVar(&cmd.Token, "token", os.Getenv("ESA_TOKEN"), "")

	cmd.StringVar(&cmd.Tags, "tags", "", "")
	cmd.BoolVar(&cmd.Wip, "wip", false, "")

	return cmd
}

func (cmd writeCommand) Run() error {
	args := cmd.Args()

	var (
		in   io.Reader
		path string
	)

	switch len(args) {
	case 1:
		in = os.Stdin
		path = args[0]

	case 2:
		fp, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
		in = bufio.NewReader(fp)
		path = args[1]

	default:
		return errors.New("require filename and esa catetory/esa pagename")
	}

	c := NewEsaClient(EsaUsingTeam(cmd.Team), EsaUsingAPIKey(cmd.Token))
	return c.WritePost(path, in, EsaPostIsWip(cmd.Wip), EsaPostUsingTags(cmd.Tags))
}