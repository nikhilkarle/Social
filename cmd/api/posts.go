package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nikhilkarle/social/internal/store"
)

type postKey string
const postCtx postKey = "post"

type CreatePostPayload struct{
	Title string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=1000"`
	Tags []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request){
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil{
		app.badRequestError(w,r,err)
		return
	}

	if err := Validate.Struct(payload); err != nil{
		app.badRequestError(w,r,err)
		return
	}

	post := &store.Post{
		Title: payload.Title,
		Content: payload.Content,
		Tags: payload.Tags,
		//TODO: change afer auth
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil{
		app.internalServerError(w,r,err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil{
		app.internalServerError(w,r,err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request){
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil{
		app.internalServerError(w,r,err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil{
		app.internalServerError(w,r,err)
		return
	}

}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request){
	postIDStr := chi.URLParam(r, "postID")

	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil{
		app.internalServerError(w,r,err)
		return
	}
	
	ctx := r.Context()

	err = app.store.Posts.Delete(ctx, postID)
	
	if err != nil {
		switch{
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w,r,err)

			default:
				app.internalServerError(w,r,err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct{
	Title *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error	"Bad request"
//	@Failure		401		{object}	error	"Unauthorized"
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [put]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request){
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil{
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil{
		app.badRequestError(w, r, err)
		return
	}

	if payload.Content != nil{
		post.Content = *payload.Content
	}

	if payload.Title != nil{
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil{
		app.internalServerError(w,r,err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil{
		app.internalServerError(w,r,err)
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		postIDStr := chi.URLParam(r, "postID")

		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil{
			app.internalServerError(w,r,err)
			return
		}
		
		ctx := r.Context()
		
		post, err := app.store.Posts.GetByID(ctx, postID)

		if err != nil {
			switch{
				case errors.Is(err, store.ErrNotFound):
					app.notFoundError(w,r,err)

				default:
					app.internalServerError(w,r,err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post{
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}

