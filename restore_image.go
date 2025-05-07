package goslide

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

type ImageExportRestore struct {
	AgentID       string          `json:"agent_id"`
	CreatedAt     time.Time       `json:"created_at"`
	DeviceID      string          `json:"device_id"`
	ImageExportID string          `json:"image_export_id"`
	ImageType     ImageExportType `json:"image_type"`
	SnapshotID    string          `json:"snapshot_id"`
}

type ImageExportRestoreService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type ImageExportType string

const (
	ImageExportType_VHDX         ImageExportType = "vhdx"
	ImageExportType_VHDX_DYNAMIC ImageExportType = "vhdx-dynamic"
	ImageExportType_VHD          ImageExportType = "vhd"
	ImageExportType_RAW          ImageExportType = "raw"
)

type ImageExportRestoreData struct {
	DiskID       string                          `json:"disk_id"`
	DownloadURIs []ImageExportRestoreDownloadURI `json:"download_uris"`
	Name         string                          `json:"name"`
	Size         uint                            `json:"size"`
}

type ImageExportRestoreDownloadURI struct {
	Type ImageExportDownloadType `json:"type"`
	URI  string                  `json:"uri"`
}

type ImageExportDownloadType string

const (
	ImageExportDownloadType_LOCAL ImageExportDownloadType = "local"
	ImageExportDownloadType_CLOUD ImageExportDownloadType = "cloud"
)

type ImageExportRestorePayload struct {
	DeviceID   string          `json:"device_id"`
	SnapshotID string          `json:"snapshot_id"`
	ImageType  ImageExportType `json:"image_type"`

	BootMods []BootMod `json:"boot_mods,omitempty"`
}

// https://docs.slide.tech/api/#tag/restores-image/GET/v1/restore/image
func (i ImageExportRestoreService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[ImageExportRestore]) error,
) error {
	return i.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/restores-image/GET/v1/restore/image
func (i ImageExportRestoreService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[ImageExportRestore]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[ImageExportRestore]{}

		endpoint := i.baseEndpoint
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

		if err := i.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/restores-image/POST/v1/restore/image
func (i ImageExportRestoreService) Create(ctx context.Context, payload ImageExportRestorePayload) (ImageExportRestore, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ImageExportRestore{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		i.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return ImageExportRestore{}, err
	}

	target := ImageExportRestore{}
	if err := i.requestClient.SlideRequest(request, &target); err != nil {
		return ImageExportRestore{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/restores-image/GET/v1/restore/image/{image_export_id}
func (i ImageExportRestoreService) Get(ctx context.Context, imageExportRestoreID string) (ImageExportRestore, error) {
	target := ImageExportRestore{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		i.baseEndpoint+"/"+imageExportRestoreID,
		http.NoBody,
	)

	if err != nil {
		return ImageExportRestore{}, err
	}

	if err := i.requestClient.SlideRequest(request, &target); err != nil {
		return ImageExportRestore{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/restores-image/GET/v1/restore/image/{image_export_id}/browse
func (i ImageExportRestoreService) Browse(
	ctx context.Context,
	imageExportRestoreID string,
	pageHandler func(response ListResponse[ImageExportRestoreData]) error,
) error {
	return i.BrowseWithQueryParameters(ctx, imageExportRestoreID, pageHandler)
}

// https://docs.slide.tech/api/#tag/restores-image/GET/v1/restore/image/{image_export_id}/browse
func (i ImageExportRestoreService) BrowseWithQueryParameters(
	ctx context.Context,
	imageExportRestoreID string,
	pageHandler func(response ListResponse[ImageExportRestoreData]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[ImageExportRestoreData]{}

		endpoint := i.baseEndpoint + "/" + imageExportRestoreID + "/browse"
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

		if err := i.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/restores-image/DELETE/v1/restore/image/{image_export_id}
func (i ImageExportRestoreService) Delete(ctx context.Context, imageExportRestoreID string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		i.baseEndpoint+"/"+imageExportRestoreID,
		http.NoBody,
	)

	if err != nil {
		return err
	}

	return i.requestClient.SlideRequest(request, nil)
}
