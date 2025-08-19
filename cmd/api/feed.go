package main

import "net/http"

func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Pagination, filters
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(324))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, r, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
