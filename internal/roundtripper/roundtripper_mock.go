package roundtripper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type MockNetwork struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (m MockNetwork) RoundTrip(request *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(request)
}

type MockRoundTripFunc func(request *http.Request) (*http.Response, error)

func MockNetworkQueue(queue []MockRoundTripFunc) http.RoundTripper {
	runNumber := 0

	return MockNetwork{
		RoundTripFunc: func(r *http.Request) (*http.Response, error) {
			defer func() {
				runNumber++
			}()

			if len(queue) <= runNumber {
				return nil, errors.New("empty queue")
			}

			return queue[runNumber](r)
		},
	}
}

type MockResponse interface {
	CreateResponse() (*http.Response, error)
}

type MockResponseFile struct {
	StatusCode        int
	FilePath          string
	ResponseModifiers []ResponseModifier
}

func (f *MockResponseFile) CreateResponse() (*http.Response, error) {
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

type MockResponseNoContent struct {
	StatusCode int
}

func (f *MockResponseNoContent) CreateResponse() (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}

func Serve(r MockResponse) MockRoundTripFunc {
	return func(request *http.Request) (*http.Response, error) {
		return r.CreateResponse()
	}
}
