package slide_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/equalsgibson/slide"
	"github.com/equalsgibson/slide/internal/roundtripper"
)

func TestAlert_List(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/alert/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/alert",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/alert/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/alert",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []slide.Alert{}

	ctx := context.Background()
	if err := testService.Alerts().List(ctx,
		func(response slide.ListResponse[slide.Alert]) error {
			actual = append(actual, response.Data...)

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 2 {
		t.Fatal(actual)
	}
}

func TestAlert_Update(t *testing.T) {
	alertID := "al_0123456789ab"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/alert/update_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPatch,
							Path:   "/v1/alert/" + alertID,
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/alert/update_200.json")
								if err != nil {
									return fmt.Errorf("error during test setup - could not read file: %w", err)
								}

								actualBody, err := io.ReadAll(r.Body)
								if err != nil {
									return fmt.Errorf("error during test setup - could not read request body: %w", err)
								}
								r.Body = io.NopCloser(bytes.NewBuffer(actualBody))

								var actualBodyFormatted bytes.Buffer
								if err := json.Indent(&actualBodyFormatted, actualBody, "", "    "); err != nil {
									return fmt.Errorf("error during test setup - could not format request body: %w", err)
								}

								if !bytes.Equal(expectedBody, actualBodyFormatted.Bytes()) {
									return fmt.Errorf("request body does not match expected request format - expected: %v, actual: %v", string(expectedBody), actualBodyFormatted.String())
								}

								return nil
							},
						},
					),
				},
			),
		),
	)

	expected := slide.Alert{
		AgentID:     "a_0123456789ab",
		AlertFields: "string",
		AlertID:     alertID,
		AlertType:   "device_not_checking_in",
		CreatedAt:   "2024-08-23T01:25:08Z",
		DeviceID:    "d_0123456789ab",
		Resolved:    true,
		ResolvedAt:  "2024-08-23T01:25:08Z",
		ResolvedBy:  "John Smith",
	}

	ctx := context.Background()

	actual, err := testService.Alerts().Update(ctx, alertID, true)
	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}

}

func TestAlert_Get(t *testing.T) {
	alertID := "al_0123456789ab"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/alert/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/alert/" + alertID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.Alerts().Get(ctx, alertID)
	if err != nil {
		t.Fatal(err)
	}

	expected := slide.Alert{
		AgentID:     "a_0123456789ab",
		AlertFields: "string",
		AlertID:     alertID,
		AlertType:   "device_not_checking_in",
		CreatedAt:   "2024-08-23T01:25:08Z",
		DeviceID:    "d_0123456789ab",
		Resolved:    false,
		ResolvedAt:  "2024-08-23T01:25:08Z",
		ResolvedBy:  "John Smith",
	}

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}
