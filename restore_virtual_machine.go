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

type VirtualMachineRestore struct {
	AgentID      string                     `json:"agent_id"`
	CPUCount     uint                       `json:"cpu_count"`
	CreatedAt    time.Time                  `json:"created_at"`
	DeviceID     string                     `json:"device_id"`
	DiskBus      DiskBus                    `json:"disk_bus"`
	ExpiresAt    time.Time                  `json:"expires_at"`
	MemoryInMB   uint                       `json:"memory_in_mb"`
	NetworkModel VirtualMachineNetworkModel `json:"network_model"`
	NetworkType  VirtualMachineNetworkType  `json:"network_type"`
	SnapshotID   string                     `json:"snapshot_id"`
	State        VirtualMachineState        `json:"state"`
	VirtID       string                     `json:"virt_id"`
	VNC          []VirtualMachineVNC        `json:"vnc"`
	VNCPassword  string                     `json:"vnc_password"`
}

type VirtualMachineRestoreService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type VirtualMachineVNC struct {
	Host         string                `json:"host"`
	Port         uint                  `json:"port"`
	Type         VirtualMachineVNCType `json:"type"`
	WebsocketURI string                `json:"websocket_uri"`
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

// https://docs.slide.tech/api/#tag/restores-virtual-machine/POST/v1/restore/virt
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

// https://docs.slide.tech/api/#tag/restores-virtual-machine/GET/v1/restore/virt/{virt_id}
func (v VirtualMachineRestoreService) Get(ctx context.Context, virtID string) (VirtualMachineRestore, error) {
	target := VirtualMachineRestore{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		v.baseEndpoint+"/"+virtID,
		http.NoBody,
	)

	if err != nil {
		return VirtualMachineRestore{}, err
	}

	if err := v.requestClient.SlideRequest(request, &target); err != nil {
		return VirtualMachineRestore{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/restores-virtual-machine/DELETE/v1/restore/virt/{virt_id}
func (v VirtualMachineRestoreService) Delete(ctx context.Context, virtID string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		v.baseEndpoint+"/"+virtID,
		http.NoBody,
	)

	if err != nil {
		return err
	}

	return v.requestClient.SlideRequest(request, nil)
}

// https://docs.slide.tech/api/#tag/restores-virtual-machine/PATCH/v1/restore/virt/{virt_id}
func (v VirtualMachineRestoreService) Update(ctx context.Context, virtID string, state VirtualMachineState) (VirtualMachineRestore, error) {
	type virtualMachineRestoreUpdatePayload struct {
		State VirtualMachineState `json:"state"`
	}

	payload := virtualMachineRestoreUpdatePayload{
		State: state,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return VirtualMachineRestore{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		v.baseEndpoint+"/"+virtID,
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

type VirtualMachineRestoreCreatePayload struct {
	DeviceID   string `json:"device_id"`
	SnapshotID string `json:"snapshot_id"`

	BootMods     []BootMod                  `json:"boot_mods,omitempty"`
	CPUCount     uint                       `json:"cpu_count,omitempty"`
	DiskBus      DiskBus                    `json:"disk_bus,omitempty"`
	MemoryInMB   uint                       `json:"memory_in_mb,omitempty"`
	NetworkModel VirtualMachineNetworkModel `json:"network_model,omitempty"`
	NetworkType  VirtualMachineNetworkType  `json:"network_type,omitempty"`
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

type VirtualMachineNetworkModel string

const (
	VirtualMachineNetworkModel_HYPERVISOR_DEFAULT VirtualMachineNetworkModel = "hypervisor_default"
	VirtualMachineNetworkModel_E1000              VirtualMachineNetworkModel = "e1000"
	VirtualMachineNetworkModel_RTL8139            VirtualMachineNetworkModel = "rtl8139"
)

type VirtualMachineNetworkType string

const (
	VirtualMachineNetworkType_NETWORK          VirtualMachineNetworkType = "network"
	VirtualMachineNetworkType_NETWORK_ISOLATED VirtualMachineNetworkType = "network-isolated"
	VirtualMachineNetworkType_BRIDGE           VirtualMachineNetworkType = "bridge"
)

type VirtualMachineState string

const (
	VirtualMachineState_RUNNING VirtualMachineState = "running"
	VirtualMachineState_STOPPED VirtualMachineState = "stopped"
	VirtualMachineState_PAUSED  VirtualMachineState = "paused"
)

type VirtualMachineVNCType string

const (
	VirtualMachineVNCType_LOCAL VirtualMachineVNCType = "local"
	VirtualMachineVNCType_CLOUD VirtualMachineVNCType = "cloud"
)
