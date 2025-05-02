package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nikhilkarle/social/internal/store"
)

type userKey string
const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request){
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil{
		app.internalServerError(w,r,err)
	}
}

type FollowUser struct{
	UserID int64 `json:"user_id"`
}

func(app *application) followUserHandler(w http.ResponseWriter, r *http.Request){
	followerUser := getUserFromCtx(r)
	
	//revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w,r, &payload); err != nil{
		app.badRequestError(w,r,err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, payload); err != nil{
		app.internalServerError(w,r,err)
	}

	if err := app.store.Followers.Follow(r.Context(), followerUser.ID, payload.UserID); err != nil{
		switch err{
		case store.ErrConflict:
			app.conflictError(w,r,err)
			log.Println("CHIGGA")
			return

		default:
			app.internalServerError(w,r,err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil{
		app.internalServerError(w,r,err)
		return
	}
}

func(app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request){
	unfollowedUser := getUserFromCtx(r)
	
	//revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w,r, &payload); err != nil{
		app.badRequestError(w,r,err)
		return
	}
	

	if err := app.jsonResponse(w, http.StatusNoContent, payload); err != nil{
		app.internalServerError(w,r,err)
	}

	app.store.Followers.Unfollow(r.Context(), unfollowedUser.ID, payload.UserID)

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil{
		app.internalServerError(w,r,err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		userIDStr := chi.URLParam(r, "userID")

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil{
			app.internalServerError(w,r,err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)

		if err != nil{
			switch{
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w,r,err)

			default:
				app.internalServerError(w,r,err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User{
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}