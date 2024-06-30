package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router           /CreateSchedule [post]
// @Summary          Create schedule
// @Description      API for creating schedule
// @Tags             schedule
// @Accept           json
// @Produce          json
// @Param            schedule body schedule_service.CreateSchedule true "Schedule"
// @Success          200 {object} schedule_service.GetSchedule
// @Failure          404 {object} models.ResponseError
// @Failure          500 {object} models.ResponseError
func (h *handler) CreateSchedule(c *gin.Context) {
	var (
		req  schedule_service.CreateSchedule
		resp *schedule_service.GetSchedule
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, err, "You are not a SuperAdmin or a manager or an admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err = h.grpcClient.ScheduleService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create schedule")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /GetSchedule/{id} [get]
// @Summary         Get a schedule by ID
// @Description     API for getting a schedule by ID
// @Tags            schedule
// @Accept          json
// @Produce         json
// @Param           id path string true "Schedule ID"
// @Success         200 {object} schedule_service.GetSchedule
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) GetScheduleByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetSchedule
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, err, "You are not a SuperAdmin or a manager or an admin")
		return
	}

	req := &schedule_service.SchedulePrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.ScheduleService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get schedule")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListSchedule [get]
// @Summary        Get list of schedules
// @Description    API for getting list of schedules
// @Tags           schedule
// @Accept         json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} schedule_service.GetListScheduleResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListSchedule(c *gin.Context) {
	var (
		req  schedule_service.GetListScheduleRequest
		resp *schedule_service.GetListScheduleResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, err, "You are not a SuperAdmin or a manager or an admin")
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

	resp, err = h.grpcClient.ScheduleService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get list of schedules")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /UpdateSchedule/{id} [put]
// @Summary        Update a schedule by ID
// @Description    API for updating a schedule by ID
// @Tags           schedule
// @Accept         json
// @Produce        json
// @Param          id path string true "Schedule ID"
// @Param          schedule body schedule_service.UpdateSchedule true "Schedule"
// @Success        200 {object} schedule_service.GetSchedule
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) UpdateSchedule(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateSchedule
		resp *schedule_service.GetSchedule
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, err, "You are not a SuperAdmin or a manager or an admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.ScheduleService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to update schedule")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /DeleteSchedule/{id} [delete]
// @Summary         Delete a schedule by ID
// @Description     API for deleting a schedule by ID
// @Tags            schedule
// @Accept          json
// @Produce         json
// @Param           id path string true "Schedule ID"
// @Success         200 {object} schedule_service.EmptySchedule
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) DeleteSchedule(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp = &schedule_service.EmptySchedule{}
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	req := &schedule_service.SchedulePrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.ScheduleService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to delete schedule")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetScheduleForWeek [get]
// @Summary        Get schedules for a specific week
// @Description    API for getting schedules for a specific week
// @Tags           schedule
// @Accept         json
// @Produce        json
// @Param          weekStartDate query string true "Week Start Date"
// @Param          weekEndDate query string true "Week End Date"
// @Success        200 {object} schedule_service.GetListScheduleResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetScheduleForWeek(c *gin.Context) {
	var (
		resp *schedule_service.GetListScheduleResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin or a teacher")
		return
	}

	weekStartDate := c.Query("weekStartDate")
	weekEndDate := c.Query("weekEndDate")

	if weekStartDate == "" || weekEndDate == "" {
		handleGrpcErrWithDescription(c, h.log, fmt.Errorf("missing required parameters"), "missing required parameters")
		return
	}

	req := &schedule_service.GetScheduleForWeekRequest{
		WeekStartDate: weekStartDate,
		WeekEndDate:   weekEndDate,
	}

	resp, err = h.grpcClient.ScheduleService().GetScheduleForWeek(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "failed to get schedules for the week")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetScheduleForMonth [get]
// @Summary        Get schedules for a specific month
// @Description    API for getting schedules for a specific month
// @Tags           schedule
// @Accept         json
// @Produce        json
// @Param          monthStartDate query string true "Month Start Date"
// @Param          monthEndDate query string true "Month End Date"
// @Success        200 {object} schedule_service.GetListScheduleResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetScheduleForMonth(c *gin.Context) {
	var (
		resp *schedule_service.GetListScheduleResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin or a teacher")
		return
	}

	monthStartDate := c.Query("monthStartDate")
	monthEndDate := c.Query("monthEndDate")

	if monthStartDate == "" || monthEndDate == "" {
		handleGrpcErrWithDescription(c, h.log, fmt.Errorf("missing required parameters"), "missing required parameters")
		return
	}

	req := &schedule_service.GetScheduleForMonthRequest{
		MonthStartDate: monthStartDate,
		MonthEndDate:   monthEndDate,
	}

	resp, err = h.grpcClient.ScheduleService().GetScheduleForMonth(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get schedules for the month")
		return
	}

	c.JSON(http.StatusOK, resp)
}
