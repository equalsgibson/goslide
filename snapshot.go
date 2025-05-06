package slide

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Snapshot struct {
	AgentID                 string             `json:"agent_id"`
	BackupEndedAt           time.Time          `json:"backup_ended_at"`
	BackupStartedAt         time.Time          `json:"backup_started_at"`
	Locations               []SnapshotLocation `json:"locations"`
	SnapshotID              string             `json:"snapshot_id"`
	VerifyBootScreenshotURL string             `json:"verify_boot_screenshot_url"`
	VerifyBootStatus        SnapshotBootStatus `json:"verify_boot_status"`
	VerifyFSStatus          SnapshotFSStatus   `json:"verify_fs_status"`
}

type SnapshotService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type SnapshotLocation struct {
	DeviceID string               `json:"device_id"`
	Type     SnapshotLocationType `json:"type"`
}

type SnapshotLocationType string

const (
	SnapshotLocationType_LOCAL SnapshotLocationType = "local"
	SnapshotLocationType_CLOUD SnapshotLocationType = "cloud"
)

type SnapshotFSStatus string

const (
	SnapshotFSStatus_SUCCESS SnapshotFSStatus = "success"
	SnapshotFSStatus_WARNING SnapshotFSStatus = "warning"
	SnapshotFSStatus_ERROR   SnapshotFSStatus = "error"
	SnapshotFSStatus_SKIPPED SnapshotFSStatus = "skipped"
)

type SnapshotBootStatus string

const (
	SnapshotBootStatus_SUCCESS                    SnapshotBootStatus = "success"
	SnapshotBootStatus_WARNING                    SnapshotBootStatus = "warning"
	SnapshotBootStatus_ERROR                      SnapshotBootStatus = "error"
	SnapshotBootStatus_SKIPPED                    SnapshotBootStatus = "skipped"
	SnapshotBootStatus_PENDING                    SnapshotBootStatus = "pending"
	SnapshotBootStatus_PENDING_DUE_TO_DISASTER_VM SnapshotBootStatus = "pending_due_to_disaster_vm"
)

// https://docs.slide.tech/api/#tag/snapshots/GET/v1/snapshot
func (s SnapshotService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Snapshot]) error,
) error {
	return s.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/snapshots/GET/v1/snapshot
func (s SnapshotService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Snapshot]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Snapshot]{}

		endpoint := s.baseEndpoint
		if len(queryParams) > 0 {
			endpoint = endpoint + "?"
		}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s%s", endpoint, queryParams.Encode()),
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.requestClient.SlideRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		// No next offset marks the end of the paginated results
		if target.Pagination.NextOffset == nil {
			break
		}

		queryParams.Set(
			"offset",
			strconv.FormatUint(
				uint64(*target.Pagination.NextOffset),
				10,
			),
		)
	}

	return nil
}

// https://docs.slide.tech/api/#tag/snapshots/GET/v1/snapshot/{snapshot_id}
func (s SnapshotService) Get(ctx context.Context, snapshotID string) (Snapshot, error) {
	target := Snapshot{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		s.baseEndpoint+"/"+snapshotID,
		http.NoBody,
	)

	if err != nil {
		return Snapshot{}, err
	}

	if err := s.requestClient.SlideRequest(request, &target); err != nil {
		return Snapshot{}, err
	}

	return target, nil
}
