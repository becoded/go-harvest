package harvest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

const (
	baseURLPath = "/v2"
)

// setup sets up a test HTTP server along with a harvest.APIClient that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(t *testing.T) (client *harvest.APIClient, mux *http.ServeMux, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	loc, err := time.LoadLocation("UTC")
	assert.NoError(t, err)
	// handle err
	time.Local = loc // -> this is setting the global timezone

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Harvest client being tested and is
	// configured to use test server.
	client = harvest.NewAPIClient(nil)

	url, err := url.Parse(server.URL + baseURLPath + "/")
	assert.NoError(t, err)

	client.BaseURL = url
	client.AccountID = "test-account-id"

	return client, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, response %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	err := r.ParseForm()
	assert.NoError(t, err)

	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, response %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) { //nolint: deadcode,unused
	got := r.Header.Get(header)
	assert.Equal(t, want, got)
}

func testURLParseError(t *testing.T, err error) {
	assert.Error(t, err, "Expected error to be returned")

	var e *url.Error
	if !errors.As(err, &e) || e.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testBody(t *testing.T, r *http.Request, path string) {
	want, err := os.ReadFile(filepath.Join("..", "testdata", path))
	assert.NoError(t, err)
	b, err := io.ReadAll(r.Body)
	assert.NoError(t, err, "error reading request body")

	assert.Equal(t, string(want), string(b))
}

func testWriteResponse(t *testing.T, w http.ResponseWriter, path string) {
	response, err := os.ReadFile(filepath.Join("..", "testdata", path))
	assert.NoError(t, err)
	_, err = fmt.Fprint(w, string(response))
	assert.NoError(t, err)
}

// Helper function to test that a value is marshalled to JSON as expected.
func testJSONMarshal(t *testing.T, v interface{}, want string) { //nolint: deadcode,unused
	j, err := json.Marshal(v)
	assert.NoError(t, err, "unable to marshal JSON")

	w := new(bytes.Buffer)
	err = json.Compact(w, []byte(want))
	assert.NoError(t, err, "string is not valid json: %s", want)

	if w.String() != string(j) {
		t.Errorf("json.Marshal(%q) returned %s, response %s", v, j, w)
	}

	// now go the other direction and make sure things unmarshal as expected
	u := reflect.ValueOf(v).Interface()
	err = json.Unmarshal([]byte(want), u)
	assert.NoError(t, err, "unable to unmarshal JSON")

	if !reflect.DeepEqual(v, u) {
		t.Errorf("json.Unmarshal(%q) returned %s, response %s", want, u, v)
	}
}

func TestNewHarvestClient(t *testing.T) {
	t.Parallel()

	c := harvest.NewAPIClient(nil)

	assert.Equal(t, harvest.DefaultBaseURL, c.BaseURL.String())
	assert.Equal(t, harvest.UserAgent, c.UserAgent)

	c2 := harvest.NewAPIClient(nil)

	assert.NotSame(t, c, c2)
}

func TestNewRequest_invalidJSON(t *testing.T) {
	t.Parallel()

	c := harvest.NewAPIClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}

	ctx := context.Background()
	_, err := c.NewRequest(ctx, "GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}

	var e *json.UnsupportedTypeError
	if !errors.As(err, &e) {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	t.Parallel()

	c := harvest.NewAPIClient(nil)
	ctx := context.Background()
	_, err := c.NewRequest(ctx, "GET", ":", nil)
	testURLParseError(t, err)
}

// ensure that no User-Agent header is set if the client's UserAgent is empty.
// This caused a problem with Google's internal http client.
func TestNewRequest_emptyUserAgent(t *testing.T) {
	t.Parallel()

	c := harvest.NewAPIClient(nil)
	c.UserAgent = ""
	ctx := context.Background()

	req, err := c.NewRequest(ctx, "GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}

	if _, ok := req.Header["User-Agent"]; ok {
		t.Fatal("constructed request contains unexpected User-Agent header")
	}
}

// If a nil body is passed to harvest.NewRequest, make sure that nil is also
// passed to http.NewRequest. In most cases, passing an io.Reader that returns
// no content is fine, since there is no difference between an HTTP request
// body that is an empty string versus one that is not set at all. However in
// certain cases, intermediate systems may treat these differently resulting in
// subtle errors.
func TestNewRequest_emptyBody(t *testing.T) {
	t.Parallel()

	c := harvest.NewAPIClient(nil)
	ctx := context.Background()

	req, err := c.NewRequest(ctx, "GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}

	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestDo(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	ctx := context.Background()

	req, _ := client.NewRequest(ctx, "GET", ".", nil)
	body := new(foo)

	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Error(err)
	}

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, response %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "GET", ".", nil)
	resp, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}
