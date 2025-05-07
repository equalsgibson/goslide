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

func TestRestore_Virtual_Machine_List(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/restore_virtual_machine/list_page1_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/restore/virt",
							Query:  url.Values{},
						},
					),
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/restore_virtual_machine/list_page2_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/restore/virt",
							Query: url.Values{
								"offset": []string{"1"},
							},
						},
					),
				},
			),
		),
	)

	actual := []goslide.VirtualMachineRestore{}

	ctx := context.Background()
	if err := testService.VirtualMachineRestores().List(ctx,
		func(response goslide.ListResponse[goslide.VirtualMachineRestore]) error {
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

func TestRestore_Virtual_Machine_Get(t *testing.T) {
	virtualMachineRestoreID := "virt_0123456789ab"
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusOK,
							FilePath:   "testdata/responses/restore_virtual_machine/get_200.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodGet,
							Path:   "/v1/restore/virt/" + virtualMachineRestoreID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	actual, err := testService.VirtualMachineRestores().Get(ctx, virtualMachineRestoreID)
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.VirtualMachineRestore{
		AgentID:      "a_0123456789ab",
		CPUCount:     2,
		CreatedAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		DeviceID:     "d_0123456789ab",
		DiskBus:      "sata",
		ExpiresAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		MemoryInMB:   4096,
		NetworkModel: "e1000",
		NetworkType:  "bridged",
		SnapshotID:   "s_0123456789ab",
		State:        "running",
		VirtID:       "virt_0123456789ab",
		VNC: []goslide.VirtualMachineVNC{
			{
				Host:         "192.168.1.53",
				Port:         12345,
				Type:         goslide.VirtualMachineVNCType_LOCAL,
				WebsocketURI: "wss://example.com",
			},
		},
		VNCPassword: "super-secret",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestRestore_Virtual_Machine_Delete(t *testing.T) {
	virtualMachineRestoreID := "virt_0123456789ab"

	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
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
							Path:   "/v1/restore/virt/" + virtualMachineRestoreID,
							Query:  url.Values{},
						},
					),
				},
			),
		),
	)

	ctx := context.Background()
	if err := testService.VirtualMachineRestores().Delete(ctx, virtualMachineRestoreID); err != nil {
		t.Fatal(err)
	}
}

func TestRestore_Virtual_Machine_Create(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusCreated,
							FilePath:   "testdata/responses/restore_virtual_machine/create_201.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPost,
							Path:   "/v1/restore/virt",
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/restore_virtual_machine/create_201.json")
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
	actual, err := testService.VirtualMachineRestores().Create(ctx, goslide.VirtualMachineRestoreCreatePayload{
		DeviceID:   "d_0123456789ab",
		SnapshotID: "s_0123456789ab",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.VirtualMachineRestore{
		AgentID:      "a_0123456789ab",
		CPUCount:     2,
		CreatedAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		DeviceID:     "d_0123456789ab",
		DiskBus:      "sata",
		ExpiresAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		MemoryInMB:   4096,
		NetworkModel: "e1000",
		NetworkType:  "bridged",
		SnapshotID:   "s_0123456789ab",
		State:        "running",
		VirtID:       "virt_0123456789ab",
		VNC: []goslide.VirtualMachineVNC{
			{
				Host:         "192.168.1.53",
				Port:         12345,
				Type:         goslide.VirtualMachineVNCType_LOCAL,
				WebsocketURI: "wss://example.com",
			},
		},
		VNCPassword: "super-secret",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestRestore_Virtual_Machine_Create_With_Options(t *testing.T) {
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusCreated,
							FilePath:   "testdata/responses/restore_virtual_machine/create_201.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPost,
							Path:   "/v1/restore/virt",
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/restore_virtual_machine/create_with_options_201.json")
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
	actual, err := testService.VirtualMachineRestores().Create(ctx, goslide.VirtualMachineRestoreCreatePayload{
		DeviceID:   "d_0123456789ab",
		SnapshotID: "s_0123456789ab",
		BootMods: []goslide.BootMod{
			goslide.BootMod_PASSWORDLESS_ADMIN_USER,
		},
		CPUCount:     2,
		DiskBus:      goslide.DiskBus_SATA,
		MemoryInMB:   4096,
		NetworkModel: goslide.VirtualMachineNetworkModel_E1000,
		NetworkType:  goslide.VirtualMachineNetworkType_BRIDGE,
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.VirtualMachineRestore{
		AgentID:      "a_0123456789ab",
		CPUCount:     2,
		CreatedAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		DeviceID:     "d_0123456789ab",
		DiskBus:      "sata",
		ExpiresAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		MemoryInMB:   4096,
		NetworkModel: "e1000",
		NetworkType:  "bridged",
		SnapshotID:   "s_0123456789ab",
		State:        "running",
		VirtID:       "virt_0123456789ab",
		VNC: []goslide.VirtualMachineVNC{
			{
				Host:         "192.168.1.53",
				Port:         12345,
				Type:         goslide.VirtualMachineVNCType_LOCAL,
				WebsocketURI: "wss://example.com",
			},
		},
		VNCPassword: "super-secret",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestRestore_Virtual_Machine_Update(t *testing.T) {
	virtID := "virt_0123456789ab"
	testService := goslide.NewService("fakeToken",
		goslide.WithCustomRoundtripper(
			roundtripper.NetworkQueue(
				t,
				[]roundtripper.TestRoundTripFunc{
					roundtripper.ServeAndValidate(
						t,
						&roundtripper.TestResponseFile{
							StatusCode: http.StatusCreated,
							FilePath:   "testdata/responses/restore_virtual_machine/update_202.json",
						},
						roundtripper.ExpectedTestRequest{
							Method: http.MethodPatch,
							Path:   "/v1/restore/virt/" + virtID,
							Query:  url.Values{},
							Validator: func(r *http.Request) error {
								expectedBody, err := os.ReadFile("testdata/requests/restore_virtual_machine/update_202.json")
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
	actual, err := testService.VirtualMachineRestores().Update(ctx, virtID, goslide.VirtualMachineRestoreUpdatePayload{
		State:      goslide.VirtualMachineState_RUNNING,
		CPUCount:   2,
		MemoryInMB: 4096,
		ExpiresAt:  generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := goslide.VirtualMachineRestore{
		AgentID:      "a_0123456789ab",
		CPUCount:     2,
		CreatedAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		DeviceID:     "d_0123456789ab",
		DiskBus:      "sata",
		ExpiresAt:    generateRFC3389FromString(t, "2024-08-23T01:25:08Z"),
		MemoryInMB:   4096,
		NetworkModel: "e1000",
		NetworkType:  "bridged",
		SnapshotID:   "s_0123456789ab",
		State:        "running",
		VirtID:       "virt_0123456789ab",
		VNC: []goslide.VirtualMachineVNC{
			{
				Host:         "192.168.1.53",
				Port:         12345,
				Type:         goslide.VirtualMachineVNCType_LOCAL,
				WebsocketURI: "wss://example.com",
			},
		},
		VNCPassword: "super-secret",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s Returned struct mismatch (-want +got):\n%s", t.Name(), diff)
	}
}
