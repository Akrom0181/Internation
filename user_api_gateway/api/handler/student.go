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
// @Router        /CreateStudent [post]
// @Summary       Create student
// @Description   API for creating student
// @Tags          student
// @Accept        json
// @Produce       json
// @Param         student body user_service.CreateStudent true "Student"
// @Success 200   {object} user_service.CreateStudent
// @Failure 404   {object} models.ResponseError
// @Failure 500   {object} models.ResponseError
func (h *handler) CreateStudent(c *gin.Context) {
	var (
		req  user_service.CreateStudent
		resp *user_service.Student
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager or an Admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	if err := helpers.ValidatePhone(req.Phone); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "input valid phone number"+req.Phone)
		return
	}

	hashedPassword, err := etc.GeneratePasswordHash(req.Password)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while hashing password!")
		return
	}

	req.Password = string(hashedPassword)

	resp, err = h.grpcClient.StudentService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create student")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListStudent [GET]
// @Summary        Get list of students
// @Description    API for getting list of students
// @Tags           student
// @Accept         json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} user_service.GetListStudentResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListStudent(c *gin.Context) {
	var (
		req  user_service.GetListStudentRequest
		resp *user_service.GetListStudentResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager or an admin")
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

	resp, err = h.grpcClient.StudentService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetByIdStudent/{id} [GET]
// @Summary        Get a single student by ID
// @Description    API for getting a single student by ID
// @Tags           student
// @Accept         json
// @Produce        json
// @Param          id path string true "Student ID"
// @Success        200 {object} user_service.Student
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetStudentByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *user_service.Student
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager or an admin")
		return
	}

	req := &user_service.StudentPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.StudentService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateStudent/{id} [PUT]
// @Summary         Update a student by ID
// @Description     API for updating a student by ID
// @Tags            student
// @Accept          json
// @Produce         json
// @Param           id path string true "Student ID"
// @Param           student body user_service.UpdateStudent true "Student"
// @Success         200 {object} user_service.UpdateStudent
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateStudent(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  user_service.UpdateStudent
		resp *user_service.Student
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager or an admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.StudentService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteStudent/{id} [DELETE]
// @Summary       Delete a student by ID
// @Description   API for deleting a student by ID
// @Tags          student
// @Accept        json
// @Produce       json
// @Param         id path string true "Student ID"
// @Success       200 {object} user_service.StudentEmpty
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteStudent(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &user_service.StudentEmpty{}
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager or an admin")
		return
	}

	req := &user_service.StudentPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.StudentService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /StudentReportList [get]
// @Summary       Get List of Students
// @Description   API for getting a list of students
// @Tags          report
// @Accept        json
// @Produce       json
// @Param         limit query string true "Limit"
// @Param         page query string true "Page"
// @Param         search query string false "Search term"
// @Success       200 {object} user_service.GetListStudentResponse
// @Failure       400 {object} models.ResponseError "Invalid query parameters"
// @Failure       500 {object} models.ResponseError "Internal server error"
func (h *handler) GetReportListStudent(c *gin.Context) {
	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin"  {
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

	req := user_service.GetReportListStudentRequest{
		Limit:  int64(limit),
		Page:   int64(page),
		Search: c.Query("search"),
	}

	resp, err := h.grpcClient.StudentService().GetReportList(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

