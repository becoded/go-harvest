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
	"strconv"
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

func testHeader(t *testing.T, r *http.Request, header string, want string) { //nolint: unused
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
func testJSONMarshal(t *testing.T, v interface{}, want string) { //nolint: unused
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

func TestNewAPIClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		httpClient *http.Client
		wantNil    bool
	}{
		{
			name:       "With nil http client",
			httpClient: nil,
			wantNil:    false,
		},
		{
			name:       "With custom http client",
			httpClient: &http.Client{Timeout: 10 * time.Second},
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := harvest.NewAPIClient(tt.httpClient)
			assert.NotNil(t, client)
			assert.Equal(t, harvest.DefaultBaseURL, client.BaseURL.String())
			assert.Equal(t, harvest.UserAgent, client.UserAgent)
			assert.NotNil(t, client.Client)
			assert.NotNil(t, client.Company)
			assert.NotNil(t, client.Estimate)
			assert.NotNil(t, client.Expense)
			assert.NotNil(t, client.Invoice)
			assert.NotNil(t, client.Project)
			assert.NotNil(t, client.Role)
			assert.NotNil(t, client.Task)
			assert.NotNil(t, client.Timesheet)
			assert.NotNil(t, client.User)
		})
	}
}

func TestNewRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		method    string
		urlStr    string
		body      interface{}
		setupFunc func(*harvest.APIClient)
		wantErr   bool
		errType   string
	}{
		{
			name:    "Valid GET request without body",
			method:  "GET",
			urlStr:  "clients",
			body:    nil,
			wantErr: false,
		},
		{
			name:   "Valid POST request with body",
			method: "POST",
			urlStr: "clients",
			body: map[string]string{
				"name": "Test Client",
			},
			wantErr: false,
		},
		{
			name:   "Request with missing trailing slash in BaseURL",
			method: "GET",
			urlStr: "clients",
			body:   nil,
			setupFunc: func(c *harvest.APIClient) {
				u, _ := url.Parse("https://api.harvestapp.com/v2")
				c.BaseURL = u
			},
			wantErr: true,
			errType: "ErrBaseURLMissingSlash",
		},
		{
			name:    "Request with invalid URL",
			method:  "GET",
			urlStr:  ":",
			body:    nil,
			wantErr: true,
			errType: "url.Error",
		},
		{
			name:   "Request with invalid JSON body",
			method: "POST",
			urlStr: "clients",
			body: map[interface{}]interface{}{
				1: "invalid",
			},
			wantErr: true,
			errType: "json.UnsupportedTypeError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := harvest.NewAPIClient(nil)
			if tt.setupFunc != nil {
				tt.setupFunc(client)
			}

			ctx := context.Background()
			req, err := client.NewRequest(ctx, tt.method, tt.urlStr, tt.body)

			if tt.wantErr {
				assert.Error(t, err)

				if tt.errType == "ErrBaseURLMissingSlash" {
					assert.Equal(t, harvest.ErrBaseURLMissingSlash, err)
				}

				return
			}

			// Success case assertions
			assert.NoError(t, err)
			assert.NotNil(t, req)
			assert.Equal(t, tt.method, req.Method)

			if tt.body != nil {
				assert.Equal(t, harvest.DefaultMediaType, req.Header.Get("Content-Type"))
			}

			if client.UserAgent != "" {
				assert.Equal(t, client.UserAgent, req.Header.Get("User-Agent"))
			}

			if client.AccountID != "" {
				assert.Equal(t, client.AccountID, req.Header.Get("Harvest-Account-ID"))
			}
		})
	}
}

func TestCheckResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		body       string
		headers    map[string]string
		wantErr    bool
		errType    string
	}{
		{
			name:       "Success 200",
			statusCode: http.StatusOK,
			body:       `{"id": 1}`,
			wantErr:    false,
		},
		{
			name:       "Success 201",
			statusCode: http.StatusCreated,
			body:       `{"id": 1}`,
			wantErr:    false,
		},
		{
			name:       "Success 204",
			statusCode: http.StatusNoContent,
			body:       "",
			wantErr:    false,
		},
		{
			name:       "Error 400 Bad Request",
			statusCode: http.StatusBadRequest,
			body:       `{"message": "Bad Request", "errors": [{"resource": "client", "field": "name", "code": "required"}]}`,
			wantErr:    true,
			errType:    "ErrorResponse",
		},
		{
			name:       "Error 404 Not Found",
			statusCode: http.StatusNotFound,
			body:       `{"message": "Not Found"}`,
			wantErr:    true,
			errType:    "ErrorResponse",
		},
		{
			name:       "Error 429 Too Many Requests without Retry-After",
			statusCode: http.StatusTooManyRequests,
			body:       `{"message": "Rate limit exceeded"}`,
			wantErr:    true,
			errType:    "AbuseRateLimitError",
		},
		{
			name:       "Error 429 Too Many Requests with Retry-After",
			statusCode: http.StatusTooManyRequests,
			body:       `{"message": "Rate limit exceeded"}`,
			headers:    map[string]string{"Retry-After": "60"},
			wantErr:    true,
			errType:    "AbuseRateLimitError",
		},
		{
			name:       "Error 500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			body:       `{"message": "Internal Server Error"}`,
			wantErr:    true,
			errType:    "ErrorResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
				Request: &http.Request{
					Method: "GET",
					URL:    &url.URL{Path: "/test"},
				},
				Header: make(http.Header),
			}

			for k, v := range tt.headers {
				resp.Header.Set(k, v)
			}

			err := harvest.CheckResponse(resp)

			if tt.wantErr {
				assert.Error(t, err)

				switch tt.errType {
				case "ErrorResponse":
					var e *harvest.ErrorResponse
					assert.True(t, errors.As(err, &e))
				case "AbuseRateLimitError":
					var e *harvest.AbuseRateLimitError
					assert.True(t, errors.As(err, &e))

					if retryAfter, ok := tt.headers["Retry-After"]; ok {
						assert.NotNil(t, e.RetryAfter)

						expectedSeconds, _ := strconv.ParseInt(retryAfter, 10, 64)
						assert.Equal(t, time.Duration(expectedSeconds)*time.Second, *e.RetryAfter)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErrorResponse_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      *harvest.ErrorResponse
		contains []string
	}{
		{
			name: "Error with message and errors",
			err: &harvest.ErrorResponse{
				Response: &http.Response{
					StatusCode: 400,
					Request: &http.Request{
						Method: "POST",
						URL:    &url.URL{Scheme: "https", Host: "api.harvestapp.com", Path: "/v2/clients"},
					},
				},
				Message: "Validation failed",
				Errors: []harvest.Error{
					{
						Resource: "client",
						Field:    "name",
						Code:     "required",
					},
				},
			},
			contains: []string{"POST", "400", "Validation failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			errStr := tt.err.Error()
			for _, substr := range tt.contains {
				assert.Contains(t, errStr, substr)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      *harvest.Error
		contains []string
	}{
		{
			name: "Validation error",
			err: &harvest.Error{
				Resource: "client",
				Field:    "name",
				Code:     "required",
				Message:  "Name is required",
			},
			contains: []string{"required", "name", "client"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			errStr := tt.err.Error()
			for _, substr := range tt.contains {
				assert.Contains(t, errStr, substr)
			}
		})
	}
}

func TestRateLimitError_Error(t *testing.T) {
	t.Parallel()

	err := &harvest.RateLimitError{
		Response: &http.Response{
			StatusCode: 429,
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Scheme: "https", Host: "api.harvestapp.com", Path: "/v2/clients"},
			},
		},
		Message: "Rate limit exceeded",
		Rate: harvest.Rate{
			Limit:     100,
			Remaining: 0,
		},
	}

	errStr := err.Error()
	assert.Contains(t, errStr, "GET")
	assert.Contains(t, errStr, "429")
	assert.Contains(t, errStr, "rate limit")
}

func TestAbuseRateLimitError_Error(t *testing.T) {
	t.Parallel()

	retryAfter := 60 * time.Second
	err := &harvest.AbuseRateLimitError{
		Response: &http.Response{
			StatusCode: 429,
			Request: &http.Request{
				Method: "POST",
				URL:    &url.URL{Scheme: "https", Host: "api.harvestapp.com", Path: "/v2/time_entries"},
			},
		},
		Message:    "Too many requests",
		RetryAfter: &retryAfter,
	}

	errStr := err.Error()
	assert.Contains(t, errStr, "POST")
	assert.Contains(t, errStr, "429")
	assert.Contains(t, errStr, "Too many requests")
}

func TestHelperFunctions(t *testing.T) {
	t.Parallel()

	t.Run("Bool", func(t *testing.T) {
		t.Parallel()

		v := harvest.Bool(true)
		assert.NotNil(t, v)
		assert.Equal(t, true, *v)

		v2 := harvest.Bool(false)
		assert.NotNil(t, v2)
		assert.Equal(t, false, *v2)
	})

	t.Run("Int", func(t *testing.T) {
		t.Parallel()

		v := harvest.Int(42)
		assert.NotNil(t, v)
		assert.Equal(t, 42, *v)

		v2 := harvest.Int(0)
		assert.NotNil(t, v2)
		assert.Equal(t, 0, *v2)
	})

	t.Run("Int64", func(t *testing.T) {
		t.Parallel()

		v := harvest.Int64(int64(123456789))
		assert.NotNil(t, v)
		assert.Equal(t, int64(123456789), *v)

		v2 := harvest.Int64(int64(0))
		assert.NotNil(t, v2)
		assert.Equal(t, int64(0), *v2)
	})

	t.Run("Ints64", func(t *testing.T) {
		t.Parallel()

		slice := []int64{1, 2, 3}
		v := harvest.Ints64(slice)
		assert.NotNil(t, v)
		assert.Equal(t, slice, *v)

		emptySlice := []int64{}
		v2 := harvest.Ints64(emptySlice)
		assert.NotNil(t, v2)
		assert.Equal(t, emptySlice, *v2)
	})

	t.Run("Float64", func(t *testing.T) {
		t.Parallel()

		v := harvest.Float64(3.14)
		assert.NotNil(t, v)
		assert.Equal(t, 3.14, *v)

		v2 := harvest.Float64(0.0)
		assert.NotNil(t, v2)
		assert.Equal(t, 0.0, *v2)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		v := harvest.String("test")
		assert.NotNil(t, v)
		assert.Equal(t, "test", *v)

		v2 := harvest.String("")
		assert.NotNil(t, v2)
		assert.Equal(t, "", *v2)
	})

	t.Run("TimeTimeP", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		v := harvest.TimeTimeP(now)
		assert.NotNil(t, v)
		assert.Equal(t, now, *v)
	})

	t.Run("DateP", func(t *testing.T) {
		t.Parallel()

		date := harvest.Date{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}
		v := harvest.DateP(date)
		assert.NotNil(t, v)
		assert.Equal(t, date, *v)
	})

	t.Run("TimeP", func(t *testing.T) {
		t.Parallel()

		timeVal := harvest.Time{Time: time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC)}
		v := harvest.TimeP(timeVal)
		assert.NotNil(t, v)
		assert.Equal(t, timeVal, *v)
	})
}

func TestRate_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rate harvest.Rate
		want string
	}{
		{
			name: "Rate with values",
			rate: harvest.Rate{
				Limit:     100,
				Remaining: 75,
			},
			want: "harvest.Rate{Limit:100, Remaining:75}",
		},
		{
			name: "Rate with zero values",
			rate: harvest.Rate{
				Limit:     0,
				Remaining: 0,
			},
			want: "harvest.Rate{Limit:0, Remaining:0}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.rate.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRateLimits_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rateLimits harvest.RateLimits
		want       string
	}{
		{
			name: "RateLimits with core rate",
			rateLimits: harvest.RateLimits{
				Core: &harvest.Rate{
					Limit:     100,
					Remaining: 50,
				},
			},
			want: "harvest.RateLimits{Core:harvest.Rate{Limit:100, Remaining:50}}",
		},
		{
			name:       "RateLimits with nil core",
			rateLimits: harvest.RateLimits{},
			want:       "harvest.RateLimits{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.rateLimits.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPagination_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pagination harvest.Pagination
		want       string
	}{
		{
			name: "Pagination with all fields",
			pagination: harvest.Pagination{
				PerPage:      harvest.Int(100),
				TotalPages:   harvest.Int(5),
				TotalEntries: harvest.Int(450),
				NextPage:     harvest.Int(3),
				PreviousPage: harvest.Int(1),
				Page:         harvest.Int(2),
				Links: &harvest.PageLinks{
					First:    harvest.String("https://api.harvestapp.com/v2/clients?page=1"),
					Next:     harvest.String("https://api.harvestapp.com/v2/clients?page=3"),
					Previous: harvest.String("https://api.harvestapp.com/v2/clients?page=1"),
					Last:     harvest.String("https://api.harvestapp.com/v2/clients?page=5"),
				},
			},
			want: `harvest.Pagination{PerPage:100, TotalPages:5, TotalEntries:450, NextPage:3, PreviousPage:1, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/clients?page=1", Next:"https://api.harvestapp.com/v2/clients?page=3", Previous:"https://api.harvestapp.com/v2/clients?page=1", Last:"https://api.harvestapp.com/v2/clients?page=5"}}`, //nolint: lll
		},
		{
			name: "Pagination with minimal fields",
			pagination: harvest.Pagination{
				Page: harvest.Int(1),
			},
			want: "harvest.Pagination{Page:1}",
		},
		{
			name:       "Empty Pagination",
			pagination: harvest.Pagination{},
			want:       "harvest.Pagination{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := harvest.Stringify(tt.pagination)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPageLinks_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pageLinks harvest.PageLinks
		want      string
	}{
		{
			name: "PageLinks with all fields",
			pageLinks: harvest.PageLinks{
				First:    harvest.String("https://api.harvestapp.com/v2/clients?page=1"),
				Next:     harvest.String("https://api.harvestapp.com/v2/clients?page=3"),
				Previous: harvest.String("https://api.harvestapp.com/v2/clients?page=1"),
				Last:     harvest.String("https://api.harvestapp.com/v2/clients?page=5"),
			},
			want: `harvest.PageLinks{First:"https://api.harvestapp.com/v2/clients?page=1", Next:"https://api.harvestapp.com/v2/clients?page=3", Previous:"https://api.harvestapp.com/v2/clients?page=1", Last:"https://api.harvestapp.com/v2/clients?page=5"}`, //nolint: lll
		},
		{
			name:      "Empty PageLinks",
			pageLinks: harvest.PageLinks{},
			want:      "harvest.PageLinks{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := harvest.Stringify(tt.pageLinks)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDo_ioWriter(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	content := "test content"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, content)
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "GET", ".", nil)

	var buf bytes.Buffer
	_, err := client.Do(ctx, req, &buf)
	assert.NoError(t, err)
	assert.Equal(t, content, buf.String())
}

func TestDo_emptyResponse(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "DELETE", ".", nil)

	type foo struct {
		A string
	}
	body := new(foo)

	_, err := client.Do(ctx, req, body)
	assert.NoError(t, err)
}

func TestDo_contextCanceled(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// Simulate a slow response
		time.Sleep(100 * time.Millisecond)
		fmt.Fprint(w, `{"A":"a"}`)
	})

	ctx, cancel := context.WithCancel(context.Background())
	req, _ := client.NewRequest(ctx, "GET", ".", nil)

	type foo struct {
		A string
	}
	body := new(foo)

	// Cancel the context before the request completes
	cancel()

	_, err := client.Do(ctx, req, body)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestNewRequest_withAccountID(t *testing.T) {
	t.Parallel()

	client := harvest.NewAPIClient(nil)
	client.AccountID = "test-account-123"

	ctx := context.Background()
	req, err := client.NewRequest(ctx, "GET", "clients", nil)

	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "test-account-123", req.Header.Get("Harvest-Account-ID"))
}

func TestNewRequest_withCustomUserAgent(t *testing.T) {
	t.Parallel()

	client := harvest.NewAPIClient(nil)
	client.UserAgent = "CustomAgent/1.0"

	ctx := context.Background()
	req, err := client.NewRequest(ctx, "GET", "clients", nil)

	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "CustomAgent/1.0", req.Header.Get("User-Agent"))
}

func TestNewRequest_bodyEncoding(t *testing.T) {
	t.Parallel()

	client := harvest.NewAPIClient(nil)

	type testBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	body := &testBody{
		Name:  "Test User",
		Email: "test@example.com",
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, "POST", "users", body)

	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, harvest.DefaultMediaType, req.Header.Get("Content-Type"))

	// Read and verify the body
	bodyBytes, err := io.ReadAll(req.Body)
	assert.NoError(t, err)

	var decodedBody testBody
	err = json.Unmarshal(bodyBytes, &decodedBody)
	assert.NoError(t, err)
	assert.Equal(t, body.Name, decodedBody.Name)
	assert.Equal(t, body.Email, decodedBody.Email)
}

func TestCheckResponse_withBlock(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: http.StatusForbidden,
		Body: io.NopCloser(bytes.NewBufferString(
			`{"message": "Account blocked", "block": {"reason": "Payment required"}}`,
		)),
		Request: &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/test"},
		},
	}

	err := harvest.CheckResponse(resp)
	assert.Error(t, err)

	var e *harvest.ErrorResponse
	assert.True(t, errors.As(err, &e))
	assert.NotNil(t, e.Block)
	assert.Equal(t, "Payment required", e.Block.Reason)
}

func TestDo_nilResponse(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "GET", ".", nil)

	// Test with nil response object
	_, err := client.Do(ctx, req, nil)
	assert.NoError(t, err)
}

func TestDo_errorResponse(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		http.Error(
			w,
			`{"message":"Validation failed","errors":[{"resource":"client","field":"name","code":"required"}]}`,
			http.StatusBadRequest,
		)
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "POST", ".", map[string]string{"name": ""})

	type foo struct {
		A string
	}
	body := new(foo)

	resp, err := client.Do(ctx, req, body)
	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var e *harvest.ErrorResponse
	assert.True(t, errors.As(err, &e))
	assert.Equal(t, "Validation failed", e.Message)
	assert.Len(t, e.Errors, 1)
	assert.Equal(t, "client", e.Errors[0].Resource)
	assert.Equal(t, "name", e.Errors[0].Field)
	assert.Equal(t, "required", e.Errors[0].Code)
}

func TestDo_rateLimitError(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Retry-After", "120")
		http.Error(w, `{"message":"Rate limit exceeded"}`, http.StatusTooManyRequests)
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "GET", ".", nil)

	resp, err := client.Do(ctx, req, nil)
	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	var e *harvest.AbuseRateLimitError
	assert.True(t, errors.As(err, &e))
	assert.Equal(t, "Rate limit exceeded", e.Message)
	assert.NotNil(t, e.RetryAfter)
	assert.Equal(t, 120*time.Second, *e.RetryAfter)
}

func TestSanitizeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputURL string
		want     string
	}{
		{
			name:     "URL with client_secret",
			inputURL: "https://api.example.com/oauth?client_id=123&client_secret=secret123&code=abc",
			want:     "https://api.example.com/oauth?client_id=123&client_secret=REDACTED&code=abc",
		},
		{
			name:     "URL without client_secret",
			inputURL: "https://api.example.com/clients?page=1&per_page=100",
			want:     "https://api.example.com/clients?page=1&per_page=100",
		},
		{
			name:     "URL with empty client_secret",
			inputURL: "https://api.example.com/oauth?client_id=123&client_secret=",
			want:     "https://api.example.com/oauth?client_id=123&client_secret=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create an ErrorResponse to trigger sanitizeURL through Error() method
			u, _ := url.Parse(tt.inputURL)
			resp := &http.Response{
				StatusCode: 400,
				Request: &http.Request{
					Method: "GET",
					URL:    u,
				},
			}

			err := &harvest.ErrorResponse{
				Response: resp,
				Message:  "Test error",
			}

			errStr := err.Error()
			// The error string should contain the sanitized URL
			if tt.name == "URL with client_secret" {
				assert.Contains(t, errStr, "REDACTED")
				assert.NotContains(t, errStr, "secret123")
			} else {
				assert.Contains(t, errStr, tt.inputURL)
			}
		})
	}
}

func TestNewRequest_withBody_escapeHTML(t *testing.T) {
	t.Parallel()

	client := harvest.NewAPIClient(nil)

	type testBody struct {
		HTML string `json:"html"`
	}

	body := &testBody{
		HTML: "<script>alert('test')</script>",
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, "POST", "test", body)

	assert.NoError(t, err)
	assert.NotNil(t, req)

	// Read and verify the body - HTML should NOT be escaped
	bodyBytes, err := io.ReadAll(req.Body)
	assert.NoError(t, err)

	bodyStr := string(bodyBytes)
	// The encoder is set to NOT escape HTML
	assert.Contains(t, bodyStr, "<script>alert('test')</script>")
	assert.NotContains(t, bodyStr, "\\u003c") // Should not contain escaped <
}

func TestListOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts harvest.ListOptions
		want string
	}{
		{
			name: "ListOptions with page and per_page",
			opts: harvest.ListOptions{
				Page:    2,
				PerPage: 50,
			},
			want: "harvest.ListOptions{Page:2, PerPage:50}",
		},
		{
			name: "ListOptions with only page",
			opts: harvest.ListOptions{
				Page: 1,
			},
			want: "harvest.ListOptions{Page:1, PerPage:0}",
		},
		{
			name: "Empty ListOptions",
			opts: harvest.ListOptions{},
			want: "harvest.ListOptions{Page:0, PerPage:0}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := harvest.Stringify(tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}
