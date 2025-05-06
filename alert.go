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

type Alert struct {
	AgentID     string     `json:"agent_id"`
	AlertFields string     `json:"alert_fields"`
	AlertID     string     `json:"alert_id"`
	AlertType   AlertType  `json:"alert_type"`
	CreatedAt   time.Time  `json:"created_at"`
	DeviceID    string     `json:"device_id"`
	Resolved    bool       `json:"resolved"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	ResolvedBy  string     `json:"resolved_by"`
}

type AlertService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type AlertType string

const (
	AlertType_DEVICE_NOT_CHECKING_IN        AlertType = "device_not_checking_in"
	AlertType_DEVICE_OUT_OF_DATE            AlertType = "device_out_of_date"
	AlertType_DEVICE_STORAGE_NOT_HEALTHY    AlertType = "device_storage_not_healthy"
	AlertType_DEVICE_STORAGE_SPACE_LOW      AlertType = "device_storage_space_low"
	AlertType_DEVICE_STORAGE_SPACE_CRITICAL AlertType = "device_storage_space_critical"
	AlertType_AGENT_NOT_CHECKING_IN         AlertType = "agent_not_checking_in"
	AlertType_AGENT_NOT_BACKING_UP          AlertType = "agent_not_backing_up"
	AlertType_AGENT_BACKUP_FAILED           AlertType = "agent_backup_failed"
)

// https://docs.slide.tech/api/#tag/alerts/GET/v1/alert/{alert_id}
func (a AlertService) Get(ctx context.Context, alertID string) (Alert, error) {
	target := Alert{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		a.baseEndpoint+"/"+alertID,
		http.NoBody,
	)

	if err != nil {
		return Alert{}, err
	}

	if err := a.requestClient.SlideRequest(request, &target); err != nil {
		return Alert{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/alerts/GET/v1/alert
func (a AlertService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Alert]) error,
) error {
	return a.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/alerts/GET/v1/alert
func (a AlertService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Alert]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Alert]{}

		endpoint := a.baseEndpoint
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

		if err := a.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/alerts/PATCH/v1/alert/{alert_id}
func (a AlertService) Update(
	ctx context.Context,
	alertID string,
	resolved bool,
) (Alert, error) {
	type alertPayload struct {
		Resolved bool `json:"resolved"`
	}

	payloadBytes, err := json.Marshal(alertPayload{
		Resolved: resolved,
	})
	if err != nil {
		return Alert{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		a.baseEndpoint+"/"+alertID,
		requestBody,
	)

	if err != nil {
		return Alert{}, err
	}

	target := Alert{}
	if err := a.requestClient.SlideRequest(request, &target); err != nil {
		return Alert{}, err
	}

	return target, nil
}
