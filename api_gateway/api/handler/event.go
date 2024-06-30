package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/api/helpers"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router         /CreateEvent [post]
// @Summary        Create event
// @Description    API for creating event
// @Tags           event
// @Accept         json
// @Produce        json
// @Param          event body schedule_service.CreateEvent true "Event"
// @Success        200 {object} schedule_service.CreateEvent
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) CreateEvent(c *gin.Context) {
	var (
		req  schedule_service.CreateEvent
		resp *schedule_service.GetEvent
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	if err := helpers.IsSunday(req.Date); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "this date is not sunday!")
		return
	}

	resp, err = h.grpcClient.EventService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create event")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /GetListEvent [GET]
// @Summary       Get list of events
// @Description   API for getting list of events
// @Tags          event
// @Accept        json
// @Produce       json
// @Param         search query string false "Search"
// @Param         page query int false "Page"
// @Param         limit query int false "Limit"
// @Success       200 {object} schedule_service.GetListEventResponse
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) GetListEvent(c *gin.Context) {
	var (
		req  schedule_service.GetListEventRequest
		resp *schedule_service.GetListEventResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin")
		return
	}

	req.Search = c.Query("search")

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while parsing page")
		return
	}

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while parsing limit")
		return
	}

	req.Page = page
	req.Limit = limit

	resp, err = h.grpcClient.EventService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /GetByIdEvent/{id} [GET]
// @Summary         Get a single event by ID
// @Description     API for getting a single event by ID
// @Tags            event
// @Accept          json
// @Produce         json
// @Param           id path string true "Event ID"
// @Success         200 {object} schedule_service.GetEvent
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) GetEventByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetEvent
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin")
		return
	}

	req := &schedule_service.EventPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.EventService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /UpdateEvent/{id} [PUT]
// @Summary        Update an event by ID
// @Description    API for updating an event by ID
// @Tags           event
// @Accept         json
// @Produce        json
// @Param          id path string true "Event ID"
// @Param          event body schedule_service.UpdateEvent true "Event"
// @Success        200 {object} schedule_service.UpdateEvent
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) UpdateEvent(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateEvent
		resp *schedule_service.GetEvent
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, err, "You are not a SuperAdmin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.EventService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteEvent/{id} [DELETE]
// @Summary       Delete an event by ID
// @Description   API for deleting an event by ID
// @Tags          event
// @Accept        json
// @Produce       json
// @Param         id path string true "Event ID"
// @Success       200 {object} schedule_service.EmptyEvent
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteEvent(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &schedule_service.EmptyEvent{}
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin")
		return
	}

	req := &schedule_service.EventPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.EventService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
