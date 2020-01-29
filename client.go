package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/upamune/go-esa/esa"
)

const (
	requirePostsSize = 1
)

type EsaClient struct {
	*esa.Client

	Team string
}

type EsaOption struct {
	Team   string
	APIKey string
}

type EsaOptionFunc func(EsaOption) EsaOption

func EsaUsingTeam(team string) EsaOptionFunc {
	return func(op EsaOption) EsaOption {
		op.Team = team
		return op
	}
}

func EsaUsingAPIKey(apiKey string) EsaOptionFunc {
	return func(op EsaOption) EsaOption {
		op.APIKey = apiKey
		return op
	}
}

type EsaPostOptionFunc func(esa.Post) esa.Post

func EsaPostIsWip(wip bool) EsaPostOptionFunc {
	return func(op esa.Post) esa.Post {
		op.Wip = wip
		return op
	}
}

func EsaPostUsingTags(tags string) EsaPostOptionFunc {
	if len(tags) == 0 {
		return func(op esa.Post) esa.Post { return op }
	}

	return func(op esa.Post) esa.Post {
		op.Tags = append(op.Tags, strings.Split(tags, ",")...)
		return op
	}
}

func NewEsaClient(options ...EsaOptionFunc) EsaClient {
	op := EsaOption{}

	for _, o := range options {
		op = o(op)
	}

	return EsaClient{
		Client: esa.NewClient(op.APIKey),
		Team:   op.Team,
	}
}

func (c EsaClient) FindPosts(path string, out io.Writer) error {
	category := filepath.Dir(path)
	name := filepath.Base(path)

	query := url.Values{}
	query.Add("category", category)
	query.Add("name", name)

	resp, err := c.Client.Post.GetPosts(c.Team, query)
	if err != nil {
		return fmt.Errorf("failed to access esa: %w", err)
	}

	if len(resp.Posts) == 0 {
		return fmt.Errorf("not found: %s, %w", path, err)
	}

	if len(resp.Posts) != requirePostsSize {
		return fmt.Errorf("too many match posts: %s", path)
	}

	_, err = io.Copy(out, strings.NewReader(resp.Posts[0].BodyMd))
	if err != nil {
		return fmt.Errorf("failed to write contents: %w", err)
	}

	return nil
}

func (c EsaClient) WritePost(path string, in io.Reader, options ...EsaPostOptionFunc) error {
	content, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("failed to read contents: %w", err)
	}

	category := filepath.Dir(path)
	name := filepath.Base(path)

	query := url.Values{}
	query.Add("category", category)
	query.Add("name", name)

	resp, err := c.Client.Post.GetPosts(c.Team, query)
	if err != nil {
		return fmt.Errorf("failed to access esa: %w", err)
	}

	req := esa.Post{}
	req.Name = name
	req.BodyMd = string(content)
	req.Category = category

	for _, op := range options {
		req = op(req)
	}

	if len(resp.Posts) == 0 {
		created, err := c.Client.Post.Create(c.Team, req)
		if err != nil {
			return fmt.Errorf("failed to write posts: %w", err)
		}

		log.Printf("Create /posts/%d\n", created.Number)

		return nil
	}

	if len(resp.Posts) != requirePostsSize {
		return fmt.Errorf("too many match posts: %s", path)
	}

	id := resp.Posts[0].Number

	updated, err := c.Client.Post.Update(c.Team, id, req)
	if err != nil {
		return fmt.Errorf("failed to write posts: %w", err)
	}

	log.Printf("Update /posts/%d\n", updated.Number)

	return nil
}
