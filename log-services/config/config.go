package config

import (
	es8 "github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	ESClient *es8.Client
}

func NewClient(address string) (*Client, error) {
	cfg := es8.Config{
		Addresses: []string{address},
	}
	client, err := es8.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{ESClient: client}, nil
}