package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/genproto/user_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router        /CreateBranch [post]
// @Summary       Create branch
// @Description   API for creating branch
// @Tags          branch
// @Accept        json
// @Produce       json
// @Param         branch body user_service.CreateBranch true "Branch"
// @Success       200 {object} user_service.CreateBranch
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) CreateBranch(c *gin.Context) {
	var (
		req  user_service.CreateBranch
		resp *user_service.Branch
		err  error
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

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err = h.grpcClient.BranchService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create branch")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListBranch [GET]
// @Summary        Get list of branches
// @Description    API for getting list of branches
// @Tags           branch
// @Accept         json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} user_service.GetListBranchResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListBranch(c *gin.Context) {
	var (
		req  user_service.GetListBranchRequest
		resp *user_service.GetListBranchResponse
		err  error
	)
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

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin")
		return
	}

	resp, err = h.grpcClient.BranchService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetByIdBranch/{id} [GET]
// @Summary        Get a single branch by ID
// @Description    API for getting a single branch by ID
// @Tags           branch
// @Accept         json
// @Produce        json
// @Param          id path string true "Branch ID"
// @Success        200 {object} user_service.Branch
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetBranchByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *user_service.Branch
		err  error
	)

	req := &user_service.BranchPrimaryKey{
		Id: id,
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

	resp, err = h.grpcClient.BranchService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /UpdateBranch/{id} [PUT]
// @Summary         Update a branch by ID
// @Description     API for updating a branch by ID
// @Tags            branch
// @Accept          json
// @Produce         json
// @Param           id path string true "Branch ID"
// @Param           branch body user_service.UpdateBranch true "Branch"
// @Success         200 {object} user_service.UpdateBranch
// @Failure         404 {object} models.ResponseError
// @Failure         500 {object} models.ResponseError
func (h *handler) UpdateBranch(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  user_service.UpdateBranch
		resp *user_service.Branch
		err  error
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

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.BranchService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteBranch/{id} [DELETE]
// @Summary       Delete a branch by ID
// @Description   API for deleting a branch by ID
// @Tags          branch
// @Accept        json
// @Produce       json
// @Param         id path string true "Branch ID"
// @Success       200 {object} user_service.EmptyBranch
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteBranch(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &user_service.EmptyBranch{}
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

	req := &user_service.BranchPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.BranchService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
