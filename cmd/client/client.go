package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/crypto2lab/ground/internal/chainspec"
	"github.com/crypto2lab/ground/internal/runtime"
)

type Client struct {
	database *Database
	runtime  *runtime.Instance
}

func (c *Client) Stop() {
	c.database.dbClient.Close()
	c.runtime.Close()
}

func StartClient(chainspecFilePath string) (*Client, error) {
	spec := chainspec.ChainSpec{}
	chainSpecBytes, err := os.ReadFile(chainspecFilePath)
	if err != nil {
		return nil, fmt.Errorf("while reading chainspec: %w", err)
	}

	err = json.Unmarshal(chainSpecBytes, &spec)
	if err != nil {
		return nil, fmt.Errorf("while unmarshaling chainspec: %w", err)
	}

	database := NewDatabase()
	err = database.Open()
	if err != nil {
		return nil, fmt.Errorf("while opening database: %w", err)
	}

	wasmBlob, err := hex.DecodeString(spec.Genesis.Code[2:])
	if err != nil {
		return nil, fmt.Errorf("decoding wasm blob: %w", err)
	}

	instance := runtime.NewInstance(context.Background(), wasmBlob)
	err = instance.Instantiate()
	if err != nil {
		return nil, fmt.Errorf("instantiating runtime: %w", err)
	}

	return &Client{
		database: database,
		runtime:  instance,
	}, nil
}
