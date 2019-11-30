package solus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var maxRetries = 10

var errMaxRetriesReached = errors.New("exceeded retry limit")

// Func represents functions that can be retried.
type retryFunc func(attempt int) (retry bool, err error)

func (c *Client) request(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	var bodyByte []byte
	var reqBody io.ReadWriter
	if body != nil {
		var err error
		bodyByte, err = json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reqBody = bytes.NewBuffer(bodyByte)
	}

	fullUrl, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, method, fullUrl.String(), reqBody)
	if err != nil {
		return nil, 0, err
	}

	for k, values := range c.Headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	req.Header.Set("User-Agent", c.UserAgent)

	c.Logger.Println(method, fullUrl.String(), string(bodyByte))
	var resp *http.Response
	err = Retry(func(attempt int) (bool, error) {
		var err error
		resp, err = c.HttpClient.Do(req)
		if err != nil {
			time.Sleep(1 * time.Second)     // wait before next try
			return attempt < c.Retries, err // try N times
		}
		if resp.StatusCode == http.StatusBadGateway {
			return attempt < c.Retries, fmt.Errorf("HTTP %d", resp.StatusCode) // try N times
		}

		return false, nil
	})
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Println("failed to close body", method, path)
		}
	}()

	code := resp.StatusCode

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if code >= 404 {
		return respBody, code, fmt.Errorf("HTTP %d: %s", code, respBody)
	}

	return respBody, code, nil
}

func Retry(fn retryFunc) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}
