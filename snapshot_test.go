package goslide_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/equalsgibson/goslide"
	"github.com/equalsgibson/goslide/internal/roundtripper"
	"github.com/google/go-cmp/cmp"
)

func TestSnapshot_List(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
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

	actual := []goslide.Snapshot{}

	ctx := context.Background()
	if err := testService.Snapshots().List(ctx,
		func(response goslide.ListResponse[goslide.Snapshot]) error {
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

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
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

	expected := goslide.Snapshot{
		AgentID:         "a_0123456789ab",
		BackupEndedAt:   generateRFC3389FromString(t, "2024-08-23T01:40:08Z"),
		BackupStartedAt: generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		SnapshotID:      "s_0123456789ab",
		Locations: []goslide.SnapshotLocation{
			{
				DeviceID: "d_0123456789ab",
				Type:     goslide.SnapshotLocationType_LOCAL,
			},
		},
		VerifyBootScreenshotURL: "https://example.com",
		VerifyBootStatus:        goslide.SnapshotBootStatus_SUCCESS,
		VerifyFSStatus:          goslide.SnapshotFSStatus_SUCCESS,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestSnapshot_Get_Deleted(t *testing.T) {
	snapshotID := "s_0123456789ab"

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/snapshot/get_deleted_200.json",
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

	deletedAt := generateRFC3389FromString(t, "2024-08-23T01:25:08Z")

	expected := goslide.Snapshot{
		AgentID:         "a_0123456789ab",
		BackupEndedAt:   generateRFC3389FromString(t, "2024-08-23T01:40:08Z"),
		BackupStartedAt: generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		SnapshotID:      "s_0123456789ab",
		Deleted:         &deletedAt,
		Deletions: []goslide.SnapshotDeletion{
			{
				Deleted:          deletedAt,
				DeletedBy:        "John Smith",
				FirstAndLastName: "John Smith",
				Type:             goslide.SnapshotLocationType_LOCAL,
			},
		},
		Locations: []goslide.SnapshotLocation{
			{
				DeviceID: "d_0123456789ab",
				Type:     goslide.SnapshotLocationType_LOCAL,
			},
		},
		VerifyBootScreenshotURL: "https://example.com",
		VerifyBootStatus:        goslide.SnapshotBootStatus_SUCCESS,
		VerifyFSStatus:          goslide.SnapshotFSStatus_SUCCESS,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}
