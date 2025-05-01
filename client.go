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

type Client struct {
	ClientID string `json:"client_id"`
	Name     string `json:"name"`
	Comments string `json:"comments"`
}

type ClientService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type ClientPayload struct {
	Comments string `json:"comments,omitempty"`
	Name     string `json:"name,omitempty"`
}

// https://docs.slide.tech/api/#tag/clients/GET/v1/client
func (c ClientService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Client]) error,
) error {
	return c.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/clients/GET/v1/client
func (c ClientService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Client]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Client]{}

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

// https://docs.slide.tech/api/#tag/clients/POST/v1/client
func (c ClientService) Create(ctx context.Context, payload ClientPayload) (Client, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Client{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseEndpoint,
		requestBody,
	)

	if err != nil {
		return Client{}, err
	}

	target := Client{}
	if err := c.requestClient.SlideRequest(request, &target); err != nil {
		return Client{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/clients/GET/v1/client/{client_id}
func (c ClientService) Get(ctx context.Context, clientID string) (Client, error) {
	target := Client{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseEndpoint+"/"+clientID,
		http.NoBody,
	)

	if err != nil {
		return Client{}, err
	}

	if err := c.requestClient.SlideRequest(request, &target); err != nil {
		return Client{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/clients/DELETE/v1/client/{client_id}
func (c ClientService) Delete(ctx context.Context, clientID string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		c.baseEndpoint+"/"+clientID,
		http.NoBody,
	)

	if err != nil {
		return err
	}

	return c.requestClient.SlideRequest(request, nil)
}

// https://docs.slide.tech/api/#tag/clients/PATCH/v1/client/{client_id}
func (c ClientService) Update(ctx context.Context, clientID string, payload ClientPayload) (Client, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Client{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		c.baseEndpoint+"/"+clientID,
		requestBody,
	)

	if err != nil {
		return Client{}, err
	}

	target := Client{}
	if err := c.requestClient.SlideRequest(request, &target); err != nil {
		return Client{}, err
	}

	return target, nil
}
