package goslide

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type requestClient struct {
	token      *oauth2.Token
	apiURL     string
	httpClient *http.Client
}

func (rc *requestClient) do(request *http.Request, target any) error {
	if request.URL.Host == "" {
		request.URL.Host = rc.apiURL
	}

	request.URL.Scheme = "https"

	request.Header.Set("Accept", "application/json")

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := rc.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		slideAPIError := &SlideError{
			HTTPStatusCode:    response.StatusCode,
			HTTPRequestPath:   request.URL.Path,
			HTTPRequestMethod: request.Method,
		}
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bodyBytes, slideAPIError); err != nil {
			return fmt.Errorf("goslide library error while unmarshalling Slide API response - %w. HTTP Status Code %d", err, response.StatusCode)
		}

		return slideAPIError
	}

	if target != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return err
		}
	}

	return nil
}

func (rc *requestClient) SlideRequest(request *http.Request, target any) error {
	if request.Header.Get("Authorization") == "" {
		if rc.token == nil {
			return fmt.Errorf("unable to set authorization header - API Token not set")
		}
		rc.token.SetAuthHeader(request)
	}

	return rc.do(request, target)
}

func (rc *requestClient) Request(request *http.Request, target any) error {
	return rc.do(request, target)
}
