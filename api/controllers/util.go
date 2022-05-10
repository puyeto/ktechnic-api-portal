package controllers

import (
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/app"
)

const (
	DEFAULT_PAGE_SIZE int = 100
	MAX_PAGE_SIZE     int = 100000
)

func getPaginatedListFromRequest(c *routing.Context, count int) *app.PaginatedList {
	page := parseInt(c.Query("page"), 1)
	perPage := parseInt(c.Query("per_page"), DEFAULT_PAGE_SIZE)
	if perPage == -1 {
		perPage = count
	}
	if perPage == 0 {
		perPage = DEFAULT_PAGE_SIZE
	}
	if perPage > MAX_PAGE_SIZE {
		perPage = MAX_PAGE_SIZE
	}
	return app.NewPaginatedList(page, perPage, count)
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}
