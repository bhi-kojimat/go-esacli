package main

import (
	"flag"
	"os"
)

type readCommand struct {
	*flag.FlagSet

	Team  string
	Token string
}

func newReadCommand() readCommand {
	cmd := readCommand{
		FlagSet: flag.NewFlagSet("read", flag.ExitOnError),
	}

	cmd.StringVar(&cmd.Team, "team", os.Getenv("ESA_TEAM"), "")
	cmd.StringVar(&cmd.Token, "token", os.Getenv("ESA_TOKEN"), "")

	return cmd
}

func (cmd readCommand) Run() error {
	c := NewEsaClient(EsaUsingTeam(cmd.Team), EsaUsingAPIKey(cmd.Token))
	for _, path := range cmd.Args() {
		err := c.FindPosts(path, os.Stdout)
		if err != nil {
			return err
		}
	}
	return nil
}
