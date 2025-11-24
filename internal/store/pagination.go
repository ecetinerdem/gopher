package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (pfq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	queryString := r.URL.Query()

	limit := queryString.Get("limit")
	if limit != "" {
		lmt, err := strconv.Atoi(limit)
		if err != nil {
			return pfq, nil
		}
		pfq.Limit = lmt
	}

	offset := queryString.Get("offset")
	if offset != "" {
		ofst, err := strconv.Atoi(offset)
		if err != nil {
			return pfq, nil
		}
		pfq.Offset = ofst
	}

	sort := queryString.Get("sort")
	if sort != "" {
		pfq.Sort = sort
	}

	tags := queryString.Get("tags")
	if tags != "" {
		pfq.Tags = strings.Split(tags, ",")
	}

	search := queryString.Get("search")
	if search != "" {
		pfq.Search = search
	}

	since := queryString.Get("since")
	if since != "" {
		pfq.Since = parseTime(since)
	}

	until := queryString.Get("until")
	if until != "" {
		pfq.Until = parseTime(until)
	}

	return pfq, nil
}

func parseTime(since string) string {
	t, err := time.Parse(time.DateTime, since)

	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
