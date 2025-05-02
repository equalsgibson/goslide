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

type VirtualMachineRestore struct {
	AgentID      string `json:"agent_id"`
	CPUCount     string `json:"cpu_count"`
	CreatedAt    string `json:"created_at"`
	DeviceID     string `json:"device_id"`
	DiskBus      string `json:"disk_bus"`
	ExpiresAt    string `json:"expires_at"`
	MemoryInMB   string `json:"memory_in_mb"`
	NetworkModel string `json:"network_model"`
	NetworkType  string `json:"network_type"`
	SnapshotID   string `json:"snapshot_id"`
	State        string `json:"state"`
	VirtID       string `json:"virt_id"`
	VNC          string `json:"vnc"`
	VNCPassword  string `json:"vnc_password"`
}

type VirtualMachineRestoreService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type VirtualMachineVNC struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Type         string `json:"type"`
	WebsocketURI string `json:"websocket_uri"`
}

// https://docs.slide.tech/api/#tag/restores-virtual-machine/GET/v1/restore/virt
func (v VirtualMachineRestoreService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[VirtualMachineRestore]) error,
) error {
	return v.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/restores-virtual-machine/GET/v1/restore/virt
func (v VirtualMachineRestoreService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[VirtualMachineRestore]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[VirtualMachineRestore]{}

		endpoint := v.baseEndpoint
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

		if err := v.requestClient.SlideRequest(request, &target); err != nil {
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
func (v VirtualMachineRestoreService) Create(ctx context.Context, payload VirtualMachineRestoreCreatePayload) (VirtualMachineRestore, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return VirtualMachineRestore{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		v.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return VirtualMachineRestore{}, err
	}

	target := VirtualMachineRestore{}
	if err := v.requestClient.SlideRequest(request, &target); err != nil {
		return VirtualMachineRestore{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/restores-image/GET/v1/restore/image/{image_export_id}
func (v VirtualMachineRestoreService) Get(ctx context.Context, imageExportRestoreID string) (ImageExportRestore, error) {
	target := ImageExportRestore{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		v.baseEndpoint+"/"+imageExportRestoreID,
		http.NoBody,
	)

	if err != nil {
		return ImageExportRestore{}, err
	}

	if err := v.requestClient.SlideRequest(request, &target); err != nil {
		return ImageExportRestore{}, err
	}

	return target, nil
}

type VirtualMachineRestoreCreatePayload struct {
	DeviceID   string `json:"device_id"`
	SnapshotID string `json:"snapshot_id"`

	BootMods     []BootMod    `json:"boot_mods,omitempty"`
	CPUCount     uint         `json:"cpu_count,omitempty"`
	DiskBus      DiskBus      `json:"disk_bus,omitempty"`
	MemoryInMB   uint         `json:"memory_in_mb,omitempty"`
	NetworkModel NetworkModel `json:"network_model,omitempty"`
	NetworkType  NetworkType  `json:"network_type,omitempty"`
}

type BootMod string

const (
	BootMod_PASSWORDLESS_ADMIN_USER BootMod = "passwordless_admin_user"
)

type DiskBus string

const (
	DiskBus_SATA   DiskBus = "sata"
	DiskBus_VIRTIO DiskBus = "virtio"
)

type NetworkModel string

const (
	NetworkModel_HYPERVISOR_DEFAULT NetworkModel = "hypervisor_default"
	NetworkModel_E1000              NetworkModel = "e1000"
	NetworkModel_RTL8139            NetworkModel = "rtl8139"
)

type NetworkType string

const (
	NetworkType_NETWORK          NetworkType = "network"
	NetworkType_NETWORK_ISOLATED NetworkType = "network-isolated"
	NetworkType_BRIDGE           NetworkType = "bridge"
)
