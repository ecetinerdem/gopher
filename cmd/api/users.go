package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ecetinerdem/gopherSocial/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userParam := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userParam, 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.getUser(ctx, userID)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	err = app.writeJsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	followedUserID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Followers.Follow(ctx, followerUser.ID, followedUserID)
	if err != nil {
		switch err {
		case store.ErrDataConflict:
			app.conflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	err = app.writeJsonResponse(w, http.StatusNoContent, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	unFollowedUserID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Followers.Unfollow(ctx, followerUser.ID, unFollowedUserID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	err = app.writeJsonResponse(w, http.StatusNoContent, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.writeJsonResponse(w, http.StatusOK, ""); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(userParam, 10, 64)

		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetUserByID(ctx, userID)

		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}

		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)

	return user
}
