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
// @Router        /CreateSupportTeacher [post]
// @Summary       Create support teacher
// @Description   API for creating support teacher
// @Tags          support_teacher
// @Accept        json
// @Produce       json
// @Param         support_teacher body user_service.CreateSupportTeacher true "Support Teacher"
// @Success       200 {object} user_service.CreateSupportTeacher
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) CreateSupportTeacher(c *gin.Context) {
	var (
		req  user_service.CreateSupportTeacher
		resp *user_service.SupportTeacher
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, err, "You are not a SuperAdmin or a Manager")
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

	resp, err = h.grpcClient.SupportTeacherService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create support teacher")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListSupportTeacher [GET]
// @Summary        Get list of support teachers
// @Description    API for getting list of support teachers
// @Tags           support_teacher
// @Accept         json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} user_service.GetListSupportTeacherResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListSupportTeacher(c *gin.Context) {
	var (
		req  user_service.GetListSupportTeacherRequest
		resp *user_service.GetListSupportTeacherResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, err, "You are not a SuperAdmin or a Manager")
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

	resp, err = h.grpcClient.SupportTeacherService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetByIdSupportTeacher/{id} [GET]
// @Summary        Get a single support teacher by ID
// @Description    API for getting a single support teacher by ID
// @Tags           support_teacher
// @Accept         json
// @Produce        json
// @Param          id path string true "Support Teacher ID"
// @Success        200 {object} user_service.SupportTeacher
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetSupportTeacherByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *user_service.SupportTeacher
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

	req := &user_service.SupportTeacherPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.SupportTeacherService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateSupportTeacher/{id} [PUT]
// @Summary         Update a support teacher by ID
// @Description     API for updating a support teacher by ID
// @Tags            support_teacher
// @Accept          json
// @Produce         json
// @Param           id path string true "Support Teacher ID"
// @Param           support_teacher body user_service.UpdateSupportTeacher true "Support Teacher"
// @Success         200 {object} user_service.UpdateSupportTeacher
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateSupportTeacher(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  user_service.UpdateSupportTeacher
		resp *user_service.SupportTeacher
		err  error
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.SupportTeacherService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteSupportTeacher/{id} [DELETE]
// @Summary       Delete a support teacher by ID
// @Description   API for deleting a support teacher by ID
// @Tags          support_teacher
// @Accept        json
// @Produce       json
// @Param         id path string true "Support Teacher ID"
// @Success       200 {object} user_service.EmptySTeacher
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteSupportTeacher(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &user_service.EmptySTeacher{}
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

	req := &user_service.SupportTeacherPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.SupportTeacherService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /SupportTeacherReportList [get]
// @Summary        Get List of Support Teachers
// @Description    API for getting a list of support teachers
// @Tags           report
// @Accept         json
// @Produce        json
// @Param          limit query string true "Limit"
// @Param          page query string true "Page"
// @Param          search query string false "Search"
// @Success        200 {object} user_service.GetListSupportTeacherResponse
// @Failure        400 {object} models.ResponseError "Invalid query parameters"
// @Failure        500 {object} models.ResponseError "Internal server error"
func (h *handler) GetReportListSupportTeacher(c *gin.Context) {
	limitStr := c.Query("limit")

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not authorized")
		return
	}

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

	req := user_service.GetReportListSupportTeacherRequest{
		Limit:  int64(limit),
		Page:   int64(page),
		Search: c.Query("search"),
	}

	resp, err := h.grpcClient.SupportTeacherService().GetReportList(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

