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
// @Router        /CreateManager [post]
// @Summary       Create manager
// @Description   API for creating manager
// @Tags          manager
// @Accept        json
// @Produce       json
// @Param         manager body user_service.CreateManager true "Manager"
// @Success       200 {object} user_service.CreateManager
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) CreateManager(c *gin.Context) {
	var (
		req  user_service.CreateManager
		resp *user_service.Manager
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" || data.UserID != "e924cb31-e068-4062-a3b9-66790722e68a" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SUPER admin")
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

	resp, err = h.grpcClient.ManagerService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create manager")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListManager [GET]
// @Summary        Get list of managers
// @Description    API for getting list of managers
// @Tags           manager
// @Accept         json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} user_service.GetListManagerResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListManager(c *gin.Context) {
	var (
		req  user_service.GetListManagerRequest
		resp *user_service.GetListManagerResponse
		err  error
	)
	req.Search = c.Query("search")

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while parsing page")
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

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while parsing limit")
		return
	}

	req.Page = page
	req.Limit = limit

	resp, err = h.grpcClient.ManagerService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetByIdManager/{id} [GET]
// @Summary        Get a single manager by ID
// @Description    API for getting a single manager by ID
// @Tags           manager
// @Accept         json
// @Produce        json
// @Param          id path string true "Manager ID"
// @Success        200 {object} user_service.Manager
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetManagerByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *user_service.Manager
		err  error
	)

	req := &user_service.ManagerPrimaryKey{
		Id: id,
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

	resp, err = h.grpcClient.ManagerService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateManager/{id} [PUT]
// @Summary         Update a manager by ID
// @Description     API for updating a manager by ID
// @Tags            manager
// @Accept          json
// @Produce         json
// @Param           id path string true "Manager ID"
// @Param           manager body user_service.UpdateManager true "Manager"
// @Success         200 {object} user_service.UpdateManager
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateManager(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  user_service.UpdateManager
		resp *user_service.Manager
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

	if data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.ManagerService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteManager/{id} [DELETE]
// @Summary       Delete a manager by ID
// @Description   API for deleting a manager by ID
// @Tags          manager
// @Accept        json
// @Produce       json
// @Param         id path string true "Manager ID"
// @Success       200 {object} user_service.EmptyManager
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteManager(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &user_service.EmptyManager{}
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin")
		return
	}

	req := &user_service.ManagerPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.ManagerService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
