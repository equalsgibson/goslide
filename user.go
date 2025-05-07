package goslide

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type User struct {
	DisplayName string   `json:"display_name"`
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	RoleID      UserRole `json:"role_id"`
	UserID      string   `json:"user_id"`
}

type UserRole string

const (
	UserRole_AccountOwner    UserRole = "r_account_owner"
	UserRole_AccountAdmin    UserRole = "r_account_admin"
	UserRole_AccountTech     UserRole = "r_account_tech"
	UserRole_AccountReadOnly UserRole = "r_readonly"
)

type UserService struct {
	baseEndpoint  string
	requestClient *requestClient
}

// https://docs.slide.tech/api/#tag/users/GET/v1/user
func (u UserService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[User]) error,
) error {
	return u.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/users/GET/v1/user
func (u UserService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[User]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[User]{}

		endpoint := u.baseEndpoint
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

		if err := u.requestClient.SlideRequest(request, &target); err != nil {
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

// https://docs.slide.tech/api/#tag/users/GET/v1/user/{user_id}
func (u UserService) Get(ctx context.Context, userID string) (User, error) {
	target := User{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/v1/user/"+userID,
		http.NoBody,
	)

	if err != nil {
		return User{}, err
	}

	if err := u.requestClient.SlideRequest(request, &target); err != nil {
		return User{}, err
	}

	return target, nil
}
