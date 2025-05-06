package slide

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type FileRestore struct {
	AgentID       string    `json:"agent_id"`
	CreatedAt     time.Time `json:"created_at"`
	DeviceID      string    `json:"device_id"`
	ExpiresAt     time.Time `json:"expires_at"`
	FileRestoreID string    `json:"file_restore_id"`
	SnapshotID    string    `json:"snapshot_id"`
}

type FileRestoreService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type FileRestorePayload struct {
	DeviceID   string `json:"device_id"`
	SnapshotID string `json:"snapshot_id"`
}

// https://docs.slide.tech/api/#tag/restores-file/GET/v1/restore/file
func (f FileRestoreService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[FileRestore]) error,
) error {
	return f.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/restores-file/GET/v1/restore/file
func (f FileRestoreService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[FileRestore]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[FileRestore]{}

		endpoint := f.baseEndpoint
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

		if err := f.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/restores-file/POST/v1/restore/file
func (f FileRestoreService) Create(ctx context.Context, payload FileRestorePayload) (FileRestore, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return FileRestore{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		f.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return FileRestore{}, err
	}

	target := FileRestore{}
	if err := f.requestClient.SlideRequest(request, &target); err != nil {
		return FileRestore{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/restores-file/GET/v1/restore/file/{file_restore_id}
func (f FileRestoreService) Get(ctx context.Context, fileRestoreID string) (FileRestore, error) {
	target := FileRestore{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		f.baseEndpoint+"/"+fileRestoreID,
		http.NoBody,
	)

	if err != nil {
		return FileRestore{}, err
	}

	if err := f.requestClient.SlideRequest(request, &target); err != nil {
		return FileRestore{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/restores-file/GET/v1/restore/file/{file_restore_id}/browse
func (f FileRestoreService) Browse(
	ctx context.Context,
	fileRestoreID string,
	pageHandler func(response ListResponse[FileRestoreData]) error,
) error {
	return f.BrowseWithQueryParameters(ctx, fileRestoreID, pageHandler)
}

// https://docs.slide.tech/api/#tag/restores-file/GET/v1/restore/file/{file_restore_id}/browse
func (f FileRestoreService) BrowseWithQueryParameters(
	ctx context.Context,
	fileRestoreID string,
	pageHandler func(response ListResponse[FileRestoreData]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[FileRestoreData]{}

		endpoint := f.baseEndpoint + "/" + fileRestoreID + "/browse"
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

		if err := f.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/restores-file/DELETE/v1/restore/file/{file_restore_id}
func (f FileRestoreService) Delete(ctx context.Context, fileRestoreID string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		f.baseEndpoint+"/"+fileRestoreID,
		http.NoBody,
	)

	if err != nil {
		return err
	}

	return f.requestClient.SlideRequest(request, nil)
}

type FileRestoreData struct {
	DownloadURIs      []FileRestoreDownloadURI `json:"download_uris"`
	ModifiedAt        string                   `json:"modified_at"`
	Name              string                   `json:"name"`
	Path              string                   `json:"path"`
	Size              uint                     `json:"size"`
	SymlinkTargetPath string                   `json:"symlink_target_path"`
	Type              FileRestoreDataType      `json:"type"`
}

type FileRestoreDownloadURI struct {
	Type FileRestoreDownloadType `json:"type"`
	URI  string                  `json:"uri"`
}

type FileRestoreDownloadType string

const (
	FileRestoreDownloadType_LOCAL FileRestoreDownloadType = "local"
	FileRestoreDownloadType_CLOUD FileRestoreDownloadType = "cloud"
)

type FileRestoreDataType string

const (
	FileRestoreDataType_FILE    FileRestoreDataType = "file"
	FileRestoreDataType_DIR     FileRestoreDataType = "dir"
	FileRestoreDataType_SYMLINK FileRestoreDataType = "symlink"
)
