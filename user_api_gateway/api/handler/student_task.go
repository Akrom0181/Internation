package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router         /CreateStudentTask [post]
// @Summary        Create student task
// @Description    API for creating student task
// @Tags           student_task
// @Accept         json
// @Produce        json
// @Param          student_task body schedule_service.CreateStudentTask true "Student Task"
// @Success        200 {object} schedule_service.GetStudentTask
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) CreateStudentTask(c *gin.Context) {
	var (
		req  schedule_service.CreateStudentTask
		resp *schedule_service.GetStudentTask
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Teacher")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err = h.grpcClient.StudentTaskService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create student task")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetStudentTask/{id} [get]
// @Summary        Get a student task by ID
// @Description    API for getting a student task by ID
// @Tags 		   student_task
// @Accept         json
// @Produce        json
// @Param          id path string true "Student Task ID"
// @Success        200 {object} schedule_service.GetStudentTask
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetStudentTaskByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetStudentTask
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Teacher")
		return
	}

	req := &schedule_service.StudentTaskPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.StudentTaskService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get student task")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router       /GetListStudentTask [get]
// @Summary      Get list of student tasks
// @Description  API for getting list of student tasks
// @Tags         student_task
// @Accept       json
// @Produce      json
// @Param        search query string false "Search"
// @Param        page query int false "Page"
// @Param        limit query int false "Limit"
// @Success      200 {object} schedule_service.GetListStudentTaskResponse
// @Failure      404 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
func (h *handler) GetListStudentTask(c *gin.Context) {
	var (
		req  schedule_service.GetListStudentTaskRequest
		resp *schedule_service.GetListStudentTaskResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Teacher")
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

	resp, err = h.grpcClient.StudentTaskService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get list of student tasks")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /UpdateStudentTask/{id} [put]
// @Summary       Update a student task by ID
// @Description   API for updating a student task by ID
// @Tags          student_task
// @Accept        json
// @Produce       json
// @Param         id path string true "Student Task ID"
// @Param         student_task body schedule_service.UpdateStudentTask true "Student Task"
// @Success       200 {object} schedule_service.GetStudentTask
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) UpdateStudentTask(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateStudentTask
		resp *schedule_service.GetStudentTask
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Teacher")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.StudentTaskService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to update student task")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router       /DeleteStudentTask/{id} [delete]
// @Summary      Delete a student task by ID
// @Description  API for deleting a student task by ID
// @Tags         student_task
// @Accept       json
// @Produce      json
// @Param        id path string true "Student Task ID"
// @Success      200 {object} schedule_service.EmptyStudentTask
// @Failure      404 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
func (h *handler) DeleteStudentTask(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp = &schedule_service.EmptyStudentTask{}
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Teacher" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Teacher")
		return
	}

	req := &schedule_service.StudentTaskPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.StudentTaskService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to delete student task")
		return
	}

	c.JSON(http.StatusOK, resp)
}
