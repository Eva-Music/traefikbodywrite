package traefikbodywrite

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	//NewBodyContent []NewBody
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	GrantType    string `json:"grantType,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// transformer plugin.
type transformer struct {
	next   http.Handler
	config Config
}

// New created a new transformer plugin.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	if config.ClientSecret == "" && config.ClientId == "" && config.GrantType == "" {
		return nil, fmt.Errorf("some required fields are empty")
	}

	return &transformer{
		next:   next,
		config: *config,
	}, nil
}

func (a *transformer) log(format string) {
	_, writeLogError := os.Stderr.WriteString(format)
	if writeLogError != nil {
		panic(writeLogError.Error())
	}
}

func (a *transformer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	data := url.Values{}

	data.Set("client_id", a.config.ClientId)
	data.Set("client_secret", a.config.ClientSecret)
	data.Set("grant_type", a.config.GrantType)

	usernameHeader := req.Header.Values("username")[0]
	if usernameHeader == "" {
		http.Error(rw, "username header missing", http.StatusInternalServerError)
	}
	data.Set("username", usernameHeader)
	req.Header.Del("username")

	passwordHeader := req.Header.Values("password")[0]
	if passwordHeader == "" {
		http.Error(rw, "password header missing", http.StatusInternalServerError)
	}
	data.Set("password", passwordHeader)
	req.Header.Del("password")

	var body io.Reader
	body = bytes.NewBufferString(data.Encode())

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}

	req.Body = rc

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//req.Body = io.NopCloser(strings.NewReader(string(jsonBody)))
	//req.ContentLength = int64(len(jsonBody))

	a.next.ServeHTTP(rw, req)
}
