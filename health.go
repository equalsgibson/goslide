package slide

import (
	"context"
	"errors"
	"net/http"
)

// healthCheck is the service that is used to report whether an API Token is valid.
//
// TODO: This can be deprecated if Slide implements an endpoint to validate tokens.
type HealthCheck struct {
	requestClient *requestClient
}

// IsAuthenticated makes a simple GET request to the list users endpoint with a limit of 1.
// If an error is encountered, the http.Response error code is checked to validate if the error is
// authentication related (401) or a generic API error
func (h HealthCheck) IsAuthenticated(ctx context.Context) (bool, error) {
	userListReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v1/user?limit=1", http.NoBody)
	if err != nil {
		return false, err
	}

	err = h.requestClient.SlideRequest(userListReq, nil)
	if err != nil {
		var slideError *SlideError
		if errors.As(err, &slideError) {
			for _, errCode := range slideError.Codes {
				if errCode == APIErrorCode_ERR_MISSING_AUTHENTICATION ||
					errCode == APIErrorCode_ERR_UNAUTHORIZED {
					return false, slideError
				}
			}
		} else {
			return true, err
		}
	}

	return true, nil
}
