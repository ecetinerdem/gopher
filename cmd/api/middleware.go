package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ecetinerdem/gopherSocial/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuth() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("missing authorization header"))
				return
			}
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Basic" {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("malformed authorization header"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])

			if err != nil {
				app.unAuthorizedBasicError(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			password := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)

			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unAuthorizedError(w, r, fmt.Errorf("authorization header missing"))
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizedError(w, r, fmt.Errorf("malformed authorization header"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePresedence(r.Context(), user, requiredRole)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenError(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePresedence(ctx context.Context, user *store.User, requiredRole string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, requiredRole)

	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetUserByID(ctx, userID)
	}

	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = app.store.Users.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		err = app.cacheStorage.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.config.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceedResponse(w, r, retryAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
