package slide

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type BackupStatus string

const (
	BackupStatus_CREATED      BackupStatus = "created"
	BackupStatus_PENDING      BackupStatus = "pending"
	BackupStatus_STARTED      BackupStatus = "started"
	BackupStatus_PREFLIGHT    BackupStatus = "preflight"
	BackupStatus_CONTACTING   BackupStatus = "contacting"
	BackupStatus_CREATING_VSS BackupStatus = "creating_vss"
	BackupStatus_PREPARING    BackupStatus = "preparing"
	BackupStatus_TRANSFERRING BackupStatus = "transferring"
	BackupStatus_SNAPSHOT     BackupStatus = "snapshot"
	BackupStatus_FINALIZING   BackupStatus = "finalizing"
	BackupStatus_CANCELING    BackupStatus = "canceling"
	BackupStatus_FAILING      BackupStatus = "failing"
	BackupStatus_CANCELED     BackupStatus = "canceled"
	BackupStatus_FAILED       BackupStatus = "failed"
	BackupStatus_SUCCEEDED    BackupStatus = "succeeded"
)

type Backup struct {
	AgentID      string       `json:"agent_id"`
	BackupID     string       `json:"backup_id"`
	EndedAt      string       `json:"ended_at"`
	ErrorCode    uint         `json:"error_code"`
	ErrorMessage string       `json:"error_message"`
	SnapshotID   string       `json:"snapshot_id"`
	StartedAt    string       `json:"started_at"`
	Status       BackupStatus `json:"status"`
}

type BackupService struct {
	baseEndpoint  string
	requestClient *requestClient
}

// https://docs.slide.tech/api/#tag/backups/GET/v1/backup
func (b BackupService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Backup]) error,
) error {
	return b.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/backups/GET/v1/backup
func (b BackupService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Backup]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Backup]{}

		endpoint := b.baseEndpoint
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

		if err := b.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/backups/GET/v1/backup/{backup_id}
func (b BackupService) Get(ctx context.Context, backupID string) (Backup, error) {
	target := Backup{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		b.baseEndpoint+"/"+backupID,
		http.NoBody,
	)

	if err != nil {
		return Backup{}, err
	}

	if err := b.requestClient.SlideRequest(request, &target); err != nil {
		return Backup{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/backups/POST/v1/backup
func (b BackupService) StartBackup(ctx context.Context, agentID string) error {
	type backupPayload struct {
		AgentID string `json:"agent_id"`
	}

	payloadBytes, err := json.Marshal(backupPayload{
		AgentID: agentID,
	})
	if err != nil {
		return err
	}

	requestBody := bytes.NewReader(payloadBytes)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		b.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return err
	}

	return b.requestClient.SlideRequest(request, nil)
}
