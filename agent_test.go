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

func TestAgent_List(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/agent/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/agent",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/agent/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/agent",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []slide.Agent{}

	ctx := context.Background()
	if err := testService.Agents().List(ctx,
		func(response slide.ListResponse[slide.Agent]) error {
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

func TestAgent_Update(t *testing.T) {
	agentID := "a_0123456789ab"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/agent/update_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPatch,
							Path:   "/v1/agent/" + agentID,
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/agent/update_200.json")
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

	expected := slide.Agent{
		AgentID:             agentID,
		AgentVersion:        "1.2.3",
		BootedAt:            generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		ClientID:            "string",
		DeviceID:            "d_0123456789ab",
		DisplayName:         "My New Displayname",
		EncryptionAlgorithm: "aes-256-gcm",
		FirmwareType:        "UEFI",
		Hostname:            "my-hostname-1",
		LastSeenAt:          generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		Manufacturer:        "Microsoft Corporation",
		OS:                  "windows",
		OSVersion:           "10.0.19042",
		Platform:            "Microsoft Windows 10 Home",
		PublicIPAddress:     "74.83.124.111",
		Addresses: []slide.Address{
			{
				IPs: []string{
					"192.168.1.104",
				},
				MAC: "62:bb:d3:0d:db:7d",
			},
		},
	}

	ctx := context.Background()

	actual, err := testService.Agents().Update(ctx, agentID, "My New Displayname")
	if err != nil {
		t.Fatal(err)
	}

	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actualBytes, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expectedBytes, actualBytes) {
		t.Fatalf("expected did not match actual result: expected: %v, actual: %v", expected, actual)
	}
}

func TestAgent_Get(t *testing.T) {
	agentID := "a_0123456789ab"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/agent/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/agent/" + agentID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	expected := slide.Agent{
		AgentID:             agentID,
		AgentVersion:        "1.2.3",
		BootedAt:            generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		ClientID:            "string",
		DeviceID:            "d_0123456789ab",
		DisplayName:         "My First Device",
		EncryptionAlgorithm: "aes-256-gcm",
		FirmwareType:        "UEFI",
		Hostname:            "my-hostname-1",
		LastSeenAt:          generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		Manufacturer:        "Microsoft Corporation",
		OS:                  "windows",
		OSVersion:           "10.0.19042",
		Platform:            "Microsoft Windows 10 Home",
		PublicIPAddress:     "74.83.124.111",
		Addresses: []slide.Address{
			{
				IPs: []string{
					"192.168.1.104",
				},
				MAC: "62:bb:d3:0d:db:7d",
			},
		},
	}

	ctx := context.Background()

	actual, err := testService.Agents().Get(ctx, agentID)
	if err != nil {
		t.Fatal(err)
	}

	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actualBytes, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expectedBytes, actualBytes) {
		t.Fatalf("expected did not match actual result: expected: %v, actual: %v", expected, actual)
	}
}

func TestAgent_AutoPair(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusCreated,
							FilePath:   "testdata/responses/agent/auto_pair_201.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPost,
							Path:   "/v1/agent",
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/agent/auto_pair_201.json")
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

	expected := slide.AgentAutoPairResponse{
		AgentID:     "a_0123456789ab",
		DisplayName: "My Agent",
		PairCode:    "ABC123",
	}

	ctx := context.Background()
	actual, err := testService.Agents().AutoPair(ctx, slide.AgentAutoPairPayload{
		DeviceID:    "d_0123456789ab",
		DisplayName: "My Agent",
	})

	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}

func TestAgent_Pair(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/agent/pair_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPost,
							Path:   "/v1/agent",
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/agent/pair_200.json")
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

	expected := slide.Agent{
		AgentID:             "a_0123456789ab",
		AgentVersion:        "1.2.3",
		BootedAt:            generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		ClientID:            "string",
		DeviceID:            "d_0123456789ab",
		DisplayName:         "My First Device",
		EncryptionAlgorithm: "aes-256-gcm",
		FirmwareType:        "UEFI",
		Hostname:            "my-hostname-1",
		LastSeenAt:          generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		Manufacturer:        "Microsoft Corporation",
		OS:                  "windows",
		OSVersion:           "10.0.19042",
		Platform:            "Microsoft Windows 10 Home",
		PublicIPAddress:     "74.83.124.111",
		Addresses: []slide.Address{
			{
				IPs: []string{
					"192.168.1.104",
				},
				MAC: "62:bb:d3:0d:db:7d",
			},
		},
	}

	ctx := context.Background()

	actual, err := testService.Agents().Pair(ctx, slide.AgentPairPayload{
		DeviceID: "d_0123456789ab",
		PairCode: "ABC123",
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actualBytes, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expectedBytes, actualBytes) {
		t.Fatalf("expected did not match actual result: expected: %v, actual: %v", expected, actual)
	}
}
