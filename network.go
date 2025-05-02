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

type NetworkService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type Network struct {
	BridgeDeviceID   string              `json:"bridge_device_id"`
	ClientID         string              `json:"client_id"`
	Comments         string              `json:"comments"`
	ConnectedVirtIDs []string            `json:"connected_virt_ids"`
	DHCP             bool                `json:"dhcp"`
	DHCPRangeEnd     string              `json:"dhcp_range_end"`
	DHCPRangeStart   string              `json:"dhcp_range_start"`
	Internet         bool                `json:"internet"`
	Name             string              `json:"name"`
	Nameservers      string              `json:"nameservers"`
	NetworkID        string              `json:"network_id"`
	RouterPrefix     string              `json:"router_prefix"`
	Type             NetworkTypeDisaster `json:"type"`
}

type NetworkTypeDisaster string

const (
	NetworkTypeDisaster_STANDARD   NetworkTypeDisaster = "standard"
	NetworkTypeDisaster_BRIDGE_LAN NetworkTypeDisaster = "bridge-lan"
)

// https://docs.slide.tech/api/#tag/networks/GET/v1/network
func (n NetworkService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Network]) error,
) error {
	return n.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/networks/GET/v1/network
func (n NetworkService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Network]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Network]{}

		endpoint := n.baseEndpoint
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

		if err := n.requestClient.SlideRequest(request, &target); err != nil {
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

type NetworkCreatePayload struct {
	Name string              `json:"name"`
	Type NetworkTypeDisaster `json:"type"`

	BridgeDeviceID string `json:"bridge_device_id,omitempty"`
	ClientID       string `json:"client_id,omitempty"`
	Comments       string `json:"comments,omitempty"`
	DHCP           bool   `json:"dhcp,omitempty"`
	DHCPRangeEnd   string `json:"dhcp_range_end,omitempty"`
	DHCPRangeStart string `json:"dhcp_range_start,omitempty"`
	Internet       bool   `json:"internet,omitempty"`
	Nameservers    string `json:"nameservers,omitempty"`
	RouterPrefix   string `json:"router_prefix,omitempty"`
}

type NetworkUpdatePayload struct {
	Name           string              `json:"name,omitempty"`
	Type           NetworkTypeDisaster `json:"type,omitempty"`
	Comments       string              `json:"comments,omitempty"`
	DHCP           bool                `json:"dhcp,omitempty"`
	DHCPRangeEnd   string              `json:"dhcp_range_end,omitempty"`
	DHCPRangeStart string              `json:"dhcp_range_start,omitempty"`
	Internet       bool                `json:"internet,omitempty"`
	Nameservers    string              `json:"nameservers,omitempty"`
	RouterPrefix   string              `json:"router_prefix,omitempty"`
}

// https://docs.slide.tech/api/#tag/networks/POST/v1/network
func (n NetworkService) Create(ctx context.Context, payload NetworkCreatePayload) (Network, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Network{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return Network{}, err
	}

	target := Network{}
	if err := n.requestClient.SlideRequest(request, &target); err != nil {
		return Network{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/networks/GET/v1/network/{network_id}
func (n NetworkService) Get(ctx context.Context, networkID string) (Network, error) {
	target := Network{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		n.baseEndpoint+"/"+networkID,
		http.NoBody,
	)

	if err != nil {
		return Network{}, err
	}

	if err := n.requestClient.SlideRequest(request, &target); err != nil {
		return Network{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/networks/DELETE/v1/network/{network_id}
func (n NetworkService) Delete(ctx context.Context, networkID string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		n.baseEndpoint+"/"+networkID,
		http.NoBody,
	)

	if err != nil {
		return err
	}

	return n.requestClient.SlideRequest(request, nil)
}

// https://docs.slide.tech/api/#tag/networks/PATCH/v1/network/{network_id}
func (n NetworkService) Update(ctx context.Context, networkID string, payload NetworkUpdatePayload) (Network, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Network{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		n.baseEndpoint+"/"+networkID,
		requestBody,
	)

	if err != nil {
		return Network{}, err
	}

	target := Network{}
	if err := n.requestClient.SlideRequest(request, &target); err != nil {
		return Network{}, err
	}

	return target, nil
}

type NetworkPortForwardPayload struct {
	Dest      string       `json:"dest"`
	NetworkID string       `json:"network_id"`
	Proto     NetworkProto `json:"proto"`
}

type NetworkProto string

const (
	NetworkProto_UDP NetworkProto = "udp"
	NetworkProto_TCP NetworkProto = "tcp"
)

// https://docs.slide.tech/api/#tag/networks/POST/v1/network/{network_id}/port-forwards
func (n NetworkService) CreatePortForward(ctx context.Context, networkID string, payload NetworkPortForwardPayload) (NetworkPortForward, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return NetworkPortForward{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.baseEndpoint+"/"+networkID+"/port-forwards",
		requestBody,
	)

	if err != nil {
		return NetworkPortForward{}, err
	}

	target := NetworkPortForward{}
	if err := n.requestClient.SlideRequest(request, &target); err != nil {
		return NetworkPortForward{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/networks/DELETE/v1/network/{network_id}/port-forwards
func (n NetworkService) DeletePortForward(ctx context.Context, networkID string, payload NetworkPortForwardPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.baseEndpoint+"/"+networkID+"/port-forwards",
		requestBody,
	)

	if err != nil {
		return err
	}

	if err := n.requestClient.SlideRequest(request, nil); err != nil {
		return err
	}

	return nil
}

// https://docs.slide.tech/api/#tag/networks/PATCH/v1/network/{network_id}/port-forwards
func (n NetworkService) UpdatePortForward(ctx context.Context, networkID string, payload NetworkPortForwardUpdatePayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.baseEndpoint+"/"+networkID+"/port-forwards",
		requestBody,
	)

	if err != nil {
		return err
	}

	if err := n.requestClient.SlideRequest(request, nil); err != nil {
		return err
	}

	return nil
}

type NetworkPortForwardUpdatePayload struct {
	Dest      string       `json:"dest"`
	NetworkID string       `json:"network_id"`
	Port      uint         `json:"port"`
	ProtoNew  NetworkProto `json:"proto_new"`
	ProtoOld  NetworkProto `json:"proto_old"`
}

type NetworkPortForward struct {
	Dest      string       `json:"dest"`
	NetworkID string       `json:"network_id"`
	Port      uint         `json:"port"`
	Proto     NetworkProto `json:"proto"`
}

type NetworkWGPeer struct {
	NetworkID      string   `json:"network_id"`
	PeerName       string   `json:"peer_name"`
	RemoteNetworks []string `json:"remote_networks"`
	WGAddress      string   `json:"wg_address"`
	WGPrivateKey   string   `json:"wg_private_key"`
	WGPublicKey    string   `json:"wg_public_key"`
}

type NetworkWGPeerCreatePayload struct {
	NetworkID      string   `json:"network_id"`
	PeerName       string   `json:"peer_name"`
	RemoteNetworks []string `json:"remote_networks"`
}

// https://docs.slide.tech/api/#tag/networks/POST/v1/network/{network_id}/wg-peers
func (n NetworkService) CreateWGPeer(ctx context.Context, networkID string, payload NetworkWGPeerCreatePayload) (NetworkWGPeer, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return NetworkWGPeer{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.baseEndpoint+"/"+networkID+"/wg-peers",
		requestBody,
	)

	if err != nil {
		return NetworkWGPeer{}, err
	}

	target := NetworkWGPeer{}
	if err := n.requestClient.SlideRequest(request, &target); err != nil {
		return NetworkWGPeer{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/networks/DELETE/v1/network/{network_id}/wg-peers
func (n NetworkService) DeleteWGPeer(ctx context.Context, networkID, wgAddress string) error {
	type networkWGPeerDeletePayload struct {
		NetworkID string `json:"network_id"`
		WGAddress string `json:"wg_address"`
	}

	payload := networkWGPeerDeletePayload{
		NetworkID: networkID,
		WGAddress: wgAddress,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.baseEndpoint+"/"+networkID+"/wg-peers",
		requestBody,
	)

	if err != nil {
		return err
	}

	if err := n.requestClient.SlideRequest(request, nil); err != nil {
		return err
	}

	return nil
}

// https://docs.slide.tech/api/#tag/networks/PATCH/v1/network/{network_id}/wg-peers
func (n NetworkService) UpdateWGPeer(ctx context.Context, networkID string, payload NetworkPortForwardUpdatePayload) (NetworkWGPeer, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return NetworkWGPeer{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.baseEndpoint+"/"+networkID+"/wg-peers",
		requestBody,
	)

	if err != nil {
		return NetworkWGPeer{}, err
	}

	target := NetworkWGPeer{}
	if err := n.requestClient.SlideRequest(request, &target); err != nil {
		return NetworkWGPeer{}, err
	}

	return target, nil
}
