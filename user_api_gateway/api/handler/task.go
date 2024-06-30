package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router         /CreateTask [post]
// @Summary        Create task
// @Description    API for creating task
// @Tags           task
// @Accept         json
// @Produce        json
// @Param          task body schedule_service.CreateTask true "Task"
// @Success        200 {object} schedule_service.GetTask
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) CreateTask(c *gin.Context) {
	var (
		req  schedule_service.CreateTask
		resp *schedule_service.GetTask
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a teacher")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err = h.grpcClient.TaskService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create task")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetTask/{id} [get]
// @Summary        Get a task by ID
// @Description    API for getting a task by ID
// @Tags           task
// @Accept         json
// @Produce        json
// @Param          id path string true "Task ID"
// @Success        200 {object} schedule_service.GetTask
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetTaskByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetTask
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a teacher")
		return
	}

	req := &schedule_service.TaskPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.TaskService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get task")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /GetListTask [get]
// @Summary         Get list of tasks
// @Description     API for getting list of tasks
// @Tags            task
// @Accept          json
// @Produce         json
// @Param           search query string false "Search"
// @Param           page query int false "Page"
// @Param           limit query int false "Limit"
// @Success         200 {object} schedule_service.GetListTaskResponse
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) GetListTask(c *gin.Context) {
	var (
		req  schedule_service.GetListTaskRequest
		resp *schedule_service.GetListTaskResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a teacher")
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

	resp, err = h.grpcClient.TaskService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get list of tasks")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /UpdateTask/{id} [put]
// @Summary        Update a task by ID
// @Description    API for updating a task by ID
// @Tags           task
// @Accept         json
// @Produce        json
// @Param          id path string true "Task ID"
// @Param          task body schedule_service.UpdateTask true "Task"
// @Success        200 {object} schedule_service.GetTask
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) UpdateTask(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateTask
		resp *schedule_service.GetTask
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a teacher")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.TaskService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to update task")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /DeleteTask/{id} [delete]
// @Summary        Delete a task by ID
// @Description    API for deleting a task by ID
// @Tags           task
// @Accept         json
// @Produce        json
// @Param          id path string true "Task ID"
// @Success        200 {object} schedule_service.EmptyTask
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) DeleteTask(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp = &schedule_service.EmptyTask{}
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a teacher")
		return
	}

	req := &schedule_service.TaskPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.TaskService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to delete task")
		return
	}

	c.JSON(http.StatusOK, resp)
}
