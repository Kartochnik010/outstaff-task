package handler

import (
	"net/http"
	"strconv"

	"github.com/kartochnik010/outstaff-task/internal/domain/models"
)

func getSearchMeta(r *http.Request) (*models.SearchMetadata, error) {
	var (
		limit int64
		page  int64
		err   error
		id    uint64
	)
	if r.URL.Query().Get("limit") == "" {
		limit = 10
	} else {
		limit, err = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
		if limit == 0 {
			limit = 10
		}
	}

	if r.URL.Query().Get("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			return nil, err
		}
		if page == 0 {
			page = 1
		}
	}

	strID := r.URL.Query().Get("id")
	if strID != "" {
		id, err = strconv.ParseUint(strID, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &models.SearchMetadata{
		Limit: int(limit),
		Page:  int(page),
		ID:    id,
		Group: r.URL.Query().Get("group"),
		Song:  r.URL.Query().Get("song"),
		Link:  r.URL.Query().Get("link"),
		Text:  r.URL.Query().Get("text"),
	}, nil
}
