// Package traefikbodywrite plugin.
package traefikbodywrite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Config the plugin configuration.

//type NewBody struct {
//	NewKey   string `json:"newKey,omitempty"`
//	NewValue string `json:"newValue,omitempty"`
//}

type Config struct {
	//NewBodyContent []NewBody
	NewContentType string `json:"newContentType,omitempty"`
	ClientId       string `json:"clientId,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty"`
	GrantType      string `json:"grantType,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		//NewBodyContent: []NewBody{},
		NewContentType: "application/x-www-form-urlencoded",
	}
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

//newBodyContent:
//   - newKey: "client_id"
//     newValue: SerialNumber%3D

func (a *transformer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	newBody := make(map[string]string)

	//create new body
	newBody["client_id"] = a.config.ClientId
	newBody["client_secret"] = a.config.ClientSecret
	newBody["grant_type"] = a.config.GrantType

	// add headers to new body
	usernameHeader := req.Header.Get("username")
	passwordHeader := req.Header.Get("password")

	if usernameHeader == "" {
		http.Error(rw, "username header missing", http.StatusInternalServerError)
	}
	newBody["username"] = usernameHeader
	req.Header.Del("username")

	if passwordHeader == "" {
		http.Error(rw, "password header missing", http.StatusInternalServerError)
	}
	newBody["password"] = passwordHeader
	req.Header.Del("password")

	req.Header.Set("Content-Type", a.config.NewContentType)
	jsonBody, err := json.Marshal(newBody)

	if err != nil {
		a.log(err.Error())
		http.Error(rw, "can't make json body", http.StatusInternalServerError)
		return
	}

	req.Body = io.NopCloser(strings.NewReader(string(jsonBody)))
	req.ContentLength = int64(len(jsonBody))

	a.next.ServeHTTP(rw, req)
}
