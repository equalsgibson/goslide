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

type Agent struct {
	Addresses           []Address `json:"addresses"`
	AgentID             string    `json:"agent_id"`
	AgentVersion        string    `json:"agent_version"`
	BootedAt            time.Time `json:"booted_at"`
	ClientID            string    `json:"client_id"`
	DeviceID            string    `json:"device_id"`
	DisplayName         string    `json:"display_name"`
	EncryptionAlgorithm string    `json:"encryption_algorithm"`
	FirmwareType        string    `json:"firmware_type"`
	Hostname            string    `json:"hostname"`
	LastSeenAt          time.Time `json:"last_seen_at"`
	Manufacturer        string    `json:"manufacturer"`
	OS                  string    `json:"os"`
	OSVersion           string    `json:"os_version"`
	Platform            string    `json:"platform"`
	PublicIPAddress     string    `json:"public_ip_address"`
}

type AgentService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type AgentAutoPairPayload struct {
	DeviceID    string `json:"device_id"`
	DisplayName string `json:"display_name"`
}

type AgentAutoPairResponse struct {
	AgentID     string `json:"agent_id"`
	DisplayName string `json:"display_name"`
	PairCode    string `json:"pair_code"`
}

type AgentPairPayload struct {
	DeviceID string `json:"device_id"`
	PairCode string `json:"pair_code"`
}

// https://docs.slide.tech/api/#tag/agents/GET/v1/agent
func (a AgentService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Agent]) error,
) error {
	return a.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/agents/GET/v1/agent
func (c AgentService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Agent]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Agent]{}

		endpoint := c.baseEndpoint
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

		if err := c.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/agents/POST/v1/agent
func (c AgentService) AutoPair(ctx context.Context, payload AgentAutoPairPayload) (AgentAutoPairResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return AgentAutoPairResponse{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return AgentAutoPairResponse{}, err
	}

	target := AgentAutoPairResponse{}
	if err := c.requestClient.SlideRequest(request, &target); err != nil {
		return AgentAutoPairResponse{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/agents/POST/v1/agent/pair
func (c AgentService) Pair(ctx context.Context, payload AgentPairPayload) (Agent, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Agent{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return Agent{}, err
	}

	target := Agent{}
	if err := c.requestClient.SlideRequest(request, &target); err != nil {
		return Agent{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/agents/GET/v1/agent/{agent_id}
func (c AgentService) Get(ctx context.Context, agentID string) (Agent, error) {
	target := Agent{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseEndpoint+"/"+agentID,
		http.NoBody,
	)

	if err != nil {
		return Agent{}, err
	}

	if err := c.requestClient.SlideRequest(request, &target); err != nil {
		return Agent{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/agents/PATCH/v1/agent/{agent_id}
func (c AgentService) Update(ctx context.Context, agentID, displayName string) (Agent, error) {
	type agentPayload struct {
		DisplayName string `json:"display_name"`
	}

	payload := agentPayload{
		DisplayName: displayName,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Agent{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		c.baseEndpoint+"/"+agentID,
		requestBody,
	)

	if err != nil {
		return Agent{}, err
	}

	target := Agent{}
	if err := c.requestClient.SlideRequest(request, &target); err != nil {
		return Agent{}, err
	}

	return target, nil
}
