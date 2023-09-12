// Package traefikbodywrite plugin.
package traefikbodywrite

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	NewBodyValues map[int]map[string]string `json:"newBodyValues,omitempty"`
	//AddCurrentBodyToNew bool                      `json:"addCurrentBodyToNew,omitempty"`
	NewContentType string `json:"newContentType,omitempty"`
	//TransformerQueryParameterName         string `json:"transformerQueryParameterName,omitempty"`
	//JSONTransformFieldName                string `json:"jsonTransformFieldName,omitempty"`
	//TokenTransformQueryParameterFieldName string `json:"tokenTransformQueryParameterFieldName,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		//AddCurrentBodyToNew: false,
		NewBodyValues:  make(map[int]map[string]string),
		NewContentType: "application/x-www-form-urlencoded",
		//TransformerQueryParameterName:         "transformer",
		//JSONTransformFieldName:                "data",
		//TokenTransformQueryParameterFieldName: "token",
	}
}

// transformer plugin.
type transformer struct {
	next          http.Handler
	newBodyValues map[int]map[string]string
	//addCurrentBodyToNew bool
	newContentType string
	//transformerQueryParameterName         string
	//jsonTransformFieldName                string
	//tokenTransformQueryParameterFieldName string
}

// New created a new transformer plugin.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	return &transformer{
		next:          next,
		newBodyValues: config.NewBodyValues,
		//addCurrentBodyToNew: config.AddCurrentBodyToNew,
		newContentType: config.NewContentType,
		//transformerQueryParameterName:         config.TransformerQueryParameterName,
		//jsonTransformFieldName:                config.JSONTransformFieldName,
		//tokenTransformQueryParameterFieldName: config.TokenTransformQueryParameterFieldName,
	}, nil
}

func (a *transformer) log(format string) {
	_, writeLogError := os.Stderr.WriteString(format)
	if writeLogError != nil {
		panic(writeLogError.Error())
	}
}

//type PayloadInput struct {
//	Host       string                 `json:"host"`
//	Method     string                 `json:"method"`
//	Path       []string               `json:"path"`
//	Parameters url.Values             `json:"parameters"`
//	Headers    map[string][]string    `json:"headers"`
//	Body       map[string]interface{} `json:"body,omitempty"`
//	Form       url.Values             `json:"form,omitempty"`
//}

//	1:
//		"dfg":"sdfg"
//	2:
//		"dfg":"dfvg"

// - достать новые значения для body
// - узнать нужно ли достать старые значения из запроса из body
// - добавить нужный content-type
func (a *transformer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	newBody := make(map[string]string)

	//create new body from newBodyValues
	for _, m := range a.newBodyValues {
		for k, v := range m {
			newBody[k] = v
		}
	}

	//transformerOption := make(map[string]bool)

	// add headers to new body
	usernameHeader, okUser := req.Header["username"]
	if !okUser {
		http.Error(rw, "username header missing", http.StatusInternalServerError)
	}
	newBody["username"] = usernameHeader[0]

	passwordHeader, okPass := req.Header["password"]
	if !okPass {
		http.Error(rw, "password header missing", http.StatusInternalServerError)
	}
	newBody["password"] = passwordHeader[0]

	req.Header.Set("Content-Type", a.newContentType)

	//reqBody, err := io.ReadAll(req.Body)
	//if err != nil {
	//	a.log(err.Error())
	//
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	jsonBody, err := json.Marshal(newBody)

	if err != nil {
		a.log(err.Error())
		http.Error(rw, "can't make json body", http.StatusInternalServerError)
		return
	}

	req.Body = io.NopCloser(strings.NewReader(string(jsonBody)))
	req.ContentLength = int64(len(jsonBody))

	//reqBody, err := io.ReadAll(req.Body)
	//
	//if err != nil {
	//	a.log(err.Error())
	//
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//reqBodyJson, err := json.Marshal(reqBody)
	//for k, v := range reqBodyJson {
	//	newBody[string(k)] = string(v)
	//}

	//}

	//if param := req.URL.Query().Get(a.transformerQueryParameterName); len(param) > 0 {
	//	for _, opt := range strings.Split(strings.ToLower(param), "|") {
	//		transformerOption[opt] = true
	//	}
	//}
	//
	//if transformerOption["body"] {
	//	reqBody, err := io.ReadAll(req.Body)
	//	if err != nil {
	//		a.log(err.Error())
	//
	//		http.Error(rw, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	jsonBody, err := json.Marshal(map[string]string{
	//		a.jsonTransformFieldName: string(reqBody),
	//	})
	//	if err != nil {
	//		a.log(err.Error())
	//
	//		http.Error(rw, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	req.Body = io.NopCloser(strings.NewReader(string(jsonBody)))
	//	req.ContentLength = int64(len(jsonBody))
	//}
	//if transformerOption["json"] {
	//	req.Header.Set("Content-Type", "application/json")
	//}
	//if transformerOption["bearer"] {
	//	token := req.URL.Query().Get(a.tokenTransformQueryParameterFieldName)
	//	req.Header.Set("Authorization", "Bearer "+token)
	//}

	a.next.ServeHTTP(rw, req)
}
