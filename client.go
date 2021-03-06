package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
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

func New(token, baseUrl string, options ...BClientOptionsFunc) *Client {
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
	client := &Client{
		token:      fmt.Sprintf("Bearer %s", token),
		baseUrl:    setupUrl,
		httpClient: new(fasthttp.Client),
	}
	for _, fn := range options {
		if fn != nil {
			err := fn(client)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	return client
}

type BClientOptionsFunc func(*Client) error

func WithMaxConnections(maxCons int) BClientOptionsFunc {
	return func(c *Client) error {
		c.httpClient.MaxConnsPerHost = maxCons
		return nil
	}
}

func WithMaxTimeoutWait(maxTimeout time.Duration) BClientOptionsFunc {
	return func(c *Client) error {
		c.httpClient.MaxConnWaitTimeout = maxTimeout
		return nil
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
			if vStr, ok := v.(string); ok {
				if len(vStr) > 0 {
					reqUri.WriteString(fmt.Sprintf("%s=%s&", k, url.QueryEscape(vStr)))
				}
			} else {
				tt := reflect.TypeOf(v)
				stype := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
				if tt.Implements(stype) {
					vStr = (v.(fmt.Stringer)).String()
					if len(vStr) > 0 {
						reqUri.WriteString(fmt.Sprintf("%s=%s&", k, url.QueryEscape(vStr)))
					}
				} else {
					var vStr string
					switch t := v.(type) {
					case int:
					case int64:
						vStr = strconv.FormatInt(t, 10)
					case uint:
					case uint64:
						vStr = strconv.FormatUint(t, 10)
					case bool:
						vStr = strconv.FormatBool(t)
					default:
						panic("cannot format")
					}
					if len(vStr) > 0 {
						reqUri.WriteString(fmt.Sprintf("%s=%s&", k, url.QueryEscape(vStr)))
					}
				}
			}
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
