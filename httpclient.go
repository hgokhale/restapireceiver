package restapireceiver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

const (
	HEADER_KEY_AUTHORIZATION = "Authorization"
	HEADER_KEY_CONTENT_TYPE  = "Content-Type"
	CONTENT_TYPE_JSON        = "application/json"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
type HttpClientHelper struct {
	Client        HttpClient
	CommonHeaders map[string]string
}

func NewHttpClientHelper() *HttpClientHelper {
	return &HttpClientHelper{
		Client:        &http.Client{},
		CommonHeaders: make(map[string]string),
	}
}

func getBasicAuthToken(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (h *HttpClientHelper) SetBasicAuth(username, password string) {
	h.CommonHeaders[HEADER_KEY_AUTHORIZATION] = "Basic " + getBasicAuthToken(username, password)
}

func (h *HttpClientHelper) SetAuthToken(token string) {
	h.CommonHeaders[HEADER_KEY_AUTHORIZATION] = token
}

//TODO helper for building url

func (h *HttpClientHelper) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err == nil {
		for k, v := range h.CommonHeaders {
			req.Header.Set(k, v)
		}
	}
	return req, err
}

func (h *HttpClientHelper) NewGetRequest(url string) (*http.Request, error) {
	return h.NewRequest(http.MethodGet, url, nil)
}

func (h *HttpClientHelper) NewJsonRequest(method, url, body string) (*http.Request, error) {
	req, err := h.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte(body)))
	if err == nil {
		req.Header.Set(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_JSON)
	}
	return req, err
}

func (h *HttpClientHelper) NewPostJsonRequest(url, body string) (*http.Request, error) {
	return h.NewJsonRequest(http.MethodPost, url, body)
}

func (h *HttpClientHelper) NewPutJsonRequest(url, body string) (*http.Request, error) {
	return h.NewJsonRequest(http.MethodPut, url, body)
}

func (h *HttpClientHelper) ExecuteJsonRequest(req *http.Request) (map[string]interface{}, error) {
	var ret map[string]interface{} = nil
	resp, err := h.Client.Do(req)
	if err != nil {
		return ret, err
	}
	ret = make(map[string]interface{})
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ret, err
		}
		err = json.Unmarshal(body, &ret)
	}
	return ret, nil
}
