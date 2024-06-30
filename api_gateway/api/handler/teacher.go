package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/api/helpers"
	"user_api_gateway/genproto/user_service"
	"user_api_gateway/pkg/etc"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router        /CreateTeacher [post]
// @Summary       Create teacher
// @Description   API for creating teacher
// @Tags          teacher
// @Accept        json
// @Produce       json
// @Param         teacher body user_service.CreateTeacher true "Teacher"
// @Success       200 {object} user_service.CreateTeacher
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) CreateTeacher(c *gin.Context) {
	var (
		req  user_service.CreateTeacher
		resp *user_service.Teacher
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	if err := helpers.ValidatePhone(req.Phone); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while validating phone number"+req.Phone)
		return
	}

	hashedPassword, err := etc.GeneratePasswordHash(req.Password)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while hashing password!")
		return
	}

	req.Password = string(hashedPassword)

	resp, err = h.grpcClient.TeacherService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create teacher")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListTeacher [GET]
// @Summary        Get list of teachers
// @Description    API for getting list of teachers
// @Tags           teacher
// @Accept         json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} user_service.GetListTeacherResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListTeacher(c *gin.Context) {
	var (
		req  user_service.GetListTeacherRequest
		resp *user_service.GetListTeacherResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager")
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

	resp, err = h.grpcClient.TeacherService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetByIdTeacher/{id} [GET]
// @Summary        Get a single teacher by ID
// @Description    API for getting a single teacher by ID
// @Tags           teacher
// @Accept         json
// @Produce        json
// @Param          id path string true "Teacher ID"
// @Success        200 {object} user_service.Teacher
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetTeacherByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *user_service.Teacher
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager")
		return
	}

	req := &user_service.TeacherPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.TeacherService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateTeacher/{id} [PUT]
// @Summary         Update a teacher by ID
// @Description     API for updating a teacher by ID
// @Tags            teacher
// @Accept          json
// @Produce         json
// @Param           id path string true "Teacher ID"
// @Param           teacher body user_service.UpdateTeacher true "Teacher"
// @Success         200 {object} user_service.UpdateTeacher
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateTeacher(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  user_service.UpdateTeacher
		resp *user_service.Teacher
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.TeacherService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteTeacher/{id} [DELETE]
// @Summary       Delete a teacher by ID
// @Description   API for deleting a teacher by ID
// @Tags          teacher
// @Accept        json
// @Produce       json
// @Param         id path string true "Teacher ID"
// @Success       200 {object} user_service.EmptyTeacher
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteTeacher(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &user_service.EmptyTeacher{}
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager")
		return
	}

	req := &user_service.TeacherPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.TeacherService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /TeacherReportList [get]
// @Summary       Get List of Teachers
// @Description   API for getting a list of teachers
// @Tags          report
// @Accept        json
// @Produce       json
// @Param         limit query string true "Limit"
// @Param         page query string true "Page"
// @Param         search query string false "Search term"
// @Success       200 {object} user_service.GetReportListTeacherResponse
// @Failure       400 {object} models.ResponseError "Invalid query parameters"
// @Failure       500 {object} models.ResponseError "Internal server error"
func (h *handler) GetReportListTeacher(c *gin.Context) {
	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not authorized")
		return
	}
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid limit parameter")
		return
	}

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid page parameter")
		return
	}

	req := user_service.GetReportListTeacherRequest{
		Limit:  int64(limit),
		Page:   int64(page),
		Search: c.Query("search"),
	}

	resp, err := h.grpcClient.TeacherService().GetReportList(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
