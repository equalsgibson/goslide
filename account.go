package goslide

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Account struct {
	AccountID      string         `json:"account_id"`
	AccountName    string         `json:"account_name"`
	AlertEmails    []string       `json:"alert_emails"`
	BillingAddress BillingAddress `json:"billing_address"`
	PrimaryContact string         `json:"primary_contact"`
	PrimaryEmail   string         `json:"primary_email"`
	PrimaryPhone   string         `json:"primary_phone"`
}

type BillingAddress struct {
	City       string `json:"City"`
	Country    string `json:"Country"`
	Line1      string `json:"Line1"`
	Line2      string `json:"Line2"`
	PostalCode string `json:"PostalCode"`
	State      string `json:"State"`
}

type AccountService struct {
	baseEndpoint  string
	requestClient *requestClient
}

// https://docs.slide.tech/api/#tag/accounts/GET/v1/account
func (a AccountService) List(
	ctx context.Context,
	pageHandler func(response ListResponse[Account]) error,
) error {
	return a.ListWithQueryParameters(ctx, pageHandler)
}

// https://docs.slide.tech/api/#tag/accounts/GET/v1/account
func (a AccountService) ListWithQueryParameters(
	ctx context.Context,
	pageHandler func(response ListResponse[Account]) error,
	options ...paginationQueryParam,
) error {
	queryParams := url.Values{}
	for _, option := range options {
		option(queryParams)
	}

	for {
		target := ListResponse[Account]{}

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

// https://docs.slide.tech/api/#tag/accounts/GET/v1/account/{account_id}
func (a AccountService) Get(ctx context.Context, accountID string) (Account, error) {
	target := Account{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		a.baseEndpoint+"/"+accountID,
		http.NoBody,
	)

	if err != nil {
		return Account{}, err
	}

	if err := a.requestClient.SlideRequest(request, &target); err != nil {
		return Account{}, err
	}

	return target, nil
}

// https://docs.slide.tech/api/#tag/accounts/PATCH/v1/account/{account_id}
func (a AccountService) Update(
	ctx context.Context,
	accountID string,
	alertEmails []string,
) (Account, error) {
	type accountPayload struct {
		AlertEmails []string `json:"alert_emails"`
	}

	payloadBytes, err := json.Marshal(accountPayload{
		AlertEmails: alertEmails,
	})
	if err != nil {
		return Account{}, err
	}

	requestBody := bytes.NewReader(payloadBytes)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		a.baseEndpoint+"/"+accountID,
		requestBody,
	)

	if err != nil {
		return Account{}, err
	}

	target := Account{}
	if err := a.requestClient.SlideRequest(request, &target); err != nil {
		return Account{}, err
	}

	return target, nil
}
