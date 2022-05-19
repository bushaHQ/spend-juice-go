package juice

import (
	"bytes"
	"encoding/json"
	er "errors"
	"github.com/google/go-querystring/query"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
)

const (
	defaultBaseURL = "https://sandbox.spendjuice.com"
	defaultTimeout = 60 * time.Second
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client ...
type Client struct {
	httpClient HTTPClient
	baseURL    string
	apiVersion string
	apiKey     string
	debug      bool
}

// NewClient creates a new Spend-Juice API client with the default base URL.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    defaultBaseURL,
		apiKey:     os.Getenv("JUICE_PRIVATE_KEY"),
		debug:      os.Getenv("ENV") != "production",
	}
}

//SetAuth provides the client with an API key and secret.
func (cl *Client) SetAuth(apiKey string) error {
	if apiKey == "" {
		return er.New("juice: no credentials provided")
	}

	if !strings.HasPrefix(apiKey, "Bearer ") {
		apiKey = "Bearer " + apiKey
	}

	cl.apiKey = apiKey

	return nil
}

// SetHTTPClient sets the HTTP client that will be used for API calls.
func (cl *Client) SetHTTPClient(httpClient HTTPClient) {
	cl.httpClient = httpClient
}

// SetBaseURL overrides the default base URL. For internal use.
func (cl *Client) SetBaseURL(baseURL string) {
	cl.baseURL = strings.TrimRight(baseURL, "/")
}

//SetAPIVersion overrides the default base URL. For internal use.
func (cl *Client) SetAPIVersion(version string) {
	cl.apiVersion = version
}

// SetDebug enables or disables debug mode. In debug mode, HTTP requests and
// responses will be logged.
func (cl *Client) SetDebug(debug bool) {
	cl.debug = debug
}

func (cl *Client) get(path string, params interface{}, response interface{}) (err error) {
	if params != nil {

		_, err = valid.ValidateStruct(params)
		if err != nil {
			return
		}

		v, _ := query.Values(params)
		path = path + "?" + v.Encode()
	}

	url := cl.baseURL + "/" + strings.TrimLeft(path, "/")

	if cl.debug {
		log.Printf("juice: Call: %s %s", "GET", url)
		log.Printf("juice: Request Params: %#v", params)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	return cl.request(req, response)
}

func (cl *Client) post(path string, params interface{}, response interface{}) (err error) {
	url := cl.baseURL + "/" + strings.TrimLeft(path, "/")

	var req *http.Request
	var bodyBuffered io.Reader

	if params != nil {
		_, err = valid.ValidateStruct(params)
		if err != nil {
			return
		}

		data, _ := json.Marshal(params)
		bodyBuffered = bytes.NewBuffer([]byte(data))

	}

	if cl.debug {
		log.Printf("juice: Call: %s %s", "POST", url)
		log.Printf("juice: Request Params: %#v", params)
	}

	req, err = http.NewRequest(http.MethodPost, url, bodyBuffered)

	if err != nil {
		return
	}

	return cl.request(req, response)
}

func (cl *Client) patch(path string, params interface{}, response interface{}) (err error) {
	url := cl.baseURL + "/" + strings.TrimLeft(path, "/")

	var req *http.Request
	var bodyBuffered io.Reader

	if params != nil {
		_, err = valid.ValidateStruct(params)
		if err != nil {
			return
		}

		data, _ := json.Marshal(params)
		bodyBuffered = bytes.NewBuffer([]byte(data))

	}

	if cl.debug {
		log.Printf("juice: Call: %s %s", "PATCH", path)
		log.Printf("juice: Request Params: %#v", params)
	}

	req, err = http.NewRequest(http.MethodPatch, url, bodyBuffered)

	if err != nil {
		return
	}

	return cl.request(req, response)
}

func (cl *Client) request(req *http.Request, response interface{}) (err error) {

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", cl.apiKey)

	r, err := cl.httpClient.Do(req)

	if err != nil {
		return
	}

	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		e := Error{}
		err = json.NewDecoder(r.Body).Decode(&e)

		if err != nil {
			return err
		}

		return e
	}

	err = json.NewDecoder(r.Body).Decode(response)
	return
}
