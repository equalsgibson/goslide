package roundtripper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestNetwork struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (m TestNetwork) RoundTrip(request *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(request)
}

type TestRoundTripFunc func(t *testing.T, request *http.Request) (*http.Response, error)

func NetworkQueue(t *testing.T, queue []TestRoundTripFunc) http.RoundTripper {
	runNumber := 0

	return TestNetwork{
		RoundTripFunc: func(r *http.Request) (*http.Response, error) {
			defer func() {
				runNumber++
			}()

			if len(queue) <= runNumber {
				return nil, errors.New("empty queue")
			}

			return queue[runNumber](t, r)
		},
	}
}

type ExpectedTestRequest struct {
	Method    string
	Path      string
	Query     url.Values
	Validator func(r *http.Request) error
}

type TestResponse interface {
	CreateResponse() (*http.Response, error)
}

type TestResponseFile struct {
	StatusCode        int
	FilePath          string
	ResponseModifiers []ResponseModifier
}

type ResponseModifier interface {
	ModifyResponse(r *http.Response)
}

type ResponseModifierFunc func(*http.Response)

func (r ResponseModifierFunc) ModifyResponse(response *http.Response) {
	r(response)
}

func (f *TestResponseFile) CreateResponse() (*http.Response, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("response body file not found: %s", f.FilePath)
	}

	response := &http.Response{
		StatusCode: f.StatusCode,
		Body:       io.NopCloser(file),
		Header:     make(http.Header),
	}

	for _, responseModifier := range f.ResponseModifiers {
		responseModifier.ModifyResponse(response)
	}

	headers := response.Header

	// Check if the Content-Type header has been set in the Header map. If not - default to application/json
	if _, ok := headers["Content-Type"]; !ok {
		response.Header.Set("Content-Type", "application/json")
	}

	return response, nil
}

type TestResponseNoContent struct {
	StatusCode int
}

func (f *TestResponseNoContent) CreateResponse() (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}

func ServeAndValidate(t *testing.T, r TestResponse, expected ExpectedTestRequest) TestRoundTripFunc {
	return func(t *testing.T, request *http.Request) (*http.Response, error) {
		if diff := cmp.Diff(expected.Method, request.Method); diff != "" {
			t.Logf("%s Request Method mismatch (-want +got):\n%s", t.Name(), diff)
			t.Fail()
		}

		if diff := cmp.Diff(expected.Path, request.URL.Path); diff != "" {
			t.Logf("%s Request Path mismatch (-want +got):\n%s", t.Name(), diff)
			t.Fail()
		}

		if diff := cmp.Diff(expected.Query, request.URL.Query()); diff != "" {
			t.Errorf("%s URL Parameter mismatch (-want +got):\n%s", t.Name(), diff)
			t.Fail()
		}

		if expected.Validator != nil {
			if err := expected.Validator(request); err != nil {
				t.Logf("validation check on expected request failed - error: %s", err.Error())
				t.Fail()
			}
		}

		return r.CreateResponse()
	}
}
