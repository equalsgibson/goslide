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

func TestNetwork_List(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/network/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/network",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/network/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/network",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []slide.Network{}

	ctx := context.Background()
	if err := testService.Networks().List(ctx,
		func(response slide.ListResponse[slide.Network]) error {
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

func TestNetwork_Get(t *testing.T) {
	networkID := "net_012345"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/network/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/network/" + networkID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	expected := slide.Network{
		BridgeDeviceID: "d_0123456789ab",
		ClientID:       "string",
		Comments:       "This is a test network",
		ConnectedVirtIDs: []string{
			"virt_0123456789ab",
		},
		DHCP:           true,
		DHCPRangeEnd:   "10.0.0.200",
		DHCPRangeStart: "10.0.0.100",
		Internet:       true,
		Name:           "Bridge to LAN Network",
		Nameservers:    "1.1.1.1,1.0.0.1",
		NetworkID:      "net_012345",
		RouterPrefix:   "10.0.0.1/24",
		Type:           slide.NetworkTypeDisaster_STANDARD,
	}

	ctx := context.Background()

	actual, err := testService.Networks().Get(ctx, networkID)
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

func TestNetwork_Delete(t *testing.T) {
	networkID := "net_012345"

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
							Path:   "/v1/network/" + networkID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	if err := testService.Networks().Delete(ctx, networkID); err != nil {
		t.Fatal(err)
	}
}

func TestNetwork_Create(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusCreated,
							FilePath:   "testdata/responses/network/create_201.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPost,
							Path:   "/v1/network",
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/network/create_201.json")
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

	expected := slide.Network{
		BridgeDeviceID: "d_0123456789ab",
		ClientID:       "string",
		Comments:       "This is a test network",
		ConnectedVirtIDs: []string{
			"virt_0123456789ab",
		},
		DHCP:           true,
		DHCPRangeEnd:   "10.0.0.200",
		DHCPRangeStart: "10.0.0.100",
		Internet:       true,
		Name:           "Bridge to LAN Network",
		Nameservers:    "1.1.1.1,1.0.0.1",
		NetworkID:      "net_012345",
		RouterPrefix:   "10.0.0.1/24",
		Type:           slide.NetworkTypeDisaster_BRIDGE_LAN,
	}

	ctx := context.Background()

	actual, err := testService.Networks().Create(ctx, slide.NetworkCreatePayload{
		Name: "Bridge to LAN Network",
		Type: slide.NetworkTypeDisaster_BRIDGE_LAN,
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
