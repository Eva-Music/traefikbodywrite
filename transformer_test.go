package traefikbodytransform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	cfg := CreateConfig()
	cfg.NewContentType = "application/x-www-form-urlencoded"
	cfg.NewBodyValues = map[int]map[string]string{
		1: {"grant_type": "password"},
		2: {"client_id": "id"},
		3: {"client_secret": "secret"},
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "transformer-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:80", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	req.Header["username"] = []string{"user"}
	req.Header["password"] = []string{"password"}

	newBody := make(map[string]string)

	//create new body from newBodyValues
	for _, m := range cfg.NewBodyValues {
		for k, v := range m {
			newBody[k] = v
		}
	}

	handler.ServeHTTP(recorder, req)

	fmt.Print(newBody)

	//assertJSONBody(t, req, map[string]string{"data": "RAWSTRING"})
	//assertHeader(t, req, "Authorization", "Bearer TOKENDATA")
	//assertHeader(t, req, "Content-Type", "application/json")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}

func assertJSONBody(t *testing.T, req *http.Request, expected map[string]string) {
	t.Helper()

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	jsonBody := make(map[string]string)
	err = json.Unmarshal(reqBody, &jsonBody)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(jsonBody, expected) {
		t.Errorf("invalid json body: %s", jsonBody)
	}
}
