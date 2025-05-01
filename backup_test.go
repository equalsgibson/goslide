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

func TestBackup_List(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/backup/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/backup",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/backup/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/backup",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []slide.Backup{}

	ctx := context.Background()
	if err := testService.Backups().List(ctx,
		func(response slide.ListResponse[slide.Backup]) error {
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

func TestBackup_StartBackup(t *testing.T) {
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseNoContent{
							StatusCode: http.StatusAccepted,
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPost,
							Path:   "/v1/backup",
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/backup/start_backup_202.json")
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

	ctx := context.Background()

	if err := testService.Backups().StartBackup(ctx, "a_0123456789ab"); err != nil {
		t.Fatal(err)
	}
}

func TestBackup_Get(t *testing.T) {
	agentID := "al_0123456789ab"
	testService := slide.NewService("fakeToken",
		slide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/backup/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/backup/" + agentID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.Backups().Get(ctx, agentID)
	if err != nil {
		t.Fatal(err)
	}

	expected := slide.Backup{
		AgentID:      "a_0123456789ab",
		BackupID:     "b_0123456789ab",
		EndedAt:      "2024-08-23T01:40:08Z",
		ErrorCode:    1,
		ErrorMessage: "string",
		SnapshotID:   "s_0123456789ab",
		StartedAt:    "2024-08-23T01:25:08Z",
		Status:       slide.BackupStatus_SUCCEEDED,
	}

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}
