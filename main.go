package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	var (
		team  string
		token string
		path  string

		input string
		tags  string
		wip   bool
	)

	flag.StringVar(&team, "team", os.Getenv("ESA_TEAM"), "")
	flag.StringVar(&token, "token", os.Getenv("ESA_TOKEN"), "")
	flag.StringVar(&path, "path", "", "")

	flag.StringVar(&input, "input", "", "")
	flag.StringVar(&tags, "tags", "", "")
	flag.BoolVar(&wip, "wip", false, "")
	flag.Parse()

	c := NewEsaClient(EsaUsingTeam(team), EsaUsingAPIKey(token))

	if len(input) == 0 {
		err := c.FindPosts(path, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	var in io.Reader

	if input == "-" {
		in = os.Stdin
	} else {
		fp, err := os.Open(input)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
		in = bufio.NewReader(fp)
	}

	err := c.WritePost(path, in, EsaPostIsWip(wip), EsaPostUsingTags(tags))
	if err != nil {
		log.Fatalf("%+v\n", err)
		return
	}
}
