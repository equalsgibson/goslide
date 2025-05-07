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

func TestDevice_List(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/device/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/device",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/device/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/device",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []goslide.Device{}

	ctx := context.Background()
	if err := testService.Devices().List(ctx,
		func(response goslide.ListResponse[goslide.Device]) error {
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

func TestDevice_Get(t *testing.T) {
	deviceID := "d_0123456789ab"

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/device/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/device/" + deviceID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.Devices().Get(ctx, deviceID)
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.Device{
		Addresses: []goslide.Address{
			{
				IPs: []string{
					"192.168.1.104",
				},
				MAC: "62:bb:d3:0d:db:7d",
			},
		},
		BootedAt:              generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		ClientID:              "…",
		DeviceID:              deviceID,
		DisplayName:           "My First Device",
		HardwareModelName:     "Slide Z1, 1 TB",
		Hostname:              "my-hostname-1",
		ImageVersion:          "1.0.0",
		LastSeenAt:            generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		NFR:                   false,
		PackageVersion:        "1.2.3",
		PublicIPAddress:       "74.83.124.111",
		SerialNumber:          "SN123456",
		ServiceModelName:      "Slide Z1 Subscription, 1 TB, 1 Year Cloud Retention",
		ServiceModelNameShort: "1 Year Cloud Retention",
		ServiceStatus:         "active",
		StorageTotalBytes:     1099511627776,
		StorageUsedBytes:      274877906944,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestDevice_Update(t *testing.T) {
	deviceID := "d_0123456789ab"

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/device/update_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPatch,
							Path:   "/v1/device/" + deviceID,
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/device/update_200.json")
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

	ctx := context.Background()
	actual, err := testService.Devices().Update(ctx, deviceID, goslide.DevicePayload{
		DisplayName: "My Device",
		Hostname:    "my-device",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.Device{
		Addresses: []goslide.Address{
			{
				IPs: []string{
					"192.168.1.104",
				},
				MAC: "62:bb:d3:0d:db:7d",
			},
		},
		BootedAt:              generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		ClientID:              "…",
		DeviceID:              deviceID,
		DisplayName:           "My Device",
		HardwareModelName:     "Slide Z1, 1 TB",
		Hostname:              "my-device",
		ImageVersion:          "1.0.0",
		LastSeenAt:            generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		NFR:                   false,
		PackageVersion:        "1.2.3",
		PublicIPAddress:       "74.83.124.111",
		SerialNumber:          "SN123456",
		ServiceModelName:      "Slide Z1 Subscription, 1 TB, 1 Year Cloud Retention",
		ServiceModelNameShort: "1 Year Cloud Retention",
		ServiceStatus:         "active",
		StorageTotalBytes:     1099511627776,
		StorageUsedBytes:      274877906944,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}
