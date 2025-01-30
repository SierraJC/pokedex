package pokeapi

import (
	"fmt"
	"io"
	"net/http"

	"pokedexcli/internal/pokecache"
)

type Client struct {
	cache *pokecache.Cache
}

func NewClient(cache *pokecache.Cache) *Client {
	return &Client{cache: cache}
}

func (c *Client) Get(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}

	return body, nil
}

func (c *Client) GetWithCache(url string) ([]byte, error) {
	if c.cache == nil {
		return c.Get(url)
	}

	if body, cached := c.cache.Get(url); cached {
		return body, nil
	}

	body, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	c.cache.Add(url, body)

	return body, nil
}
