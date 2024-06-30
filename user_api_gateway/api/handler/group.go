package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router         /CreateGroup [post]
// @Summary        Create group
// @Description    API for creating group
// @Tags           group
// @Accept         json
// @Produce        json
// @Param          group body schedule_service.CreateGroup true "Group"
// @Success        200 {object} schedule_service.CreateGroup
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) CreateGroup(c *gin.Context) {
	var (
		req  schedule_service.CreateGroup
		resp *schedule_service.GetGroup
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

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err = h.grpcClient.GroupService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create group")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetListGroup [GET]
// @Summary        Get list of groups
// @Description    API for getting list of groups
// @Tags           group
// @Accept         json
// @Produce        json
// @Param          search query string false "Search"
// @Param          page query int false "Page"
// @Param          limit query int false "Limit"
// @Success        200 {object} schedule_service.GetListGroupResponse
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetListGroup(c *gin.Context) {
	var (
		req  schedule_service.GetListGroupRequest
		resp *schedule_service.GetListGroupResponse
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

	resp, err = h.grpcClient.GroupService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /GetByIdGroup/{id} [GET]
// @Summary        Get a single group by ID
// @Description    API for getting a single group by ID
// @Tags           group
// @Accept         json
// @Produce        json
// @Param          id path string true "Group ID"
// @Success        200 {object} schedule_service.GetGroup
// @Failure        404 {object} models.ResponseError
// @Failure        500 {object} models.ResponseError
func (h *handler) GetGroupByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetGroup
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

	req := &schedule_service.GroupPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.GroupService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router       /UpdateGroup/{id} [PUT]
// @Summary      Update a group by ID
// @Description  API for updating a group by ID
// @Tags         group
// @Accept       json
// @Produce      json
// @Param        id path string true "Group ID"
// @Param        group body schedule_service.UpdateGroup true "Group"
// @Success      200 {object} schedule_service.UpdateGroup
// @Failure      404 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
func (h *handler) UpdateGroup(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateGroup
		resp *schedule_service.GetGroup
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

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.GroupService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router       /DeleteGroup/{id} [DELETE]
// @Summary      Delete a group by ID
// @Description  API for deleting a group by ID
// @Tags         group
// @Accept       json
// @Produce      json
// @Param        id path string true "Group ID"
// @Success      200 {object} schedule_service.EmptyGroup
// @Failure      404 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
func (h *handler) DeleteGroup(c *gin.Context) {
	var (
		id   = c.Param("id")
		err  error
		resp = &schedule_service.EmptyGroup{}
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager")
		return
	}

	req := &schedule_service.GroupPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.GroupService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router          /GroupTeacher/{id} [get]
// @Summary         Get Groups by Teacher ID
// @Description     Get Groups associated with a Teacher by Teacher ID
// @Produce         json
// @Tags            group
// @Param           id path string true "Teacher ID"
// @Success         200 {object} schedule_service.GetListGroupResponse
// @Failure         400 {object} models.ResponseError "Invalid request body"
// @Failure         500 {object} models.ResponseError "Internal server error"
func (h *handler) GetGroupByIDTeacher(c *gin.Context) {
	id := c.Param("id")
	req := &schedule_service.TeacherID{Id: id}

	resp, err := h.grpcClient.GroupService().GetByIDTeacher(c, req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
