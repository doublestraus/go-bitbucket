package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"net/url"
	"strings"
)

type Client struct {
	token, baseUrl string
	httpClient     *fasthttp.Client
}

type Errors struct {
	StatusCode int           `json:"-"`
	Errors     []ClientError `json:"errors"`
}

type ClientError struct {
	Message       string `json:"message"`
	Context       string `json:"context"`
	ExceptionName string `json:"exceptionName"`
}

func (e Errors) Error() string {
	errorStr := fmt.Sprintf("status code = %d\n", e.StatusCode)
	for _, e := range e.Errors {
		errorStr += fmt.Sprintf("error = %s\n\n", e.Message)
	}
	return errorStr
}

func New(token, baseUrl string) *Client {
	setupUrl := ""
	if strings.HasPrefix(baseUrl, "https://") {
		setupUrl = baseUrl
	} else {
		setupUrl = "https://" + baseUrl
	}
	if strings.HasSuffix(setupUrl, "/") {
		setupUrl = setupUrl + "rest/api/1.0"
	} else {
		setupUrl = setupUrl + "/rest/api/1.0/"
	}
	return &Client{
		token:      fmt.Sprintf("Bearer %s", token),
		baseUrl:    setupUrl,
		httpClient: new(fasthttp.Client),
	}
}

type Response struct {
	pagination Pagination
}

func (c *Client) get(path string, pagination *Pagination, filter map[string]interface{}) ([]byte, error) {
	var resp *fasthttp.Response
	var body []byte
	var err error
	if pagination == nil {
		err = errors.New("provide pagination to paged api")
	} else {
		start := pagination.Start
		limit := pagination.Limit
		if filter != nil {
			resp, err = c.do("GET", path, start, limit, filter)
		} else {
			resp, err = c.do("GET", path, start, limit, nil)
		}
		if err == nil {
			if resp.StatusCode() >= 400 {
				var jerr Errors
				err = json.Unmarshal(resp.Body(), &jerr)
				jerr.StatusCode = resp.StatusCode()
				if err != nil {
					return nil, jerr
				} else {
					return nil, errors.New(fmt.Sprintf("status code = %d. error = %s\n\n", resp.StatusCode(), resp.Body()))
				}
			} else {
				writer := bytes.NewBuffer([]byte{})
				reader := bytes.NewReader(resp.Body())
				_, err = io.Copy(writer, reader)
				if err == nil {
					body = writer.Bytes()
				}
				fasthttp.ReleaseResponse(resp)
			}
		}
	}
	return body, err
}

func (c *Client) post() {

}

func (c *Client) do(method, path string, start, limit int, filter map[string]interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.Header.Set("Authorization", c.token)
	if filter != nil {
		reqUri := new(strings.Builder)
		reqUri.WriteString(fmt.Sprintf("%s%s?", c.baseUrl, path))
		for k, v := range filter {
			reqUri.WriteString(fmt.Sprintf("%s=%s&", k, url.QueryEscape(v.(string))))
		}
		reqUri.WriteString(fmt.Sprintf("start=%d&limit=%d", start, limit))
		req.SetRequestURI(reqUri.String())
	} else {
		req.SetRequestURI(fmt.Sprintf("%s%s?start=%d&limit=%d", c.baseUrl, path, start, limit))
	}
	req.Header.SetMethod(method)
	resp := fasthttp.AcquireResponse()
	err := c.httpClient.Do(req, resp)
	if err != nil {
		return nil, err
	}
	fasthttp.ReleaseRequest(req)
	return resp, nil
}
