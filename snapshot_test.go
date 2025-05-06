package slide_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/equalsgibson/slide"
	"github.com/equalsgibson/slide/internal/roundtripper"
)

func TestSnapshot_List(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/snapshot/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/snapshot",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/snapshot/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/snapshot",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []slide.Snapshot{}

	ctx := context.Background()
	if err := testService.Snapshots().List(ctx,
		func(response slide.ListResponse[slide.Snapshot]) error {
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

func TestSnapshot_Get(t *testing.T) {
	snapshotID := "s_0123456789ab"

	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/snapshot/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/snapshot/" + snapshotID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.Snapshots().Get(ctx, snapshotID)
	if err != nil {
		t.Fatal(err)
	}

	expected := slide.Snapshot{
		AgentID:         "a_0123456789ab",
		BackupEndedAt:   generateRFC3389FromString(t, "2024-08-23T01:40:08Z"),
		BackupStartedAt: generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		SnapshotID:      "s_0123456789ab",
		Locations: []slide.SnapshotLocation{
			{
				DeviceID: "d_0123456789ab",
				Type:     slide.SnapshotLocationType_LOCAL,
			},
		},
		VerifyBootScreenshotURL: "https://example.com",
		VerifyBootStatus:        slide.SnapshotBootStatus_SUCCESS,
		VerifyFSStatus:          slide.SnapshotFSStatus_SUCCESS,
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
