package slide

import (
	"fmt"
	"strings"
)

type APIErrorCode string

const (
	APIErrorCode_ERR_ENDPOINT_NOT_FOUND            APIErrorCode = "err_endpoint_not_found"
	APIErrorCode_ERR_ENTITY_NOT_FOUND              APIErrorCode = "err_entity_not_found"
	APIErrorCode_ERR_VALIDATION_ERROR              APIErrorCode = "err_validation_error"
	APIErrorCode_ERR_MISSING_AUTHENTICATION        APIErrorCode = "err_missing_authentication"
	APIErrorCode_ERR_UNAUTHORIZED                  APIErrorCode = "err_unauthorized"
	APIErrorCode_ERR_INTERNAL_SERVER_ERROR         APIErrorCode = "err_internal_server_error"
	APIErrorCode_ERR_RATE_LIMIT_EXCEEDED           APIErrorCode = "err_rate_limit_exceeded"
	APIErrorCode_ERR_AGENT_NOT_CONNECTED_TO_DEVICE APIErrorCode = "err_agent_not_connected_to_device"
	APIErrorCode_ERR_DEVICE_NOT_CONNECTED_TO_CLOUD APIErrorCode = "err_device_not_connected_to_cloud"
	APIErrorCode_ERR_BACKUP_ALREADY_RUNNING        APIErrorCode = "err_backup_already_running"
	APIErrorCode_ERR_CLIENT_NOT_FOUND              APIErrorCode = "err_client_not_found"
)

type SlideError struct {
	HTTPStatusCode  int
	HTTPRequestPath string
	Codes           []APIErrorCode `json:"codes"`
	Details         []string       `json:"details"`
	Message         string         `json:"message"`
}

func (e *SlideError) Error() string {
	var sb strings.Builder

	sb.Write([]byte("slide api request error"))

	if len(e.Codes) > 0 {
		for _, errcode := range e.Codes {
			sb.WriteString(fmt.Sprintf(" %s ", errcode))
		}
	}

	return sb.String()
}
