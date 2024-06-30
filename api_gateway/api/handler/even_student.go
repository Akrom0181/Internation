package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"user_api_gateway/api/helpers"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router         /CreateEventStudent [post]
// @Summary        Create event student
// @Description    API for creating event student
// @Tags           event_student
// @Accept         json
// @Produce        json
// @Param          event_student body schedule_service.CreateEventStudent true "Event Student"
// @Success        200 {object} schedule_service.CreateEventStudent
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) CreateEventStudent(c *gin.Context) {
	var (
		req  schedule_service.CreateEventStudent
		resp *schedule_service.GetEventStudent
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Student" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a Student")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	currentTime := time.Now()

	if err := helpers.CheckEventRegistration(currentTime); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "you are late to register")
	}

	resp, err = h.grpcClient.EventStudentService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create event student")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /GetListEventStudent [GET]
// @Summary         Get list of event students
// @Description     API for getting list of event students
// @Tags            event_student
// @Accept          json
// @Produce         json
// @Param           search query string false "Search"
// @Param           page query int false "Page"
// @Param           limit query int false "Limit"
// @Success         200 {object} schedule_service.GetListEventStudentResponse
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) GetListEventStudent(c *gin.Context) {
	var (
		req  schedule_service.GetListEventStudentRequest
		resp *schedule_service.GetListEventStudentResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Student" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a Student")
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

	resp, err = h.grpcClient.EventStudentService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /GetByIdEventStudent/{id} [GET]
// @Summary         Get a single event student by ID
// @Description     API for getting a single event student by ID
// @Tags            event_student
// @Accept          json
// @Produce         json
// @Param           id path string true "Event Student ID"
// @Success         200 {object} schedule_service.GetEventStudent
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) GetEventStudentByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetEventStudent
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Student" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a Student")
		return
	}

	req := &schedule_service.EventStudentPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.EventStudentService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateEventStudent/{id} [PUT]
// @Summary         Update an event student by ID
// @Description     API for updating an event student by ID
// @Tags            event_student
// @Accept          json
// @Produce         json
// @Param           id path string true "Event Student ID"
// @Param           event_student body schedule_service.UpdateEventStudent true "Event Student"
// @Success         200 {object} schedule_service.UpdateEventStudent
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateEventStudent(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateEventStudent
		resp *schedule_service.GetEventStudent
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Student" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a Student")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.EventStudentService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /DeleteEventStudent/{id} [DELETE]
// @Summary        Delete an event student by ID
// @Description    API for deleting an event student by ID
// @Tags           event_student
// @Accept         json
// @Produce        json
// @Param          id path string true "Event Student ID"
// @Success        200 {object} schedule_service.EmptyEventStudent
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) DeleteEventStudent(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &schedule_service.EmptyEventStudent{}
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a Student")
		return
	}

	req := &schedule_service.EventStudentPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.EventStudentService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /EventStudent/{id} [get]
// @Summary       Get student with events by student ID
// @Description   API for getting a student with their events by student ID
// @Tags          event_student
// @Accept        json
// @Produce       json
// @Param         id path string true "Student ID"
// @Success       200 {object} schedule_service.GetStudentWithEventsResponse
// @Failure       404 {object} models.ResponseError "Student not found"
// @Failure       500 {object} models.ResponseError "Internal server error"
func (h *handler) GetStudentWithEventsByID(c *gin.Context) {

	id := c.Param("id")
	req := &schedule_service.EventStudentPrimaryKey{Id: id}

	resp, err := h.grpcClient.EventStudentService().GetStudentByID(c, req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
