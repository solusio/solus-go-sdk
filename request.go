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

type requestOpts struct {
	params map[string][]string
	body   interface{}
}

func newRequestOpts() requestOpts {
	return requestOpts{}
}

type requestWithOpt func(requestOpts) requestOpts

func withParams(params requestOpts) func(requestOpts) requestOpts {
	return func(o requestOpts) requestOpts {
		o.params = params.params
		return o
	}
}

func withBody(body requestOpts) func(requestOpts) requestOpts {
	return func(o requestOpts) requestOpts {
		o.body = body.body
		return o
	}
}

func (c *Client) request(ctx context.Context, method, path string, opts ...requestWithOpt) ([]byte, int, error) {
	reqOpts := newRequestOpts()
	for _, o := range opts {
		reqOpts = o(reqOpts)
	}

	var bodyByte []byte
	var reqBody io.ReadWriter
	if reqOpts.body != nil {
		var err error
		bodyByte, err = json.Marshal(reqOpts.body)
		if err != nil {
			return nil, 0, err
		}
		reqBody = bytes.NewBuffer(bodyByte)
	}

	fullUrl, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, 0, err
	}

	if reqOpts.params != nil {
		query := fullUrl.Query()
		for param, values := range reqOpts.params {
			for _, value := range values {
				query.Add(param, value)
			}
		}

		fullUrl.RawQuery = query.Encode()
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

func filterToParams(filter map[string]string) map[string][]string {
	var params map[string][]string
	if filter != nil {
		params = map[string][]string{}
		for field, value := range filter {
			params[field] = append(params[field], value)
		}
	}
	return params
}
