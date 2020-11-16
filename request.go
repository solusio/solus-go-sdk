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

type requestOption func(*requestOpts)

func withFilter(f map[string]string) requestOption {
	return func(o *requestOpts) {
		if f == nil {
			return
		}

		if o.params == nil {
			o.params = map[string][]string{}
		}
		for field, value := range f {
			o.params[field] = append(o.params[field], value)
		}
	}
}

func withBody(b interface{}) requestOption {
	return func(o *requestOpts) {
		o.body = b
	}
}

func (c *Client) create(ctx context.Context, path string, data, resp interface{}) error {
	body, code, err := c.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusCreated {
		return newHTTPError(http.MethodPost, path, code, body)
	}

	return unmarshal(body, &resp)
}

func (c *Client) list(ctx context.Context, path string, resp interface{}, opts ...requestOption) error {
	body, code, err := c.request(ctx, http.MethodGet, path, opts...)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodGet, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) get(ctx context.Context, path string, resp interface{}, opts ...requestOption) error {
	body, code, err := c.request(ctx, http.MethodGet, path, opts...)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodGet, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) update(ctx context.Context, path string, data, resp interface{}) error {
	body, code, err := c.request(ctx, http.MethodPut, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodPut, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) patch(ctx context.Context, path string, data, resp interface{}) error {
	body, code, err := c.request(ctx, http.MethodPatch, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodPatch, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) delete(ctx context.Context, path string) error {
	body, code, err := c.request(ctx, http.MethodDelete, path)
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return newHTTPError(http.MethodDelete, path, code, body)
	}
	return nil
}

func (c *Client) request(ctx context.Context, method, path string, opts ...requestOption) ([]byte, int, error) {
	reqOpts := requestOpts{}
	for _, o := range opts {
		o(&reqOpts)
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

	c.Logger.Debugf("[%s] %s with body %q", method, fullUrl.String(), string(bodyByte))
	var resp *http.Response
	err = retry(func(attempt int) (bool, error) {
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
			c.Logger.Errorf("failed to close response body for %s %s: %s", method, path, err)
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

func retry(fn retryFunc) error {
	var (
		err   error
		retry bool
	)
	attempt := 1
	for {
		retry, err = fn(attempt)
		if !retry || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}

func unmarshal(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to decode %q: %w", data, err)
	}
	return nil
}
