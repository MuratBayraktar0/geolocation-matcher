package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type HttpClient struct {
	ctx     context.Context
	Client  *http.Client
	Method  string
	Headers map[string][]string
	CB      *gobreaker.CircuitBreaker
}

func NewHttpClient(ctx context.Context) *HttpClient {
	c := &HttpClient{ctx: ctx, Headers: make(map[string][]string)}
	return c
}

func (c *HttpClient) Auth(token string) *HttpClient {
	c.Headers["Authorization"] = []string{"Bearer " + token}

	return c
}

func (c *HttpClient) Post() *HttpClient {
	c.Method = http.MethodPost

	c.Client = &http.Client{
		Timeout: 7 * time.Second,
	}

	return c
}

func (c *HttpClient) Get() *HttpClient {
	c.Method = http.MethodGet

	c.Client = &http.Client{
		Timeout: 7 * time.Second,
	}

	return c
}

func (c *HttpClient) CircuitBreaker(cb *gobreaker.CircuitBreaker) *HttpClient {
	c.CB = cb

	return c
}

func (c *HttpClient) Request(url string) (*Respose, error) {
	req, err := http.NewRequestWithContext(c.ctx, c.Method, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header[k] = v
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	if c.CB != nil {
		result, err := c.CB.Execute(func() (interface{}, error) {
			resp, err = c.Client.Do(req)
			if err != nil {
				return nil, err
			}

			return resp, nil
		})
		if err != nil {
			return nil, err
		}
		resp = result.(*http.Response)
	} else {
		resp, err = c.Client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respModel Respose
	err = json.Unmarshal(b, &respModel)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return &respModel, errors.New(http.StatusText(resp.StatusCode))
	}

	return &respModel, nil
}

type Respose struct {
	Status int
	Data   interface{}
	Error  string
	Meta   map[string]string
}
