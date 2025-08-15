package main

import (
	"net/http"

	"github.com/szoumoc/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPosthandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := ReadJson(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := 1
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags: payload.Tags,
		//TODO: CHANGE AFTER AUTH
		UserID: 1,
	}
	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
