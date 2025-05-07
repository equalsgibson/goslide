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

func TestUser_List(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/user/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/user",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/user/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/user",
							Query: url.Values{
								"offset": []string{"2"},
							},
						},
					),
				},
			),
		),
	)

	actual := []goslide.User{}

	ctx := context.Background()
	if err := testService.Users().List(ctx,
		func(response goslide.ListResponse[goslide.User]) error {
			actual = append(actual, response.Data...)

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 4 {
		t.Fatal(actual)
	}
}

func TestUser_Get(t *testing.T) {
	userID := "u_1111111111ab"

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/user/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/user/" + userID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.Users().Get(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.User{
		DisplayName: "John Doe",
		Email:       "john.doe@gmail.com",
		FirstName:   "John",
		LastName:    "Doe",
		RoleID:      goslide.UserRole_AccountOwner,
		UserID:      userID,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}
