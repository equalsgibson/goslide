package slide

import (
	"net/url"
	"strconv"
)

type OffsetPagination struct {
	Total      uint  `json:"total"`
	NextOffset *uint `json:"next_offset"`
}

type ListResponse[Record any] struct {
	Pagination OffsetPagination `json:"pagination"`
	Data       []Record         `json:"data"`
}

type paginationQueryParam func(u url.Values)

func WithOffset(offset uint) paginationQueryParam {
	return func(u url.Values) {
		u.Set("offset", strconv.FormatUint(uint64(offset), 10))
	}
}

func WithLimit(limit uint) paginationQueryParam {
	return func(u url.Values) {
		u.Set("limit", strconv.FormatUint(uint64(limit), 10))
	}
}

func WithSortDirection(ascending bool) paginationQueryParam {
	return func(u url.Values) {
		u.Set("sort_asc", strconv.FormatBool(ascending))
	}
}

func WithSortBy(field string) paginationQueryParam {
	return func(u url.Values) {
		u.Set("sort_by", url.QueryEscape(field))
	}
}

func WithAgentID(agentID string) paginationQueryParam {
	return func(u url.Values) {
		u.Set("agent_id", url.QueryEscape(agentID))
	}
}

func WithDeviceID(deviceID string) paginationQueryParam {
	return func(u url.Values) {
		u.Set("device_id", url.QueryEscape(deviceID))
	}
}

func WithIncludeResolvedAlerts(b bool) paginationQueryParam {
	return func(u url.Values) {
		u.Set("resolved", strconv.FormatBool(b))
	}
}

func WithSnapshotID(snapshotID string) paginationQueryParam {
	return func(u url.Values) {
		u.Set("snapshot_id", url.QueryEscape(snapshotID))
	}
}

func WithPath(path string) paginationQueryParam {
	return func(u url.Values) {
		u.Set("path", url.QueryEscape(path))
	}
}
