package harvest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"
)

const (
	LibraryVersion   = "1"
	DefaultBaseURL   = "https://api.harvestapp.com/v2/"
	UserAgent        = "becoded/go-harvest/v" + LibraryVersion
	DefaultMediaType = "application/json"
	baseDecimal      = 10
	bitSize64        = 64
)

// A APIClient manages communication with the Harvest API.
type APIClient struct {
	httpClient *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public Harvest API.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	AccountID string

	// User agent used when communicating with the Harvest API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Harvest API.
	Client    *ClientService
	Company   *CompanyService
	Estimate  *EstimateService
	Expense   *ExpenseService
	Invoice   *InvoiceService
	Project   *ProjectService
	Role      *RoleService
	Task      *TaskService
	Timesheet *TimesheetService
	User      *UserService
}

type service struct {
	client *APIClient
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	// The number of records to return per page. Can range between 1 and 100. (Default: 100)
	PerPage int `url:"per_page,omitempty"`
}

type Pagination struct {
	PerPage      *int       `json:"per_page,omitempty"`
	TotalPages   *int       `json:"total_pages,omitempty"`
	TotalEntries *int       `json:"total_entries,omitempty"`
	NextPage     *int       `json:"next_page,omitempty"`
	PreviousPage *int       `json:"previous_page,omitempty"`
	Page         *int       `json:"page,omitempty"`
	Links        *PageLinks `json:"links,omitempty"`
}

type PageLinks struct {
	First    *string `json:"first,omitempty"`
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
	Last     *string `json:"last,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()

	return u.String(), nil
}

// NewAPIClient returns a new Harvest API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewAPIClient(httpClient *http.Client) *APIClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &APIClient{httpClient: httpClient, BaseURL: baseURL, UserAgent: UserAgent}
	c.common.client = c
	c.Client = (*ClientService)(&c.common)
	c.Company = (*CompanyService)(&c.common)
	c.Estimate = (*EstimateService)(&c.common)
	c.Expense = (*ExpenseService)(&c.common)
	c.Invoice = (*InvoiceService)(&c.common)
	c.Project = (*ProjectService)(&c.common)
	c.Role = (*RoleService)(&c.common)
	c.Task = (*TaskService)(&c.common)
	c.Timesheet = (*TimesheetService)(&c.common)
	c.User = (*UserService)(&c.common)

	return c
}

var ErrBaseURLMissingSlash = errors.New("BaseURL must have a trailing slash")

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *APIClient) NewRequest(
	ctx context.Context,
	method,
	urlStr string,
	body interface{},
) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, ErrBaseURLMissingSlash
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", DefaultMediaType)
	}
	// req.Header.Set("Accept", mediaTypeV3)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	if c.AccountID != "" {
		req.Header.Set("Harvest-Account-ID", c.AccountID)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it. If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *APIClient) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		var e *url.Error
		if errors.As(err, &e) {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()

				return nil, e
			}
		}

		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		drainBytes := 512
		if _, err := io.CopyN(io.Discard, resp.Body, int64(drainBytes)); err != nil && !errors.Is(err, io.EOF) {
			logrus.Error(err)
		}

		if err := resp.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	if err := CheckResponse(resp); err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			if _, err := io.Copy(w, resp.Body); err != nil {
				return resp, err
			}

			return resp, err
		}

		if err = json.NewDecoder(resp.Body).Decode(v); errors.Is(err, io.EOF) {
			err = nil // ignore EOF errors caused by empty response body
		}
	}

	return resp, err
}

/*
An ErrorResponse reports one or more errors caused by an API request.
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Errors   []Error        `json:"errors"`  // more detail on individual errors

	Block *struct {
		Reason string `json:"reason,omitempty"`
	} `json:"block,omitempty"`

	DocumentationURL string `json:"documentation_url,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message, r.Errors)
}

// RateLimitError occurs when Harvest returns 429 Forbidden response with a rate limit
// remaining value of 0, and error message starts with "API rate limit exceeded for ".
type RateLimitError struct {
	Rate     Rate           // Rate specifies last known rate limit for the client
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message, "rate limit")
}

// AbuseRateLimitError occurs when Harvest returns 429 Too many requests response with the
// "documentation_url" field value equal to
// "https://help.getharvest.com/api-v2/introduction/overview/general/#rate-limiting".
type AbuseRateLimitError struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message

	// RetryAfter is provided with some abuse rate limit errors. If present,
	// it is the amount of time that the client should wait before retrying.
	// Otherwise, the client should try again later (after an unspecified amount of time).
	RetryAfter *time.Duration
}

func (r *AbuseRateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message)
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}

	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}

	return uri
}

/*
An Error reports more details on an individual error in an ErrorResponse.
*/
type Error struct {
	// resource on which the error occurred
	Resource string `json:"resource"`
	// field on which the error occurred
	Field string `json:"field"`
	// validation error code
	Code string `json:"code"`
	// Message describing the error. Errors with Code == "custom" will always have this set.
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error caused by %v field on %v resource",
		e.Code, e.Field, e.Resource)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}

	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		if err := json.Unmarshal(data, errorResponse); err != nil {
			logrus.Error(err)
		}
	}

	switch r.StatusCode {
	case http.StatusTooManyRequests:
		abuseRateLimitError := &AbuseRateLimitError{
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
		}

		if v := r.Header["Retry-After"]; len(v) > 0 {
			// The "Retry-After" header value will be
			// an integer which represents the number of seconds that one should
			// wait before resuming making requests.
			retryAfterSeconds, _ := strconv.ParseInt(v[0], baseDecimal, bitSize64) // Error handling is noop.
			retryAfter := time.Duration(retryAfterSeconds) * time.Second
			abuseRateLimitError.RetryAfter = &retryAfter
		}

		return abuseRateLimitError
	default:
		return errorResponse
	}
}

// Rate represents the rate limit for the current client.
type Rate struct {
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`
}

func (r Rate) String() string {
	return Stringify(r)
}

// RateLimits represents the rate limits for the current client.
type RateLimits struct {
	// 100 requests per 15 seconds
	// Harvest API docs: https://help.getharvest.com/api-v2/introduction/overview/general/#rate-limiting
	Core *Rate `json:"core"`
}

func (r RateLimits) String() string {
	return Stringify(r)
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Ints64 is a helper routine that allocates a new int slice
// to store v and returns a pointer to it.
func Ints64(v []int64) *[]int64 { return &v }

// Float64 is a helper routine that allocates a new float64 value
// to store v and returns a pointer to it.
func Float64(v float64) *float64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// DateP is a helper routine that allocates a new Date value
// to store v and returns a pointer to it.
func DateP(v Date) *Date { return &v }

// TimeTimeP is a helper routine that allocates a new time.Time value
// to store v and returns a pointer to it.
func TimeTimeP(v time.Time) *time.Time { return &v }

// TimeP is a helper routine that allocates a new Time value
// to store v and returns a pointer to it.
func TimeP(v Time) *Time { return &v }
