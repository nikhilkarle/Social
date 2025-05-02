package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nikhilkarle/social/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request){
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

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil{
		app.internalServerError(w,r,err)
	}
}