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
// @Router        /CreateAdministration [post]
// @Summary       Create administration
// @Description   API for creating administration
// @Tags          administration
// @Accept        json
// @Produce       json
// @Param         administration body user_service.CreateAdministration true "administration"
// @Success 200   {object} user_service.CreateAdministration
// @Failure 404   {object} models.ResponseError
// @Failure 500   {object} models.ResponseError
func (h *handler) CreateAdministration(c *gin.Context) {
	var (
		req  user_service.CreateAdministration
		resp *user_service.Administration
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

	resp, err = h.grpcClient.AdministrationService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create customer")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListAdministration [GET]
// @Summary        Get List administration
// @Description    API for getting list administration
// @Tags           administration
// @Accept         json
// @Produce        json
// @Param		   seller query string false "administration"
// @Param		   page query int false "page"
// @Param		   limit query int false "limit"
// @Success 200    {object} user_service.GetListAdministrationResponse
// @Failure 404    {object} models.ResponseError
// @Failure 500    {object} models.ResponseError
func (h *handler) GetListAdministration(c *gin.Context) {
	var (
		req  user_service.GetListAdministrationRequest
		resp *user_service.GetListAdministrationResponse
		err  error
	)
	req.Search = c.Query("search")

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a Manager")
		return
	}

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

	resp, err = h.grpcClient.AdministrationService().GetList(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetByIdAdministration/{id} [GET]
// @Summary        Get a single administration by ID
// @Description    API for getting a single administration by ID
// @Tags           administration
// @Accept         json
// @Produce        json
// @Param          id path string true "administration ID"
// @Success        200 {object} user_service.Administration
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetAdministrationByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *user_service.Administration
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

	req := &user_service.AdministrationPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.AdministrationService().GetByID(c, req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateAdministration/{id} [PUT]
// @Summary         Update a administration by ID
// @Description     API for updating a administration by ID
// @Tags            administration
// @Accept          json
// @Produce         json
// @Param           id path string true "administration ID"
// @Param           seller body user_service.UpdateAdministration true "Administration"
// @Success         200 {object} user_service.UpdateAdministration
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateAdministrationr(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  user_service.UpdateAdministration
		resp *user_service.Administration
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

	req.Id = id
	resp, err = h.grpcClient.AdministrationService().Update(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteAdministration/{id} [DELETE]
// @Summary       Delete a administration by ID
// @Description   API for deleting a administration by ID
// @Tags          administration
// @Accept        json
// @Produce       json
// @Param         id path string true "administration ID"
// @Success       200 {object} user_service.EmptyAdmin
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteAdministration(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &user_service.EmptyAdmin{}
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

	req := &user_service.AdministrationPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.AdministrationService().Delete(c, req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /AdministrationReportList [get]
// @Summary        Get List of Administrations
// @Description    API for getting a list of administrations
// @Tags           report
// @Accept         json
// @Produce        json
// @Param          limit query string true "Limit"
// @Param          page query string true "Page"
// @Param          search query string false "Search term"
// @Success        200 {object} user_service.GetReportAdministrationResponse
// @Failure        400 {object} models.ResponseError "Invalid query parameters"
// @Failure        500 {object} models.ResponseError "Internal server error"
func (h *handler) GetReportListAdministration(c *gin.Context) {
	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
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

	req := user_service.GetReportListAdministrationRequest{
		Limit:  int64(limit),
		Page:   int64(page),
		Search: c.Query("search"),
	}

	resp, err := h.grpcClient.AdministrationService().GetReportList(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
