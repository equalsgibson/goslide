package goslide_test

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

	"github.com/equalsgibson/goslide"
	"github.com/equalsgibson/goslide/internal/roundtripper"
	"github.com/google/go-cmp/cmp"
)

func TestAccount_List(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/account/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/account",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/account/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/account",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []goslide.Account{}

	ctx := context.Background()
	if err := testService.Accounts().List(ctx,
		func(response goslide.ListResponse[goslide.Account]) error {
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

func TestAccount_Update(t *testing.T) {
	accountID := "act_0123456789ab"

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/account/update_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPatch,
							Path:   "/v1/account/" + accountID,
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/account/update_200.json")
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

								if diff := cmp.Diff(string(expectedBody), actualBodyFormatted.String()); diff != "" {
									t.Fatalf("%s Expected Request Body mismatch (-want +got):\n%s", t.Name(), diff)
								}

								return nil
							},
						},
					),
				},
			),
		),
	)

	expected := goslide.Account{
		AccountID:   "act_0123456789ab",
		AccountName: "Slide Inc.",
		AlertEmails: []string{
			"john.doe@gmail.com",
		},
		BillingAddress: goslide.BillingAddress{
			City:       "New York",
			Country:    "US",
			Line1:      "123 Main St",
			Line2:      "Apt 4B",
			PostalCode: "10001",
			State:      "NY",
		},
		PrimaryContact: "John Smith",
		PrimaryEmail:   "john.doe@gmail.com",
		PrimaryPhone:   "+1 555-555-5555",
	}

	ctx := context.Background()

	actual, err := testService.Accounts().Update(ctx, accountID, []string{"john.doe@gmail.com"})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestAccount_Get(t *testing.T) {
	accountID := "act_0123456789ab"

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/account/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/account/" + accountID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.Accounts().Get(ctx, accountID)
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.Account{
		AccountID:   "act_0123456789ab",
		AccountName: "Slide Inc.",
		AlertEmails: []string{
			"john.doe@gmail.com",
			"jane.smith@example.com",
			"user123@domain.net",
		},
		BillingAddress: goslide.BillingAddress{
			City:       "New York",
			Country:    "US",
			Line1:      "123 Main St",
			Line2:      "Apt 4B",
			PostalCode: "10001",
			State:      "NY",
		},
		PrimaryContact: "John Smith",
		PrimaryEmail:   "john.doe@gmail.com",
		PrimaryPhone:   "+1 555-555-5555",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}
