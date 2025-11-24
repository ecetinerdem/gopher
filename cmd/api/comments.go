package main

import (
	"net/http"

	"github.com/ecetinerdem/gopherSocial/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=200"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateCommentPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := getPostFromCtx(r)

	cms := &store.Comment{
		Content: payload.Content,
		UserID:  post.UserID,
		PostID:  post.ID,
	}

	err := app.store.Comments.Create(r.Context(), cms)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
