package goslide_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"testing"

	"github.com/equalsgibson/goslide"
	"github.com/equalsgibson/goslide/internal/roundtripper"
)

func TestHealth_IsAuthenticated_InvalidToken(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusUnauthorized,
							FilePath:   "testdata/responses/health/is_authenticated_invalid_auth_401.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/user",
							Query: url.Values{
								"limit": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	authenticated, err := testService.CheckAuthenticationToken(ctx)

	// Confirm that the error is a slide pkg error and not a net/http or other error
	var slideError *goslide.SlideError
	if errors.As(err, &slideError) {
		// Validate the StatusCode is being set on the error correctly.
		if slideError.HTTPStatusCode != http.StatusUnauthorized {
			t.Fatalf("expected to receive HTTP code: %d, actual HTTP code: %d", http.StatusUnauthorized, slideError.HTTPStatusCode)
		}

		// Validate the error code enum for unauthorized is in the slice of returned error codes.
		if !slices.Contains(slideError.Codes, goslide.APIErrorCode_ERR_UNAUTHORIZED) {
			t.Fatalf("the Slide API Error did not contain the correct error code - expecting %s, actual codes: %+v", goslide.APIErrorCode_ERR_UNAUTHORIZED, slideError.Codes)
		}
	} else {
		t.Fatal("expected to receive a Slide API error")
	}

	// Validate that we are returning false, as we received http.StatusUnauthorized
	if authenticated {
		t.Fatal("expected to not be authenticated due to bad token")
	}
}

func TestHealth_IsAuthenticated_MissingToken(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusUnauthorized,
							FilePath:   "testdata/responses/health/is_authenticated_missing_auth_401.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/user",
							Query: url.Values{
								"limit": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	authenticated, err := testService.CheckAuthenticationToken(ctx)

	// Confirm that the error is a slide pkg error and not a net/http or other error
	var slideError *goslide.SlideError
	if errors.As(err, &slideError) {
		// Validate the StatusCode is being set on the error correctly.
		if slideError.HTTPStatusCode != http.StatusUnauthorized {
			t.Fatalf("expected to receive HTTP code: %d, actual HTTP code: %d", http.StatusUnauthorized, slideError.HTTPStatusCode)
		}

		// Validate the error code enum for missing auth token is in the slice of returned error codes.
		if !slices.Contains(slideError.Codes, goslide.APIErrorCode_ERR_MISSING_AUTHENTICATION) {
			t.Fatalf("the Slide API Error did not contain the correct error code - expecting %s, actual codes: %+v", goslide.APIErrorCode_ERR_MISSING_AUTHENTICATION, slideError.Codes)
		}
	} else {
		t.Fatal("expected to receive a Slide API error")
	}

	// Validate that we are returning false, as we received http.StatusUnauthorized
	if authenticated {
		t.Fatal("expected to not be authenticated due to bad token")
	}
}

func TestHealth_IsAuthenticated_OK(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/health/is_authenticated_success_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/user",
							Query: url.Values{
								"limit": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	authenticated, err := testService.CheckAuthenticationToken(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !authenticated {
		t.Fatal("expected to be authenticated")
	}
}
