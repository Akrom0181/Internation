package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router        /CreateStudentPayment [post]
// @Summary       Create student payment
// @Description   API for creating student payment
// @Tags          student_payment
// @Accept        json
// @Produce       json
// @Param         student_payment body schedule_service.CreateStudentPayment true "Student Payment"
// @Success       200 {object} schedule_service.GetStudentPayment
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) CreateStudentPayment(c *gin.Context) {
	var (
		req  schedule_service.CreateStudentPayment
		resp *schedule_service.GetStudentPayment
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Administration" && data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or an admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err = h.grpcClient.StudentPaymentService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create student payment")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /GetStudentPayment/{id} [get]
// @Summary       Get a student payment by ID
// @Description   API for getting a student payment by ID
// @Tags          student_payment
// @Accept        json
// @Produce       json
// @Param 		  id path string true "Student Payment ID"
// @Success       200 {object} schedule_service.GetStudentPayment
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) GetStudentPaymentByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetStudentPayment
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Administration" && data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or an admin")
		return
	}

	req := &schedule_service.StudentPaymentPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.StudentPaymentService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get student payment")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListStudentPayment [get]
// @Summary        Get list of student payments
// @Description    API for getting list of student payments
// @Tags           student_payment
// @Accept 		   json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} schedule_service.GetListStudentPaymentResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListStudentPayment(c *gin.Context) {
	var (
		req  schedule_service.GetListStudentPaymentRequest
		resp *schedule_service.GetListStudentPaymentResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Administration" && data.UserRole != "Manager" && data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
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

	resp, err = h.grpcClient.StudentPaymentService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get list of student payments")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateStudentPayment/{id} [put]
// @Summary         Update a student payment by ID
// @Description     API for updating a student payment by ID
// @Tags            student_payment
// @Accept          json
// @Produce         json
// @Param           id path string true "Student Payment ID"
// @Param           student_payment body schedule_service.UpdateStudentPayment true "Student Payment"
// @Success         200 {object} schedule_service.GetStudentPayment
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateStudentPayment(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateStudentPayment
		resp *schedule_service.GetStudentPayment
		err  error
	)
	
	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Administration" && data.UserRole != "Manager" && data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.StudentPaymentService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to update student payment")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /DeleteStudentPayment/{id} [delete]
// @Summary         Delete a student payment by ID
// @Description     API for deleting a student payment by ID
// @Tags            student_payment
// @Accept          json
// @Produce         json
// @Param           id path string true "Student Payment ID"
// @Success         200 {object} schedule_service.EmptyStudentPayment
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) DeleteStudentPayment(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp = &schedule_service.EmptyStudentPayment{}
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "Manager" && data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	req := &schedule_service.StudentPaymentPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.StudentPaymentService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to delete student payment")
		return
	}

	c.JSON(http.StatusOK, resp)
}
