package response

import (
	"math"
	"net/http"
)

// Ok returns a 200 OK JSON response
func (resp *respInfo) Ok(body any) {
	_ = resp.WriteResponse(http.StatusOK, body, nil)
}

// Created returns a 201 Created JSON response
func (resp *respInfo) Created(body any) {
	_ = resp.WriteResponse(http.StatusCreated, body, nil)
}

// OkPagination returns a 200 OK With Pagination
func (resp *respInfo) OkPagination(body any, page, limit, total int) {
	if page <= 0 {
		page = 1
	}
	resPag := &Pagination{
		Page:  page,
		Total: 1,
	}
	if limit > 0 {
		resPag.Total = int(math.Ceil(float64(total) / float64(limit)))
	}
	_ = resp.WriteResponse(http.StatusOK, body, resPag)
}

// Redirect301 returns a 301
func (resp *respInfo) Redirect301(body any) {
	_ = resp.WriteResponse(301, body, nil)
}

// BadRequest returns a 400 Bad Request JSON response
func (resp *respInfo) BadRequest(body RespBodyError) {
	_ = resp.WriteResponse(http.StatusBadRequest, body, nil)
}

// Unauthorized returns a 401 Unauthorized JSON response
func (resp *respInfo) Unauthorized(body RespBodyError) {
	_ = resp.WriteResponse(http.StatusUnauthorized, body, nil)
}

// Forbidden returns a 403 Forbidden JSON response
func (resp *respInfo) Forbidden(body RespBodyError) {
	_ = resp.WriteResponse(http.StatusForbidden, body, nil)
}

// NotFound returns a 404 Not Found JSON response
func (resp *respInfo) NotFound(body RespBodyError) {
	_ = resp.WriteResponse(http.StatusNotFound, body, nil)
}

// InternalServerError returns a 500 Internal Server Error JSON response
func (resp *respInfo) InternalServerError(body RespBodyError) {
	_ = resp.WriteResponse(http.StatusInternalServerError, body, nil)
}
