package main

import (
	"net/http"

	"github.com/szoumoc/social/internal/store"
)

func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Pagination, filters, sorting
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(324), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, r, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
