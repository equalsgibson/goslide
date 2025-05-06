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

type Device struct {
	Addresses             []Address `json:"addresses"`
	BootedAt              time.Time `json:"booted_at"`
	ClientID              string    `json:"client_id"`
	DeviceID              string    `json:"device_id"`
	DisplayName           string    `json:"display_name"`
	HardwareModelName     string    `json:"hardware_model_name"`
	Hostname              string    `json:"hostname"`
	ImageVersion          string    `json:"image_version"`
	PublicIPAddress       string    `json:"public_ip_address"`
	LastSeenAt            time.Time `json:"last_seen_at"`
	NFR                   bool      `json:"nfr"`
	PackageVersion        string    `json:"package_version"`
	SerialNumber          string    `json:"serial_number"`
	ServiceModelName      string    `json:"service_model_name"`
	ServiceModelNameShort string    `json:"service_model_name_short"`
	ServiceStatus         string    `json:"service_status"`
	StorageTotalBytes     uint64    `json:"storage_total_bytes"`
	StorageUsedBytes      uint64    `json:"storage_used_bytes"`
}

type Address struct {
	MAC string   `json:"mac"`
	IPs []string `json:"ips"`
}

type DeviceService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type DevicePayload struct {
	DisplayName string `json:"display_name,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
}

// https://docs.slide.tech/api/#tag/devices
func (d DeviceService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Device]) error,
) error {
	return d.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/devices
func (d DeviceService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Device]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Device]{}

		endpoint := d.baseEndpoint
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

		if err := d.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/devices/GET/v1/device/{device_id}
func (d DeviceService) Get(ctx context.Context, deviceID string) (Device, error) {
	target := Device{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		d.baseEndpoint+"/"+deviceID,
		http.NoBody,
	)

	if err != nil {
		return Device{}, err
	}

	if err := d.requestClient.SlideRequest(request, &target); err != nil {
		return Device{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/devices/PATCH/v1/device/{device_id}
func (d DeviceService) Update(ctx context.Context, deviceID string, payload DevicePayload) (Device, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Device{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		d.baseEndpoint+"/"+deviceID,
		requestBody,
	)

	if err != nil {
		return Device{}, err
	}

	target := Device{}
	if err := d.requestClient.SlideRequest(request, &target); err != nil {
		return Device{}, err
	}

	return target, nil
}
