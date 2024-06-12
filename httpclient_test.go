package restapireceiver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock HTTP client to simulate responses
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
func TestGetBasicAuthToken(t *testing.T) {
	username := "testuser"
	password := "testpass"

	// Create a request and set basic auth
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.SetBasicAuth(username, password)
	expectedHeader := req.Header.Get(HEADER_KEY_AUTHORIZATION)
	generatedHeader := "Basic " + getBasicAuthToken(username, password)
	assert.Equal(t, expectedHeader, generatedHeader)
}

func TestSetBasicAuth(t *testing.T) {
	helper := NewHttpClientHelper()
	username := "testuser"
	password := "testpass"
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.SetBasicAuth(username, password)
	expectedHeader := req.Header.Get(HEADER_KEY_AUTHORIZATION)

	helper.SetBasicAuth(username, password)
	assert.Equal(t, expectedHeader, helper.CommonHeaders[HEADER_KEY_AUTHORIZATION])
}

func TestSetAuthToken(t *testing.T) {
	helper := NewHttpClientHelper()
	token := "testtoken"

	helper.SetAuthToken(token)
	assert.Equal(t, token, helper.CommonHeaders[HEADER_KEY_AUTHORIZATION])
}

func TestNewRequest(t *testing.T) {
	helper := NewHttpClientHelper()
	helper.SetAuthToken("testtoken")

	req, err := helper.NewRequest(http.MethodGet, "http://example.com", nil)
	assert.NoError(t, err)
	assert.Equal(t, "testtoken", req.Header.Get(HEADER_KEY_AUTHORIZATION))
}

func TestNewJsonRequest(t *testing.T) {
	helper := NewHttpClientHelper()
	helper.SetAuthToken("testtoken")

	body := `{"key": "value"}`
	req, err := helper.NewJsonRequest(http.MethodPost, "http://example.com", body)
	assert.NoError(t, err)
	assert.Equal(t, "testtoken", req.Header.Get(HEADER_KEY_AUTHORIZATION))
	assert.Equal(t, CONTENT_TYPE_JSON, req.Header.Get(HEADER_KEY_CONTENT_TYPE))

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	assert.Equal(t, body, buf.String())
}

func TestExecuteJsonRequest(t *testing.T) {
	mockClient := new(MockHTTPClient)
	helper := NewHttpClientHelper()
	helper.Client = mockClient

	expectedResponse := map[string]interface{}{
		"key": "value",
	}
	respBody, _ := json.Marshal(expectedResponse)
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(respBody)),
	}

	req, _ := helper.NewJsonRequest(http.MethodGet, "http://example.com", "")
	mockClient.On("Do", req).Return(mockResponse, nil)

	response, err := helper.ExecuteJsonRequest(req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}
