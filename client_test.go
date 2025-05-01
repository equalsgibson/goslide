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

func TestClient_List(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/client/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/client",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/client/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/client",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []slide.Client{}

	ctx := context.Background()
	if err := testService.Clients().List(ctx,
		func(response slide.ListResponse[slide.Client]) error {
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

func TestClient_Create(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusCreated,
							FilePath:   "testdata/responses/client/create_201.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPost,
							Path:   "/v1/client",
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/client/create_201.json")
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

	expected := slide.Client{
		Name:     "My Client",
		Comments: "",
		ClientID: "c_123456789abc",
	}

	payload := slide.ClientPayload{
		Name: "My Client",
	}

	ctx := context.Background()

	actual, err := testService.Clients().Create(ctx, payload)
	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}

}

func TestClient_Get(t *testing.T) {
	clientID := "c_123456789abc"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/client/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/client/" + clientID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.Clients().Get(ctx, clientID)
	if err != nil {
		t.Fatal(err)
	}

	expected := slide.Client{
		ClientID: clientID,
		Comments: "This is a test client",
		Name:     "Slide Office",
	}

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}

func TestClient_Delete(t *testing.T) {
	clientID := "c_123456789abc"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseNoContent{
							StatusCode: http.StatusNoContent,
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodDelete,
							Path:   "/v1/client/" + clientID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	if err := testService.Clients().Delete(ctx, clientID); err != nil {
		t.Fatal(err)
	}
}

func TestClient_Update(t *testing.T) {
	clientID := "c_123456789abc"
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/client/update_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPatch,
							Path:   "/v1/client/" + clientID,
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/client/update_200.json")
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

	expected := slide.Client{
		Name:     "My Client",
		Comments: "",
		ClientID: clientID,
	}

	payload := slide.ClientPayload{
		Name: "My Client",
	}

	ctx := context.Background()

	actual, err := testService.Clients().Update(ctx, clientID, payload)
	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}
