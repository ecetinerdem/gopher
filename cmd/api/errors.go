package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusNotFound, "the server cannot find the resource")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusConflict, "the server has a resource conflict")
}

func (app *application) unAuthorizedError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unAuthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("www-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request) {

	app.logger.Warnw("forbidden error", "method", r.Method, "path", r.URL.Path, "error")
	writeJsonError(w, http.StatusForbidden, "forbidden")
}

func (app *application) rateLimitExceedResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceed", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	writeJsonError(w, http.StatusTooManyRequests, "rate limit exceed, retry after: "+retryAfter)
}
